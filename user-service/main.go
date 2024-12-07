package main

import (
    "log"
    "net/http"
    "user-service/handlers"
	"user-service/database"
	"github.com/gorilla/mux"
)

func main() {
	database.Initialize()

	router := mux.NewRouter()
    router.HandleFunc("/api/v1/register", handlers.RegisterHandler).Methods("POST")
    router.HandleFunc("/api/v1/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/verify", handlers.VerifyCodeHandler).Methods("POST")
    router.HandleFunc("/api/v1/update-profile/{user_id}", handlers.UpdateProfileHandler).Methods("PUT")
    router.HandleFunc("/api/v1/get-profile/{user_id}", handlers.GetProfileHandler).Methods("GET")

    // Wrap the router with the CORS middleware
    corsRouter := enableCORS(router)

    log.Println("Server is running on port 8081")
    log.Fatal(http.ListenAndServe(":8081", corsRouter))
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
