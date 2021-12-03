package user

import (
	"encoding/json"
	"log"
	"net/http"

	"dreamz.com/api/common"
	"github.com/gin-gonic/gin"
)

type userServer struct {
	store *common.Store
}

func (server *userServer) getUsers(ctx *gin.Context) {
	users := dbGetUsers(server.store)
	ctx.JSON(http.StatusOK, users)
}

func (server *userServer) createUser(c *gin.Context) {
	jsonBody, err := c.GetRawData()
	if err != nil {
		log.Fatal(err)
	}
	var user User
	json.Unmarshal(jsonBody, &user)
	nUser := dbUpdateUser(server.store, &user)
	c.JSON(http.StatusCreated, nUser)
}
