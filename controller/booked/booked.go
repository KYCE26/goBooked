package booked

import (
	"encoding/json"
	"net/http"
	"strconv"

	"belajar/database"
	"belajar/model/booked"

	"github.com/gorilla/mux"
)

func GetBooked(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT * FROM bookeds")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var bookeds []booked.Booked
    for rows.Next() {
        var s booked.Booked
        if err := rows.Scan(&s.BookedId,&s.MotorId,&s.Nama,&s.Alamat,&s.Harga); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        bookeds = append(bookeds, s)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bookeds)
}

func PostBooked(w http.ResponseWriter, r *http.Request) {
	var ps booked.Booked
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Query untuk memasukan mahasiswa ke dalam table
	query := `
		INSERT INTO bookeds (motor_id, nama, alamat, harga) 
		VALUES (?, ?,?,?)`

	// Mengeksekusi query
	res, err := database.DB.Exec(query, ps.MotorId, ps.Nama, ps.Alamat, ps.Harga)
	if err != nil {
		http.Error(w, "Failed to insert booked: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ambil id terakhir
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booked added successfully",
		"id":      id,
	})
}

func PutBooked(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Decode JSON body
	var ps booked.Booked
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Query ubah booked
	query := `
		UPDATE bookeds 
		SET motor_id=?, nama=?,alamat=?,harga=?
		WHERE booked_id=?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query,ps.MotorId,ps.Nama,ps.Alamat, ps.Harga,id)
	if err != nil {
		http.Error(w, "Failed to update booked: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were updated
	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booked updated successfully",
	})
}

func DeleteBooked(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for deleting a category admin
	query := `
		DELETE FROM bookeds
		WHERE booked_id = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete booked: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booked deleted successfully",
	})
}

