package motor

import (
    "encoding/json"
    "net/http"
    "strconv"

    "belajar/database"
    "github.com/gorilla/mux"
    "belajar/model/motor"
)

func GetMotor(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT * FROM motors")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var motors []motor.Motor
    for rows.Next() {
        var c motor.Motor
        if err := rows.Scan(&c.MotorId,&c.MotorName,&c.Plat); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        motors = append(motors, c)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(motors)
}

func PostMotor(w http.ResponseWriter, r *http.Request) {
	var pc motor.Motor
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new motor
	query := `
		INSERT INTO motors (nama_motor, plat) 
		VALUES (?, ?)`

	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.MotorName, pc.Plat)
	if err != nil {
		http.Error(w, "Failed to insert motor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Motor added successfully",
		"id":      id,
	})
}

func PutMotor(w http.ResponseWriter, r *http.Request) {
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
	var pc motor.Motor
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for updating the category admin
	query := `
		UPDATE motors
		SET nama_motor=?, plat=?
		WHERE motor_id=?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, pc.MotorName, pc.Plat,id)
	if err != nil {
		http.Error(w, "Failed to update motor: "+err.Error(), http.StatusInternalServerError)
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
		"message": "Motor updated successfully",
	})
}

func DeleteMotor(w http.ResponseWriter, r *http.Request) {
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
		DELETE FROM motors
		WHERE motor_id = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete motor: "+err.Error(), http.StatusInternalServerError)
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
		"message": "Motor deleted successfully",
	})
}

