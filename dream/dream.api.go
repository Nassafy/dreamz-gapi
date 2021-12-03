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

func AddDreamRoute(r *gin.Engine, store *common.Store) {
	s := dreamServer{store: store}
	dreamRouter := r.Group("/dream", auth.AuthMiddleware())
	dreamRouter.GET("", auth.AuthMiddleware(), s.getDreams)
	dreamRouter.GET(":id", auth.AuthMiddleware(), s.getTodayDream)
	dreamRouter.POST("", auth.AuthMiddleware(), s.updateDream)
	dreamRouter.DELETE("", auth.AuthMiddleware(), s.deleteDream)
}

func (s *dreamServer) getDreams(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		noUserIdError(c)
		return
	}
	dreamDays := dbGetDreamDays(s.store, userId)
	c.JSON(http.StatusOK, gin.H{"results": dreamDays})
}

func (s *dreamServer) getTodayDream(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		noUserIdError(c)
		return
	}
	dreamdDay := dbGetTodayDream(s.store, userId)
	if dreamdDay == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Today dream not found"})
	} else {
		c.JSON(http.StatusOK, dreamdDay)
	}
}

func (s *dreamServer) updateDream(c *gin.Context) {
	var dream DreamDay
	err := c.ShouldBind(&dream)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	userId, err := GetUserId(c)
	if err != nil {
		noUserIdError(c)
		return
	}
	dream.UserId = userId

	updated := dbUpdateDreamDay(s.store, &dream)
	if updated == nil {
		c.String(http.StatusNotFound, "Dream not found")
	} else {
		c.JSON(http.StatusOK, dream)
	}
}

func (s *dreamServer) deleteDream(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		noUserIdError(c)
	}
	dreamId := c.Param("id")
	DbDeleteDreamDay(s.store, dreamId, userId)
	c.Status(http.StatusNoContent)
}

func noUserIdError(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "No user ID"})
}

func GetUserId(c *gin.Context) (string, error) {
	userId := c.Keys["userId"]
	if userId == nil {
		return "", errors.New("no user id")
	}
	return userId.(string), nil
}
