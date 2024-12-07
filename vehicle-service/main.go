package main

import (
    "log"
    "net/http"
    "vehicle-service/handlers"
	"vehicle-service/database"
	"github.com/gorilla/mux"
)

func main() {
	database.Initialize()
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/vehicles", handlers.ViewAvailableVehicles).Methods("GET")
	router.HandleFunc("/api/v1/reservations/{user_id}", handlers.GetReservationsHandler).Methods("GET")
	router.HandleFunc("/api/v1/reservations", handlers.CreateReservationsHandler).Methods("POST")
	router.HandleFunc("/api/v1/reservations/{user_id}", handlers.UpdateReservationHandler).Methods("PUT")
	router.HandleFunc("/api/v1/reservations/{user_id}", handlers.DeleteReservationHandler).Methods("DELETE")

    // Wrap the router with the CORS middleware
    corsRouter := enableCORS(router)

    log.Println("Vehicle server is running on port 8082")
    log.Fatal(http.ListenAndServe(":8082", corsRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
