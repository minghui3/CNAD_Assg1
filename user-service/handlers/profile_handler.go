package handlers

import (
    "encoding/json"
    "net/http"
    "user-service/models"
    "user-service/services"
    "strconv"
    "log"
	"github.com/gorilla/mux"
)

// UpdateProfileHandler handles updating the user's profile
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    params := mux.Vars(r)
	userID := params["user_id"]
    userIDInt, err := strconv.Atoi(userID)
	if err != nil {
        log.Printf("Error converting userID: %v", err)
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    updatedUser, err := services.UpdateProfile(userIDInt, user.Name, user.Email, user.PhoneNumber)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedUser)
}

// GetProfileHandler handles fetching the user's profile
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
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
    user, err := services.GetProfile(userIDInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// GetRentalHistoryHandler handles fetching the rental history of a user
func GetRentalHistoryHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    userID := params["user_id"]
    if userID == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    userIDInt, err := strconv.Atoi(userID)
    if err != nil {
        log.Printf("Error converting userID: %v", err)
        http.Error(w, "Invalid User ID", http.StatusBadRequest)
        return
    }

    rentalHistory, err := services.GetRentalHistory(userIDInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rentalHistory)
}
