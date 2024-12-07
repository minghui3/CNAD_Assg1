package models

import "time"

type User struct {
    ID           int       `json:"user_id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    PhoneNumber  string    `json:"phone_number"`
    Membership   string    `json:"membership_tier"`
    PasswordHash string    `json:"password"` // Excluded from JSON for security
    VerificationCode string `json:"verification_code"`
}

type RentalHistory struct {
    RentalID  int       `json:"rental_id"`
    UserID    int       `json:"user_id"`
    CarID     int       `json:"car_id"`
    StartDate time.Time `json:"start_date"`
    EndDate   time.Time `json:"end_date"`
    TotalCost float64   `json:"total_cost"`
}