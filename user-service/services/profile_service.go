package services

import (
	"errors"
	"user-service/database"
	"user-service/models"
    "fmt"
)

// UpdateProfile updates a user's profile with new information
func UpdateProfile(userID int, name string, email string, phone_number string) (*models.User, error) {
    err := database.UpdateUserByID(userID, name, email, phone_number)
    if err != nil {
        return nil, fmt.Errorf("failed to update user with ID %d: %w", userID, err)
    }
    updatedUser, err := database.GetUserByID(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch updated user with ID %d: %w", userID, err)
    }
    return updatedUser, nil
}

// GetProfile fetches a user's profile information
func GetProfile(userID int) (*models.User, error) {
    var user *models.User
	var err error
    user, err = database.GetUserByID(userID)
    if err != nil {
        if errors.Is(err, database.ErrUserNotFound) { // Check for error in the case that user is not found
            return nil, fmt.Errorf("user with ID %d not found", userID)
        }
        return nil, fmt.Errorf("failed to fetch user profile: %w", err) // Wrap the error
    }
    return user, nil
}

// GetRentalHistory retrieves the rental history for a specific user
func GetRentalHistory(userID int) ([]models.RentalHistory, error) {
    var rentalHistory []models.RentalHistory
    var err error
	// Fetch the rental history from the database
	rentalHistory, err = database.GetRentalHistoryByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to fetch rental history")
	}

	if len(rentalHistory) == 0 {
		return nil, errors.New("no rental history found for the user")
	}

	return rentalHistory, nil
}

