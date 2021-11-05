package api

import (
	"encoding/json"
	"log"
	"net/http"

	"dreamz.com/api/db"
	"dreamz.com/api/model"
	"github.com/gin-gonic/gin"
)

func (server *Server) getDreams(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		noUserIdError(c)
		return
	}
	dreamDays := db.GetDreamDays(server.store, userId)
	c.JSON(http.StatusOK, gin.H{"results": dreamDays})
}

func (server *Server) getTodayDream(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		noUserIdError(c)
		return
	}
	dreamdDay := db.GetTodayDream(server.store, userId)
	if dreamdDay == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Today dream not found"})
	} else {
		c.JSON(http.StatusOK, dreamdDay)
	}
}

func (server *Server) updateDream(c *gin.Context) {
	jsonBody, err := c.GetRawData()
	if err != nil {
		log.Panic(err)
	}
	var dream model.DreamDay
	json.Unmarshal(jsonBody, &dream)

	userId, err := GetUserId(c)
	if err != nil {
		noUserIdError(c)
		return
	}
	dream.UserId = userId

	updated := db.UpdateDreamDay(server.store, &dream)
	if updated == nil {
		c.String(http.StatusNotFound, "Dream not found")
	} else {
		c.JSON(http.StatusOK, dream)
	}
}

func noUserIdError(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "No user ID"})
}
