package handlers

import (
    "encoding/json"
    "net/http"

    "vehicle-service/services"
)

// ViewAvailableVehicles handles requests to fetch available vehicles.
func ViewAvailableVehicles(w http.ResponseWriter, r *http.Request) {
    vehicles, err := services.FetchAvailableVehicles()
    if err != nil {
        http.Error(w, "Failed to fetch vehicles", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(vehicles)
}
