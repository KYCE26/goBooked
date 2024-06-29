package main

import (
	"belajar/database"
	"fmt"
	"belajar/controller/auth"
	"belajar/controller/booked"
	"belajar/controller/motor"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()	
	fmt.Println("Hello World")

	router := mux.NewRouter()

	router.HandleFunc("/regis", auth.Registration).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")

	// Router handler motors
	router.HandleFunc("/motors", motor.GetMotor).Methods("GET")
	router.HandleFunc("/motors", auth.JWTAuth(motor.PostMotor)).Methods("POST")
	router.HandleFunc("/motors", auth.JWTAuth(motor.PutMotor)).Methods("PUT")
	router.HandleFunc("/motors/{id}", auth.JWTAuth(motor.DeleteMotor)).Methods("DELETE")

	// Router handler bookeds
	router.HandleFunc("/bookeds", booked.GetBooked).Methods("GET")
	router.HandleFunc("/bookeds", auth.JWTAuth(booked.PostBooked)).Methods("POST")
	router.HandleFunc("/bookeds/{id}", auth.JWTAuth(booked.PutBooked)).Methods("PUT")
	router.HandleFunc("/bookeds/{id}", auth.JWTAuth(booked.DeleteBooked)).Methods("DELETE")

	c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://127.0.0.1:5500"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
        Debug: true,
    })
	
    handler := c.Handler(router)
	
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))


}