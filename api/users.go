package api

import (
	"encoding/json"
	"log"
	"net/http"

	"dreamz.com/api/db"
	"dreamz.com/api/model"
	"github.com/gin-gonic/gin"
)

func (server *Server) getUsers(ctx *gin.Context) {
	users := db.GetUsers(server.store)
	ctx.JSON(http.StatusOK, users)
}

func (server *Server) createUser(c *gin.Context) {
	jsonBody, err := c.GetRawData()
	if err != nil {
		log.Panic(err)
	}
	var user model.User
	json.Unmarshal(jsonBody, &user)
	nUser := db.UpdateUser(server.store, &user)
	c.JSON(http.StatusCreated, nUser)
}
