package models

import "time"

type Vehicle struct {
	VehicleID int    `json:"vehicle_id"`
	Model     string `json:"model"`
	Location  string `json:"location"`
}

type Reservation struct {
	ReservationID int       `json:"reservation_id"`
	UserID        int       `json:"user_id"`
	VehicleID     int       `json:"vehicle_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Status        string    `json:"status"`
	TotalAmount   float64   `json:"total_amount"`
}