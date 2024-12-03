package services

import (
    "errors"
    "user-service/models"
)

// ProfileService manages user profile-related actions
var users = []models.User{} // Simulated in-memory user storage

// UpdateProfile updates a user's profile with new information
func UpdateProfile(userID int, newName, newMembership string) (*models.User, error) {
    for i, u := range users {
        if u.ID == userID {
            // Update user information
            users[i].Name = newName
            users[i].Membership = newMembership
            return &users[i], nil
        }
    }
    return nil, errors.New("user not found")
}

// GetProfile fetches a user's profile information
func GetProfile(userID int) (*models.User, error) {
    for _, u := range users {
        if u.ID == userID {
            return &u, nil
        }
    }
    return nil, errors.New("user not found")
}
