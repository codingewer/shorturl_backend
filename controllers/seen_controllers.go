package controllers

import (
	"net/http"
	"short-link/auth"
	"short-link/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserSeenData(c *gin.Context) {
	days := c.Param("days")
	daysInt, err := strconv.ParseInt(days, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gün verisi geçersiz formatta"})
		return
	}
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	seen := models.Seen{}
	data, err := seen.GetSeenData(tokenUser.ID, int(daysInt))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"balanceChart": data.Balance,
		"viewsChart":   data.Views,
		"data":         data,
		"user":         tokenUser,
		"days":         days,
	})
}
