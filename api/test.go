package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) helloWorld(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello world")
}
