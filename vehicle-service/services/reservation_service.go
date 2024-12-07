package services

import (
	"vehicle-service/models"
	"vehicle-service/database"
)

// GetAllReservations retrieves all reservations from the database
func GetAllReservationsByID(userID int) ([]models.Reservation, error) {
    return database.GetReservationsByID(userID)
}

// CreateReservation adds a new reservation
func CreateReservation(reservation models.Reservation) (*models.Reservation, error) {
	reservation_entry,err := database.InsertReservation(reservation)
	if err != nil {
		return	reservation_entry,nil
	}
	return reservation_entry, nil
}

// UpdateReservation updates an existing reservation
func UpdateReservation(user_id int, updated models.Reservation) error {
	return database.UpdateReservationsByID(user_id,updated)
}

// DeleteReservation removes a reservation by ID
func DeleteReservation(user_id int, rev_id int) error {
	return database.DeleteReservationByID(user_id, rev_id)
}
