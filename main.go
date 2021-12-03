package main

import (
	"log"
	"os"

	"dreamz.com/api/auth"
	"dreamz.com/api/common"
	"dreamz.com/api/dream"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No valid .env file")
	}

	port := os.Getenv("PORT")
	serverAdress := "0.0.0.0:" + port

	s := common.NewStore()
	r := gin.Default()
	auth.AddAuthRoute(r, s)
	dream.AddDreamRoute(r, s)
	r.Run(serverAdress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
