package controllers

import (
	"net/http"
	"short-link/auth"
	"short-link/models"

	"github.com/gin-gonic/gin"
)

func NewBalanceRequests(c *gin.Context) {
	balance := models.BalanceRequest{}
	c.BindJSON(&balance)
	claims, _ := auth.ValidateUseToken(c)
	user := auth.ClaimsToUser(claims)
	c.JSON(http.StatusOK, user)
}
