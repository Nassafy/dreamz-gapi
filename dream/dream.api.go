package dream

import (
	"errors"
	"net/http"

	"dreamz.com/api/auth"
	"dreamz.com/api/common"

	"github.com/gin-gonic/gin"
)

type dreamServer struct {
	store *common.Store
}

// AddDreamRoute add the dream route to the main router
func AddDreamRoute(r *gin.Engine, store *common.Store) {
	s := dreamServer{store: store}
	dreamRouter := r.Group("/dream", auth.Middleware())
	dreamRouter.GET("", s.getDreams)
	dreamRouter.GET("/today", s.getTodayDream)
	dreamRouter.GET(":id", s.getDreamDay)
	dreamRouter.POST("", s.updateDream)
	dreamRouter.DELETE(":id", s.deleteDream)
}

func (s *dreamServer) getDreams(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		noUserIDError(c)
		return
	}
	dreamDays := dbGetDreamDays(s.store, userID)
	c.JSON(http.StatusOK, gin.H{"results": dreamDays})
}

func (s *dreamServer) getTodayDream(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		noUserIDError(c)
		return
	}
	dreamsDay := dbGetTodayDream(s.store, userID)
	if dreamsDay == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Today dream not found"})
	} else {
		c.JSON(http.StatusOK, dreamsDay)
	}
}

func (s *dreamServer) getDreamDay(c *gin.Context) {
	dayID := c.Param("id")
	userID, err := getUserID(c)
	if err != nil {
		noUserIDError(c)
		return
	}
	dreamDay := dbGetDreamDay(s.store, userID, dayID)
	if dreamDay == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dream not found"})
	} else {
		c.JSON(http.StatusOK, dreamDay)
	}
}

func (s *dreamServer) updateDream(c *gin.Context) {
	var dream Day
	err := c.ShouldBind(&dream)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	userID, err := getUserID(c)
	if err != nil {
		noUserIDError(c)
		return
	}
	dream.UserID = userID

	updated := dbUpdateDreamDay(s.store, &dream)
	if updated == nil {
		c.String(http.StatusNotFound, "Dream not found")
	} else {
		c.JSON(http.StatusOK, dream)
	}
}

func (s *dreamServer) deleteDream(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		noUserIDError(c)
	}
	dreamID := c.Param("id")
	DbDeleteDreamDay(s.store, dreamID, userID)
	c.Status(http.StatusNoContent)
}

func noUserIDError(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "No user ID"})
}

func getUserID(c *gin.Context) (string, error) {
	userID := c.Keys["userID"]
	if userID == nil {
		return "", errors.New("no user id")
	}
	return userID.(string), nil
}
