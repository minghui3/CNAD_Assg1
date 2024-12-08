package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
	"vehicle-service/models"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// ErrUserNotFound is returned when a user cannot be found in the database
var ErrUserNotFound = errors.New("user not found")

// Initialize initializes the MySQL database connection
func Initialize() {
	var err error
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/CarSharingSystem")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Ping the database to verify the connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	DB = db

	log.Println("Connected to the MySQL database successfully!")
}

// GetAvailableVehicles retrieves all available vehicles from the database.
func GetAvailableVehicles() ([]models.Vehicle, error) {
	query := "select vehicle_id, model, location from vehicles"

	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var vehicle models.Vehicle
		err := rows.Scan(&vehicle.VehicleID, &vehicle.Model, &vehicle.Location)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Dynamically check reservation status
		err = UpdateReservationStatusIfNeeded(vehicle.VehicleID)
		if err != nil {
			log.Printf("Failed to update reservation status: %v", err)
			continue
		}

		vehicles = append(vehicles, vehicle)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		return nil, err
	}

	return vehicles, nil
}

// GetAvailableVehicles retrieves all available vehicles from the database.
func GetReservationsByID(userID int) ([]models.Reservation, error) {
	query := "SELECT reservation_id, user_id, vehicle_id, start_time, end_time, status, total_amount FROM Reservations WHERE user_id = ?"
	rows, err := DB.Query(query, userID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var (
		reservationID int
		userIDTemp    int
		vehicleID     int
		startTimeStr  string
		endTimeStr    string
		status        string
		totalAmount   float64
	)

	var reservations []models.Reservation
	for rows.Next() {
		err := rows.Scan(&reservationID, &userIDTemp, &vehicleID, &startTimeStr, &endTimeStr, &status, &totalAmount)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		// Convert strings to time.Time
		startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			log.Printf("Failed to parse start_time: %v", err)
			return nil, err
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			log.Printf("Failed to parse end_time: %v", err)
			return nil, err
		}

		// Append to the list of reservations
		reservations = append(reservations, models.Reservation{
			ReservationID: reservationID,
			UserID:        userIDTemp,
			VehicleID:     vehicleID,
			StartTime:     startTime,
			EndTime:       endTime,
			Status:        status,
			TotalAmount:   totalAmount,
		})
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		return nil, err
	}

	return reservations, nil
}

// InsertUser inserts a new reservation into the database
func InsertReservation(reservation models.Reservation) (*models.Reservation, error) {
	query := "INSERT INTO Reservations (user_id, vehicle_id, start_time, end_time, total_amount) VALUES (?, ?, ?, ?, ?)"
	result, err := DB.Exec(query, reservation.UserID, reservation.VehicleID, reservation.StartTime, reservation.EndTime, reservation.TotalAmount)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("no reservation/user found with the given ID")
	}

	// Get the auto-incremented ID
	reservationID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error fetching inserted ID: %v", err)
		return nil, err
	}

	// Set the fetched reservation ID to the model
	reservation.ReservationID = int(reservationID)
	return &reservation, err
}

// UpdateReservationsByID updates an existing in the database
func UpdateReservationsByID(user_id int, updated models.Reservation) error {
	query := "UPDATE Reservations SET vehicle_id = ?, start_time = ?, end_time = ?, total_amount = ? WHERE reservation_id = ? AND user_id = ?"
	result, err := DB.Exec(query, updated.VehicleID, updated.StartTime, updated.EndTime, updated.TotalAmount, updated.ReservationID, user_id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no reservation/user found with the given ID")
	}
	return err
}

func DeleteReservationByID(user_id int, rev_id int) error {
	query := "DELETE FROM Reservations WHERE reservation_id = ? AND user_id = ? AND (start_time < NOW() - INTERVAL 7 DAY OR start_time > NOW() + INTERVAL 7 DAY)"
	result, err := DB.Exec(query, rev_id, user_id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no reservation/user found with the given ID or user tried to delete within 7 days")
	}
	return err
}

// CheckVehicleReservation checks if the vehicle already has a reservation with overlapping time range
func CheckVehicleReservation(vehicleID int, startTime time.Time, endTime time.Time) (bool, error) {
	query := `SELECT COUNT(*) FROM Reservations WHERE vehicle_id = ? AND start_time > ? AND end_time < ? AND status = 'active'`

	// Execute the query
	var count int
	err := DB.QueryRow(query, vehicleID, endTime, startTime).Scan(&count)
	if err != nil {
		log.Printf("Error executing vehicle reservation check: %v", err)
		return false, err
	}

	if count > 0 {
		// Vehicle is already reserved
		return true, nil
	}

	// No existing reservation conflicts
	return false, nil
}

