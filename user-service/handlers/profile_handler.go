package handlers

import (
    "encoding/json"
    "net/http"
    "user-service/models"
    "user-service/services"
    "strconv"
    "log"
)

// UpdateProfileHandler handles updating the user's profile
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    updatedUser, err := services.UpdateProfile(user.ID, user.Name, user.Membership)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(updatedUser)
}

// GetProfileHandler handles fetching the user's profile
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
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

    json.NewEncoder(w).Encode(user)
}
