console.log("Calendar JS loaded");

let pendingBookings = [];

// ===============================
// CLICK EXISTING RESERVATIONS
// ===============================
document.addEventListener("click", function (e) {

    const target = e.target.closest(".reservation");
    if (!target) return;

    const id = target.dataset.id;
    const roomID = target.dataset.roomId;
    const guest = target.dataset.guest;
    const start = target.dataset.start;
    const end = target.dataset.end;

    document.getElementById("room-select").value = roomID;
    document.getElementById("guest-name").value = guest;
    document.getElementById("start-date").value = start;
    document.getElementById("end-date").value = end;

    window.editingReservationID = id;
});

// ===============================
// PREVIEW PENDING BOOKINGS
// ===============================
function renderPreview() {
    const box = document.getElementById("booking-preview");
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
document.getElementById("add-booking").addEventListener("click", function () {

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
        guest: guest,
        start: start,
        end: end
    });

    renderPreview();
});

// ===============================
// SAVE / UPDATE TO SERVER
// ===============================
document.getElementById("save-calendar").addEventListener("click", async function () {

    if (pendingBookings.length === 0 && !window.editingReservationID) {
        alert("No changes to save");
        return;
    }

    const token = document.getElementById("csrf-token").value;

    const url = window.editingReservationID
        ? `/admin/calendar/update/${window.editingReservationID}`
        : "/admin/calendar/save";

    const payload = window.editingReservationID
        ? {
            roomID: document.getElementById("room-select").value,
            guest: document.getElementById("guest-name").value,
            start: document.getElementById("start-date").value,
            end: document.getElementById("end-date").value
        }
        : pendingBookings;

    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": token
        },
        body: JSON.stringify(payload)
    });

    if (response.ok) {
        alert("Changes saved successfully");
        location.reload();
    } else {
        const msg = await response.text();
        alert("Error: " + msg);
    }
});


document.getElementById("delete-reservation").addEventListener("click", async function () {

    if (!window.editingReservationID) {
        alert("Select a reservation first");
        return;
    }

    if (!confirm("Are you sure you want to delete this reservation?")) return;

    const token = document.getElementById("csrf-token").value;

    const response = await fetch(`/admin/calendar/delete/${window.editingReservationID}`, {
        method: "POST",
        headers: {
            "X-CSRF-Token": token
        }
    });

    if (response.ok) {
        alert("Reservation deleted");
        location.reload();
    } else {
        alert("Delete failed");
    }
});

const deleteBtn = document.getElementById("delete-reservation");

if (deleteBtn) {
    deleteBtn.onclick = async function () {

        if (!window.editingReservationID) {
            notify("Select a reservation first", "warning");
            return;
        }

        document.getElementById("delete-confirm-box").classList.remove("d-none");
    };
}


document.getElementById("confirm-delete").onclick = async function () {

    const token = document.getElementById("csrf-token").value;

    const response = await fetch(`/admin/calendar/delete/${window.editingReservationID}`, {
        method: "POST",
        headers: { "X-CSRF-Token": token }
    });

    if (response.ok) {
        notify("Reservation deleted successfully", "success");
        setTimeout(() => location.reload(), 800);
    } else {
        notify("Delete failed", "danger");
    }
};

document.getElementById("cancel-delete").onclick = function () {
    document.getElementById("delete-confirm-box").classList.add("d-none");
};