// CheckUserReservation checks if the user already has a reservation with overlapping time range
func CheckUserReservation(userID int, startTime time.Time, endTime time.Time) (bool, error) {
	query := `SELECT COUNT(*) FROM Reservations WHERE user_id = ? AND start_time > ? AND end_time < ? AND status = 'active'`

	// Execute the query
	var count int
	err := DB.QueryRow(query, userID, endTime, startTime).Scan(&count)
	if err != nil {
		log.Printf("Error executing user reservation check: %v", err)
		return false, err
	}

	if count > 0 {
		// user has existing reservation on that time
		return true, nil
	}

	// No existing reservation conflicts
	return false, nil
}

// UpdateReservationStatusIfNeeded dynamically updates reservation status based on the current time
func UpdateReservationStatusIfNeeded(vehicleID int) error {
	currentTime := time.Now()

	// Update status to "active" if the current time is within the reservation range
	activeQuery := `
		UPDATE Reservations 
		SET status = 'active' 
		WHERE vehicle_id = ? AND start_time <= ? AND end_time >= ? AND status != 'active'
	`
	_, err := DB.Exec(activeQuery, vehicleID, currentTime, currentTime)
	if err != nil {
		fmt.Printf("Error setting active status: %v", err)
		return err
	}

	// Update status to "completed" if the reservation has ended
	completedQuery := `
		UPDATE Reservations 
		SET status = 'completed' 
		WHERE vehicle_id = ? AND end_time < ? AND status != 'completed'
	`
	res, err := DB.Exec(completedQuery, vehicleID, currentTime)
	if err != nil {
		fmt.Printf("Error setting completed status: %v", err)
		return err
	}

	// Insert into rental history if any rows were updated (i.e., reservation was marked as completed)
	affectedRows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error fetching affected rows: %v", err)
		return err
	}

	if affectedRows > 0 {
		// Fetch the completed reservation details to log into rental history
		var reservation = &models.Reservation{}
		var startTimeStr, endTimeStr string // Temporary variables for scanning datetime as string

		selectQuery := `
			SELECT reservation_id, user_id, start_time, end_time, total_amount
			FROM Reservations
			WHERE vehicle_id = ? AND end_time < ? AND status = 'completed'
			LIMIT 1
		`
		row := DB.QueryRow(selectQuery, vehicleID, currentTime)
		err = row.Scan(&reservation.ReservationID, &reservation.UserID, &startTimeStr, &endTimeStr, &reservation.TotalAmount)
		if err != nil {
			fmt.Printf("Failed to fetch reservation details for history: %v", err)
			return err
		}

		// Parse start_time and end_time strings into time.Time
		reservation.StartTime, err = time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			fmt.Printf("Failed to parse start_time: %v", err)
			return err
		}

		reservation.EndTime, err = time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			fmt.Printf("Failed to parse end_time: %v", err)
			return err
		}

		// Count the number of rows in the RentalHistory table to determine the new rental_id
		var maxRentalID string
		countQuery := `SELECT MAX(rental_id) FROM RentalHistory`
		err = DB.QueryRow(countQuery).Scan(&maxRentalID)
		if err != nil {
			fmt.Printf("Failed to retrieve max rental_id: %v\n", err)
			return err
		}

		// Extract the numeric part from the rental_id
		var currentMaxID int
		if maxRentalID != "" {
			currentMaxID, err = strconv.Atoi(maxRentalID[1:]) // Skip the first character ('R')
			if err != nil {
				fmt.Printf("Failed to parse max rental_id: %v\n", err)
				return err
			}
		} else {
			// If no rental_id exists, start from 0
			currentMaxID = 0
		}

		// Generate the new rental_id in the format Rxxx
		newRentalID := fmt.Sprintf("R%03d", currentMaxID+1)

		// Create another query to get model
		var modelType string
		modelQuery := `SELECT Model FROM vehicles WHERE vehicle_id = ?`
		err = DB.QueryRow(modelQuery, vehicleID).Scan(&modelType)
		if err != nil {
			fmt.Printf("Failed to count rows in RentalHistory: %v", err)
			return err
		}
		// Insert the fetched details into rental history
		insertQuery := `
			INSERT INTO RentalHistory (rental_id, user_id, vehicle_id, model, rental_start, rental_end, total_amount)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err = DB.Exec(insertQuery, newRentalID, reservation.UserID, vehicleID, modelType, reservation.StartTime, reservation.EndTime, reservation.TotalAmount)
		if err != nil {
			fmt.Printf("Error inserting into rental history: %v", err)
			return err
		}
	}

	return nil
}
