package services

import (
	"errors"
	"fmt"
	"regexp"
	"user-service/database"
	"user-service/models"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// LoginUser checks if input is email or phone number and handles it.
func LoginUser(input, password string) (*models.User, error) {
	// Define regex patterns
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{7,15}$`)

	var user *models.User
	var err error

	if emailRegex.MatchString(input) {
		// Handle email login/registration
		user, err = database.GetUserByEmail(input)
	} else if phoneRegex.MatchString(input) {
		// Handle phone login/registration
		user, err = database.GetUserByPhoneNumber(input)
	} else {
		return nil, errors.New("invalid email or phone number")
	}

	if err != nil {
		return nil, err // User not found or other database error
	}

	// Debugging: Log the stored password hash and the input password
	fmt.Println("Stored password hash:", user.PasswordHash)
	fmt.Println("Provided password:", password)
    fmt.Println("User email:", user.Email)
	fmt.Println("Password (Raw):", password)
	fmt.Println("Password Hash (Raw):", user.PasswordHash)

	// Verify the password using the VerifyPassword function
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    fmt.Println(err)
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	return user, nil
}

// RegisterUser registers a new user with a hashed password
func RegisterUser(user models.User) error {
	// Check if user exists by email
	existingUser, err := database.GetUserByEmail(user.Email)
	if err != nil {
		return errors.New("failed to check if user exists by email")
	}
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Check if user exists by phone number
	existingPhoneUser, err := database.GetUserByPhoneNumber(user.PhoneNumber)
	if err != nil {
		return errors.New("failed to check if user exists by phone number")
	}
	if existingPhoneUser != nil {
		return errors.New("user with this phone number already exists")
	}
    fmt.Println(user.PasswordHash)
	// Hash the password before saving the user
	hashedPassword, err := HashPassword(user.PasswordHash)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.PasswordHash = hashedPassword

	// Insert user into database
	err = database.InsertUser(user)
	if err != nil {
		return errors.New("failed to insert user into the database")
	}

	return nil
}
