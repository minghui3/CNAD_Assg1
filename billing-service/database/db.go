package database

import (
	"billing-service/models"
	"database/sql"
	"errors"
	"fmt"
	"log"

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

// GetDiscountPercentage checks the promotion code in the database and returns the discount percentage
func GetDiscountPercentage(promotionCode string) (float64, error) {
	var discountPercentage float64

	query := `SELECT discount FROM promotions WHERE promotion_code = ?`
	err := DB.QueryRow(query, promotionCode).Scan(&discountPercentage)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Promotion code not found
		}
		return 0, err
	}

	return discountPercentage, nil
}

// UpdateBillingDiscount updates the discount in the billing table for a specific reservation ID
func UpdateBillingDiscount(reservationID int, discountPercentage float64) error {
	query := `UPDATE billing SET discount = ? WHERE reservation_id = ?`
	_, err := DB.Exec(query, discountPercentage, reservationID)
	return err
}

// InsertBilling creates a new entry
func InsertBilling(reservationID int, amount float64) error {
	query := `INSERT INTO billing (reservation_id, initial_amount) VALUE ( ? , ? )`
	_, err := DB.Exec(query, reservationID, amount)
	return err
}

// UpdateBillingDiscount does the final update of amount to the reservation table for a specific reservation ID
func UpdateReservationAmount(revId int) error {

	// Get final amount from billing table
	var finalAmount float64
	query := `SELECT final_amount FROM billing WHERE reservation_id = ?`
	err := DB.QueryRow(query, revId).Scan(&finalAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no billing record found for reservation ID: %d", revId)
		}
		return fmt.Errorf("failed to query billing table: %v", err)
	}

	// Update the reservation table with the selected final_amount
	updateQuery := `UPDATE reservations SET total_amount = ? WHERE reservation_id = ?`
	_, err = DB.Exec(updateQuery, finalAmount, revId)
	if err != nil {
		return fmt.Errorf("failed to update reservation table: %v", err)
	}

	return nil
}

// CheckBillingExists checks if a billing entry exists for a given reservation ID
func CheckBillingExists(reservationID int) (bool, error) {
	var exists bool

	// Query to check if a record exists
	query := `SELECT EXISTS(SELECT 1 FROM billing WHERE reservation_id = ?)`
	err := DB.QueryRow(query, reservationID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to query billing table: %v", err)
	}

	return exists, nil
}

func GetReservationForEmail(reservationID int) (models.EmailItems, error) {

	var emailItems models.EmailItems
	// Query to check if a record exists
	query := `SELECT initial_amount, discount, final_amount FROM Billing WHERE reservation_id = ?`
	err := DB.QueryRow(query, reservationID).Scan(&emailItems.InitialAmount, &emailItems.Discount, &emailItems.FinalAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found
			log.Printf("No billing data found for reservation ID %d", reservationID)
			return models.EmailItems{}, nil
		}
		// Log the error and return it
		log.Printf("Error executing query: %v", err)
		return models.EmailItems{}, err
	}

	return emailItems, nil
}
