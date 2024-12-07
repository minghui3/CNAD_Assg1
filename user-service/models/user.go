package models

import "time"

type User struct {
	ID               int    `json:"user_id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	Membership       string `json:"membership_tier"`
	PasswordHash     string `json:"password"` // Excluded from JSON for security
	VerificationCode string `json:"verification_code"`
}

type RentalHistory struct {
	RentalID    string `json:"rental_id"`
	UserID      int `json:"user_id"`
	VehicleID   int `json:"vehicle_id"`
	Model       string `json:"model"`
	StartTime   time.Time `json:"rental_start"`
	EndTime     time.Time `json:"rental_end"`
	TotalAmount float64   `json:"total_cost"`
}
