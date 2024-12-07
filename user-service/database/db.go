package database

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"user-service/models"

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

// GetUserByEmail retrieves a user by email.
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow("SELECT user_id, name, email, phone_number, membership_id, password, verification_code FROM Users WHERE email = ?", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Membership, &user.PasswordHash, &user.VerificationCode)
	if err != nil {
		// Check for "no rows" error
		if err == sql.ErrNoRows {
			// No user found, return nil
			return nil, nil
		}
		// Some other error occurred
		return nil, err
	}
	return &user, nil
}

// GetUserByPhoneNumber retrieves a user by phone number.
func GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow("SELECT user_id, name, email, phone_number, membership_id, password FROM Users WHERE phone_number = ?", phoneNumber).
		Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Membership, &user.PasswordHash)
	if err != nil {
		// Check for "no rows" error
		if err == sql.ErrNoRows {
			// No user found, return nil
			return nil, nil
		}
		// Some other error occurred
		return nil, err
	}
	return &user, nil
}

// InsertUser inserts a new user into the database
func InsertUser(user models.User) error {
	query := "INSERT INTO Users (name, email, phone_number, password, verification_code) VALUES (?, ?, ?, ?, ?)"
	_, err := DB.Exec(query, user.Name, user.Email, user.PhoneNumber, user.PasswordHash, user.VerificationCode)
	return err
}

func UpdateUserVerifiedStatus(email string, verified bool) error {
	query := "UPDATE Users SET verified = ? WHERE email = ?"
	_, err := DB.Exec(query, verified, email) // Assuming `db` is your database connection
	return err
}

func GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := DB.QueryRow("SELECT user_id, name, email, phone_number, membership_id, password FROM Users WHERE user_id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Membership, &user.PasswordHash)
	if err != nil {
		// Check for "no rows" error
		if err == sql.ErrNoRows {
			// No user found, return nil
			return nil, ErrUserNotFound
		}
		// Some other error occurred
		return nil, err
	}
	return &user, nil
}

func UpdateUserByID(user_id int, name string, email string, phone_number string) error {
	query := "UPDATE Users SET name = ?, email = ?, phone_number = ? WHERE user_id = ?"
	_, err := DB.Exec(query, name, email, phone_number, user_id)
	return err
}

// GetRentalHistoryByUserID retrieves the rental history for a specific user from the database
func GetRentalHistoryByUserID(userID int) ([]models.RentalHistory, error) {
	query := "SELECT rental_id, vehicle_id, model, rental_start, rental_end, total_amount FROM RentalHistory WHERE user_id = ?"

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	log.Printf("Fetching rental history for userID: %d", userID)
	var (
		rentalID    string
		vehicleID   int
		model       string
		rentalStart string
		rentalEnd   string
		totalAmount float64
	)
	
	var rentalHistory []models.RentalHistory
	for rows.Next() {
		err := rows.Scan(&rentalID, &vehicleID, &model, &rentalStart, &rentalEnd, &totalAmount)
		if err != nil {
			log.Println("Query error:", err)
			return nil, err
		}
		// Convert strings to time.Time
		startTime, err := time.Parse("2006-01-02 15:04:05", rentalStart)
		if err != nil {
			log.Printf("Failed to parse rental_start: %v", err)
			return nil, err
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", rentalEnd)
		if err != nil {
			log.Printf("Failed to parse rental_end: %v", err)
			return nil, err
		}

		log.Printf("Row scanned - rental_id: %s, vehicle_id: %d, rental_start: %s, rental_end: %s, total_amount: %.2f", 
			rentalID, vehicleID, rentalStart, rentalEnd, totalAmount)

		rentalHistory = append(rentalHistory, models.RentalHistory{
			RentalID:     rentalID,
			VehicleID:    vehicleID,
			StartTime:    startTime,
			Model:		  model,
			EndTime:      endTime,
			TotalAmount:  totalAmount,
		})
	}
	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		return nil, err
	}

	return rentalHistory, nil
}
