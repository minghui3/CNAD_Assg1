package main

import (
    "log"
    "net/http"
    "billing-service/handlers"
	"billing-service/database"
	"github.com/gorilla/mux"
)

func main() {
	database.Initialize()

	router := mux.NewRouter()
    router.HandleFunc("/api/v1/calculate", handlers.CalculateAmountHandler).Methods("POST")
	router.HandleFunc("/api/v1/promotion", handlers.ApplyPromotionHandler).Methods("PUT")
	router.HandleFunc("/api/v1/updateamount", handlers.UpdateReservationHandler).Methods("PUT")
	router.HandleFunc("/api/v1/billing", handlers.InsertBillingHandler).Methods("POST")

    // Wrap the router with the CORS middleware
    corsRouter := enableCORS(router)

    log.Println("Billing server is running on port 8083")
    log.Fatal(http.ListenAndServe(":8083", corsRouter))
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
