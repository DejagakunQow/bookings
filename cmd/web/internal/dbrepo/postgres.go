package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/DejagakunQow/bookings/cmd/web/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 🛡️ Prevent double booking
	ok, err := m.IsRoomAvailable(res.RoomID, res.StartDate, res.EndDate)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("room is already booked for the selected dates")
	}

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
			end_date, room_id, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err = m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,	
			created_at, updated_at, restriction_id) 
			values
			($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `
		select
			count(id)
		from
			room_restrictions
		where
			room_id = $1
			and $2 < end_date and $3 > start_date;`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		select
			r.id, r.room_name
		from
			rooms r
		where r.id not in 
		(select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);
		`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		select id, room_name, created_at, updated_at from rooms where id = $1
`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}

func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, created_at, updated_at
	from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := ` 
		  update users set first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5
		  where id = $6
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil

}
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	query := `
		select 
			r.id, r.first_name, r.last_name, r.email, r.phone,
			r.start_date, r.end_date, r.room_id,
			rm.room_name,
			r.created_at, r.updated_at
		from reservations r
		left join rooms rm on rm.id = r.room_id
		order by r.start_date desc
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []models.Reservation

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(
			&r.ID,
			&r.FirstName,
			&r.LastName,
			&r.Email,
			&r.Phone,
			&r.StartDate,
			&r.EndDate,
			&r.RoomID,
			&r.Room.RoomName,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, r)
	}

	return reservations, nil
}

func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var r models.Reservation

	query := `
        SELECT id, first_name, last_name, email, phone, start_date, end_date, room_id, processed, created_at, updated_at
        FROM reservations
        WHERE id = $1
    `

	row := m.DB.QueryRow(query, id)

	err := row.Scan(
		&r.ID,
		&r.FirstName,
		&r.LastName,
		&r.Email,
		&r.Phone,
		&r.StartDate,
		&r.EndDate,
		&r.RoomID,
		&r.Processed,
		&r.CreatedAt,
		&r.UpdatedAt,
	)

	return r, err
}
func (m *postgresDBRepo) UpdateReservation(r models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
    UPDATE reservations
    SET first_name = $1,
        last_name = $2,
        room_id = $3,
        start_date = $4,
        end_date = $5,
        updated_at = $6
    WHERE id = $7
`

	_, err := m.DB.ExecContext(ctx, query,
		r.FirstName,
		r.LastName,
		r.RoomID,
		r.StartDate,
		r.EndDate,
		time.Now(),
		r.ID,
	)

	return err
}

func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := ` 
		  update reservations set processed = $1
		  where id = $2
	`

	_, err := m.DB.ExecContext(ctx, query,
		processed,
		id,
	)
	if err != nil {
		return err
	}
	return nil

}
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	query := `
		select 
			r.id,
			r.last_name,
			rm.room_name,
			r.start_date,
			r.end_date
		from reservations r
		left join rooms rm on rm.id = r.room_id
		where r.processed = 0
		order by r.start_date
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []models.Reservation

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(
			&r.ID,
			&r.LastName,
			&r.Room.RoomName,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, r)
	}

	return reservations, nil

}

func (m *postgresDBRepo) AllReservationsForDate(date time.Time) ([]models.Reservation, error) {
	query := `
		SELECT 
			r.id, r.first_name, r.last_name, r.email, r.phone,
			r.start_date, r.end_date, r.room_id,
			rm.room_name,
			r.created_at, r.updated_at
		FROM reservations r
		LEFT JOIN rooms rm ON rm.id = r.room_id
		WHERE $1 >= r.start_date AND $1 < r.end_date
		ORDER BY rm.room_name
	`

	rows, err := m.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []models.Reservation

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(
			&r.ID,
			&r.FirstName,
			&r.LastName,
			&r.Email,
			&r.Phone,
			&r.StartDate,
			&r.EndDate,
			&r.RoomID,
			&r.Room.RoomName,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, r)
	}

	return reservations, nil
}

func (m *postgresDBRepo) IsRoomAvailable(roomID int, start, end time.Time) (bool, error) {

	query := `
		SELECT COUNT(*)
		FROM reservations
		WHERE room_id = $1
		  AND $2 < end_date
		  AND $3 > start_date
	`
	var count int
	err := m.DB.QueryRow(query, roomID, start, end).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	query := `select id, room_name from rooms order by room_name`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var r models.Room
		err := rows.Scan(&r.ID, &r.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, r)
	}
	return rooms, nil
}

func (m *postgresDBRepo) CountReservations() (int, error) {
	var count int
	query := `select count(id) from reservations`
	err := m.DB.QueryRow(query).Scan(&count)
	return count, err
}

func (m *postgresDBRepo) CountNewReservations() (int, error) {
	var count int
	query := `select count(id) from reservations where processed = 0`
	err := m.DB.QueryRow(query).Scan(&count)
	return count, err
}
