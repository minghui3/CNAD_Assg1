package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"billing-files/models"
	"billing-files/services"
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