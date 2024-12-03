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
    router.HandleFunc("/api/v1/update-profile", handlers.UpdateProfileHandler)
    router.HandleFunc("/api/v1/get-profile", handlers.GetProfileHandler)

    log.Println("Server is running on port 8081")
    log.Fatal(http.ListenAndServe(":8081", router))
}
