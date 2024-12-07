package services

import (
    "vehicle-service/database"
    "vehicle-service/models"
    "fmt"
)

// FetchAvailableVehicles fetches all vehicles that are currently available for rental
func FetchAvailableVehicles() ([]models.Vehicle, error) {
	vehicles, err := database.GetAvailableVehicles()
	if err != nil {
		fmt.Printf("Error fetching available vehicles: %v", err)
		return nil, err
	}

	// Loop through the available vehicles and update reservation status for each
	for _, vehicle := range vehicles {
		err := database.UpdateReservationStatusIfNeeded(vehicle.VehicleID)
		if err != nil {
			fmt.Printf("Error updating reservation status for vehicle_id %d: %v", vehicle.VehicleID, err)
			// Decide if you want to continue processing other vehicles even if an error occurs
			// If not, you can `return nil, err`
		}
	}

	return vehicles, nil
}