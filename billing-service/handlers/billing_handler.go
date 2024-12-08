package handlers

import (
	"billing-service/services"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CalculateAmountHandler calculates the amount based on time given (0.30 per min)
func CalculateAmountHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	// Parse and decode the JSON payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert starttime and endtime to time.Time
	startTime, err := time.Parse(time.RFC3339, payload.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, payload.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	// Call the service function
	amount, err := services.CalculateAmount(startTime, endTime)
	if err != nil {
		http.Error(w, "Error calculating amount: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"amount": amount})
}

// ApplyPromotionHandler applies promotion
func ApplyPromotionHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RevID         int `json:"reservation_id"`
		PromotionCode string `json:"promotion_code"`
	}

	// Parse and decode the JSON payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service function
	discountedAmount, err := services.ApplyPromotion(payload.RevID, payload.PromotionCode)
	if err != nil {
		http.Error(w, "Error applying promotion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Discount Applied!", "discountedAmount": discountedAmount})
}

// UpdateReservationHandler handles the HTTP endpoint for updating reservation amount
func UpdateReservationHandler(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		RevID int `json:"reservation_id"`
	}
	// Parse and decode the JSON payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service/database logic to perform the update
	err := services.UpdateReservationAmount(payload.RevID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update reservation amount: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation amount updated successfully"})
}

// InsertBillingHandler handles the HTTP endpoint for inserting new bill
func InsertBillingHandler(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		RevID  int  `json:"reservation_id"`
		Amount float64 `json:"total_amount"`
	}
	// Parse and decode the JSON payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service/database logic to perform the update
	err := services.InsertBilling(payload.RevID,payload.Amount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update reservation amount: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Billing inserted successfully"})
}
