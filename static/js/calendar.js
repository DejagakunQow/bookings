document.addEventListener("DOMContentLoaded", function () {

    console.log("Calendar JS loaded");

    let pendingBookings = [];
    let editingReservationID = null;

    // ===============================
    // CLICK EXISTING RESERVATIONS
    // ===============================
    document.addEventListener("click", function (e) {

        const target = e.target.closest(".reservation");
        if (!target) return;

        editingReservationID = target.dataset.id;

        document.getElementById("room-select").value = target.dataset.roomId;
        document.getElementById("guest-name").value = target.dataset.guest;
        document.getElementById("start-date").value = target.dataset.start;
        document.getElementById("end-date").value = target.dataset.end;

        pendingBookings = []; // clear pending when editing
        renderPreview();
    });

    // ===============================
    // PREVIEW PENDING BOOKINGS
    // ===============================
    function renderPreview() {
        const box = document.getElementById("booking-preview");
        if (!box) return;

        box.innerHTML = "";

        pendingBookings.forEach((b) => {
            const div = document.createElement("div");
            div.className = "reservation-badge";
            div.innerText = `${b.roomName} — ${b.guest} (${b.start} → ${b.end})`;
            box.appendChild(div);
        });
    }

    // ===============================
    // ADD BOOKING (NEW)
    // ===============================
    const addBtn = document.getElementById("add-booking");
    if (addBtn) {
        addBtn.addEventListener("click", function () {

            const roomSelect = document.getElementById("room-select");
            const guest = document.getElementById("guest-name").value;
            const start = document.getElementById("start-date").value;
            const end = document.getElementById("end-date").value;

            if (!roomSelect.value || !guest || !start || !end) {
                alert("Please complete all fields");
                return;
            }

            editingReservationID = null; // switching to add mode

            pendingBookings.push({
                room_id: parseInt(roomSelect.value),
                guest: guest,
                start_date: start,
                end_date: end,
                roomName: roomSelect.options[roomSelect.selectedIndex].text
            });

            renderPreview();
        });
    }

    // ===============================
    // SAVE CALENDAR (ADD OR UPDATE)
    // ===============================
    const saveBtn = document.getElementById("save-calendar");
    if (saveBtn) {
        saveBtn.addEventListener("click", async function () {

            const token = document.getElementById("csrf-token").value;

            // UPDATE EXISTING
            if (editingReservationID) {

                const payload = {
                    room_id: document.getElementById("room-select").value,
                    guest: document.getElementById("guest-name").value,
                    start_date: document.getElementById("start-date").value,
                    end_date: document.getElementById("end-date").value
                };

                const response = await fetch(`/admin/calendar/${editingReservationID}`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "X-CSRF-Token": token
                    },
                    body: JSON.stringify(payload)
                });

                if (response.ok) {
                    alert("Reservation updated");
                    location.reload();
                } else {
                    alert("Failed to update reservation");
                }

                return;
            }

            // ADD NEW BOOKINGS
            if (pendingBookings.length === 0) {
                alert("No bookings to save");
                return;
            }

            const response = await fetch(`/admin/calendar/0`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-CSRF-Token": token
                },
                body: JSON.stringify(pendingBookings)
            });

            if (response.ok) {
                alert("Bookings saved");
                location.reload();
            } else {
                alert("Failed to save bookings");
            }
        });
    }

    // ===============================
    // DELETE RESERVATION
    // ===============================
    const deleteBtn = document.getElementById("delete-reservation");
    if (deleteBtn) {
        deleteBtn.addEventListener("click", function () {

            if (!editingReservationID) {
                alert("Select a reservation first");
                return;
            }

            document.getElementById("delete-confirm-box").classList.remove("d-none");
        });
    }

    const confirmDelete = document.getElementById("confirm-delete");
    if (confirmDelete) {
        confirmDelete.addEventListener("click", async function () {

            const token = document.getElementById("csrf-token").value;

            const response = await fetch(`/admin/calendar/${editingReservationID}`, {
                method: "DELETE",
                headers: { "X-CSRF-Token": token }
            });

            if (response.ok) {
                alert("Reservation deleted");
                location.reload();
            } else {
                alert("Delete failed");
            }
        });
    }

    const cancelDelete = document.getElementById("cancel-delete");
    if (cancelDelete) {
        cancelDelete.addEventListener("click", function () {
            document.getElementById("delete-confirm-box").classList.add("d-none");
        });
    }

});
