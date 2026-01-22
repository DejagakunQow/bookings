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
    // ADD BOOKING
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

            pendingBookings.push({
                roomID: parseInt(roomSelect.value),
                roomName: roomSelect.options[roomSelect.selectedIndex].text,
                guest,
                start,
                end
            });

            renderPreview();
        });
    }

    // ===============================
    // SAVE CALENDAR
    // ===============================
    const saveBtn = document.getElementById("save-calendar");
    if (saveBtn) {
        saveBtn.addEventListener("click", async function () {

            if (pendingBookings.length === 0 && !editingReservationID) {
                alert("No changes to save");
                return;
            }

            const token = document.getElementById("csrf-token").value;

            const url = editingReservationID
                ? `/admin/calendar/${editingReservationID}`
                : `/admin/calendar/0`;

            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-CSRF-Token": token
                },
                body: JSON.stringify(pendingBookings)
            });

            if (response.ok) {
                alert("Changes saved");
                location.reload();
            } else {
                alert("Failed to save calendar");
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
