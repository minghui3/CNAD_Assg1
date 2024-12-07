package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vehicle-service/models"
	"vehicle-service/services"
	"vehicle-service/database"
	"github.com/gorilla/mux"
)

// GetAllReservationsHandler retrieves all reservations
func GetReservationsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Error converting userID: %v", err)
		return
	}
	reservations, err := services.GetAllReservationsByID(userIDInt)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reservations)
}

// CreateReservationHandler handles creating a new reservation
func CreateReservationsHandler(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var newReservation models.Reservation
	err := json.NewDecoder(r.Body).Decode(&newReservation)
	if err != nil {
		fmt.Printf("Error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the vehicle is already reserved
	isReserved, err := database.CheckVehicleReservation(newReservation.VehicleID, newReservation.StartTime, newReservation.EndTime)
	if err != nil {
		http.Error(w, "Error checking vehicle reservation", http.StatusInternalServerError)
		return
	}

	if isReserved {
		http.Error(w, "Vehicle is already reserved for the specified time", http.StatusConflict)
		return
	}

		// Check if the vehicle is already reserved
	TimingClash, err := database.CheckUserReservation(newReservation.UserID, newReservation.StartTime, newReservation.EndTime)
	if err != nil {
		http.Error(w, "Error checking vehicle reservation", http.StatusInternalServerError)
		return
	}

	if TimingClash {
		http.Error(w, "Vehicle is already reserved for the specified time", http.StatusConflict)
		return
	}

	// Call the service folder to create the reservation
	createdReservation, err := services.CreateReservation(newReservation)
	if err != nil {
		log.Printf("Error creating reservation: %v", err)
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdReservation)
}

// UpdateReservationHandler updates an existing reservation
func UpdateReservationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Error converting userID: %v", err)
		return
	}

	var updatedReservation models.Reservation
	err = json.NewDecoder(r.Body).Decode(&updatedReservation)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the vehicle is already reserved
	isReserved, err := database.CheckVehicleReservation(updatedReservation.VehicleID, updatedReservation.StartTime, updatedReservation.EndTime)
	if err != nil {
		http.Error(w, "Error checking vehicle reservation", http.StatusInternalServerError)
		return
	}

	if isReserved {
		http.Error(w, "Vehicle is already reserved for the specified time", http.StatusConflict)
		return
	}

	// Check if the vehicle is already reserved
	TimingClash, err := database.CheckUserReservation(updatedReservation.UserID, updatedReservation.StartTime, updatedReservation.EndTime)
	if err != nil {
		http.Error(w, "Error checking vehicle reservation", http.StatusInternalServerError)
		return
	}

	if TimingClash {
		http.Error(w, "Vehicle is already reserved for the specified time", http.StatusConflict)
		return
	}

	err = services.UpdateReservation(userIDInt, updatedReservation)
	if err != nil {
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Reservation updated")
}

// DeleteReservationHandler deletes a reservation
func DeleteReservationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Error converting userID: %v", err)
		return
	}

	type ReservationIDRequest struct {
		ReservationID int `json:"reservation_id"`
	}
	var req ReservationIDRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = services.DeleteReservation(userIDInt, req.ReservationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("Reservation deleted")
}
