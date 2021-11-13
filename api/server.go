package api

import (
	"errors"

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

	router.Use(CORSMiddleware())

	router.POST("auth/login", server.Login)

	//TODO delete, mega dangereux
	router.GET("users", AuthMiddleware(), server.getUsers)
	router.POST("users", AuthMiddleware(), server.createUser)

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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func GetUserId(c *gin.Context) (string, error) {
	userId := c.Keys["userId"]
	if userId == nil {
		return "", errors.New("no user id")
	}
	return userId.(string), nil
}
