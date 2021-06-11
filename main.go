package main

import (
	"log"
	"os"

	"dreamz.com/api/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No valid .env file")
	}

	port := os.Getenv("PORT")
	serverAdress := "0.0.0.0:" + port

	server := api.NewServer()
	defer server.CloseStore()
	err = server.Start(serverAdress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
