package api

import (
	"dreamz.com/api/db"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	store  *db.Store
}

func NewServer() *Server {
	server := &Server{}
	router := gin.Default()

	store := db.NewStore()

	router.POST("auth/login", server.Login)
	router.GET("users", AuthMiddleware(), server.getUsers)

	dreamRoute := router.Group("/dream", AuthMiddleware())
	{

		dreamRoute.GET("", AuthMiddleware(), server.getDreams)
		dreamRoute.POST("", AuthMiddleware(), server.newDream)
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
