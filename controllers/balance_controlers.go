package controllers

import (
	"fmt"
	"net/http"
	"short-link/auth"
	"short-link/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewBalanceRequests(c *gin.Context) {
	balance := models.BalanceRequest{}
	c.BindJSON(&balance)
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	user := models.User{}
	userFromDB, err := user.FindUserByID(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userFromDB.Balance < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bakiye 10 TL'den az!"})
		return
	}
	if userFromDB.Balance < balance.Amount {
		fmt.Println(tokenUser)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bakiye yetersiz"})
		return
	}
	balance.UserId = tokenUser.ID
	balance.User = userFromDB
	balanceSaved, err := balance.CreateNewRequest()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balanceSaved)
}

func GetBalanceRequests(c *gin.Context) {
	balance := models.BalanceRequest{}
	// get status param
	status := c.Param("status")
	stats, err := strconv.ParseBool(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	balanceRequests, err := balance.FindRequestsByStatus(stats)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error12": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balanceRequests)
}
