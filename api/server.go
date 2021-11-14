package api

import (
	"errors"

	"dreamz.com/api/db"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

type Server struct {
	router *gin.Engine
	store  *db.Store
}

func NewServer() *Server {
	server := &Server{}
	router := gin.Default()

	store := db.NewStore()

	router.Use(cors.Default())

	router.POST("auth/login", server.Login)

	//TODO delete, mega dangereux
	router.GET("users", AuthMiddleware(), cors.Default(), server.getUsers)
	router.POST("users", AuthMiddleware(), cors.Default(), server.createUser)

	dreamRoute := router.Group("/dream", AuthMiddleware())
	{

		dreamRoute.GET("", AuthMiddleware(), server.getDreams)
		dreamRoute.GET("/today", AuthMiddleware(), server.getTodayDream)
		dreamRoute.POST("", AuthMiddleware(), server.updateDream)
	}

	server.router = router
	server.store = store
	return server
}

func (server *Server) Start(adress string) error {
	return server.router.Run(adress)
}

func (server *Server) CloseStore() {
	server.store.CloseStore()
}

func GetUserId(c *gin.Context) (string, error) {
	userId := c.Keys["userId"]
	if userId == nil {
		return "", errors.New("no user id")
	}
	return userId.(string), nil
}
