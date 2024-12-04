package database

import (
    "database/sql"
    "log"
	_ "github.com/go-sql-driver/mysql"
	"user-service/models"
)

var DB *sql.DB

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
	err := DB.QueryRow("SELECT id, name, email, phone_number, membership_tier, password FROM Users WHERE email = ?", email).
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

// GetUserByPhoneNumber retrieves a user by phone number.
func GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow("SELECT id, name, email, phone_number, membership_tier, password FROM Users WHERE phone_number = ?", phoneNumber).
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
    query := "INSERT INTO Users (name, email, phone_number, membership_tier, password) VALUES (?, ?, ?, ?, ?)"
    _, err := DB.Exec(query, user.Name, user.Email, user.PhoneNumber, user.Membership, user.PasswordHash)
    return err
}

// UpdateUser allows an existing user to change his credentials
func UpdateUser()