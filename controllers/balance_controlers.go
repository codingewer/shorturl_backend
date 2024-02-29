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
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	user := models.User{}
	userFromDB, err := user.FindResposeUserByID(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	if userFromDB.Balance < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Bakiye 10 TL'den az!"})
		return
	}
	if userFromDB.Balance < balance.Amount {
		fmt.Println(tokenUser)
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Bakiye yetersiz"})
		return
	}
	balance.UserId = tokenUser.ID
	balance.User = userFromDB
	balanceSaved, err := balance.CreateNewRequest()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balanceSaved)
}

func GetBalanceRequests(c *gin.Context) {
	balance := models.BalanceRequest{}
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Yetkiniz yok!"})
		return
	}
	// get status param
	status := c.Param("status")
	stats, err := strconv.ParseBool(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	balanceRequests, err := balance.FindRequestsByStatus(stats)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	for i, _ := range balanceRequests {
		usr := models.User{}
		user, _ := usr.FindResposeUserByID(balanceRequests[i].UserId)
		balanceRequests[i].User = user
		balanceİnfo := models.BalanceInfo{}
		balance, _ := balanceİnfo.FindBalanceInfoByUserId(balanceRequests[i].UserId)
		balanceRequests[i].User.BalanceInfo = balance
	}
	c.JSON(http.StatusOK, balanceRequests)
}

func GetBalanceRequestsById(c *gin.Context) {
	balance := models.BalanceRequest{}
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Yetkiniz yok!"})
		return
	}
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	user := models.User{}
	userFromDB, err := user.FindResposeUserByID(tokenUser.ID)
	balanceRequests, _ := balance.FindRequestsByUserID(userFromDB.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balanceRequests)
}

func UpdateBalanceRequest(c *gin.Context) {
	balance := models.BalanceRequest{}
	c.BindJSON(&balance)
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Yetkiniz yok!"})
		return
	}
	status := c.Param("status")
	stats, err := strconv.ParseBool(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = balance.UpdateRequestStatus(stats)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balance)
}

// get info by user id
func GetBalanceInfo(c *gin.Context) {
	balance := models.BalanceInfo{}
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	fmt.Println(tokenUser.ID)
	balanceInfo, err := balance.FindBalanceInfoByUserId(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error(), "id": tokenUser.ID})
		return
	}
	c.JSON(http.StatusOK, balanceInfo)
}

// Update balance info by user id
func UpdateBalanceInfo(c *gin.Context) {
	balance := models.BalanceInfo{}
	c.BindJSON(&balance)
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)

	balancInfofromDB, err := balance.FindBalanceInfoById(balance.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	if balancInfofromDB.UserId != tokenUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Yetkiniz yok!"})
		return
	}
	balanceInfoUpdated, err := balance.UpdateBalanceInfo(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balanceInfoUpdated)
}

func FindAllBalanceInfo(c *gin.Context) {
	balance := models.BalanceInfo{}
	balanceInfo, err := balance.FindAllBalanceInfo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balanceInfo)
}
