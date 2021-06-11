package api

import (
	"net/http"

	"dreamz.com/api/db"
	"github.com/gin-gonic/gin"
)

func (server *Server) getUsers(ctx *gin.Context) {
	users := db.GetUsers(server.store)
	ctx.JSON(http.StatusOK, users)
}
