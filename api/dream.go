package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"dreamz.com/api/db"
	"dreamz.com/api/model"
	"github.com/gin-gonic/gin"
)

func (server *Server) getDreams(c *gin.Context) {
	userId := c.Keys["userId"].(string)
	dreamDays := db.GetDreamDays(server.store, userId)
	c.JSON(http.StatusOK, gin.H{"results": dreamDays})
}

func (server *Server) updateDream(c *gin.Context) {
	jsonBody, err := c.GetRawData()
	if err != nil {
		log.Fatal(err)
	}
	var dream model.DreamDay
	json.Unmarshal(jsonBody, &dream)

	fmt.Printf("%#v", dream) //with type

	userId := c.Keys["userId"].(string)
	dream.UserId = userId

	updated := db.UpdateDreamDay(server.store, &dream)
	if updated == nil {
		c.String(http.StatusNotFound, "")
	} else {
		c.JSON(http.StatusOK, dream)
	}

}
