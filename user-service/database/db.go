package database

import (
    "database/sql"
    "log"
	_ "github.com/go-sql-driver/mysql"
	"user-service/models"
	"errors"
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
