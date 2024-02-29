package controllers

import (
	"fmt"
	"net/http"
	"short-link/auth"
	"short-link/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewHelpRequest(c *gin.Context) {
	helpReq := models.HelpRequest{}
	c.BindJSON(&helpReq)
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	if helpReq.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Mesaj alanı boş olamaz!"})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	requestSaved, err := helpReq.NewHelpRequest(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requestSaved)
}

func ChangeHelpRequestStatus(c *gin.Context) {
	helpReq := models.HelpRequest{}
	c.BindJSON(&helpReq)
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Bu işlemi yapmak için yetkiniz yok!"})
		return
	}
	err := helpReq.ChangeStatus()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetHelpRequestsByStatus(c *gin.Context) {
	helpReq := models.HelpRequest{}
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Bu işlemi yapmak için yetkiniz yok!"})
		return
	}
	param := c.Param("stats")
	status, _ := strconv.ParseBool(param)
	fmt.Println(status)
	requests, err := helpReq.FindByStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	for i, _ := range requests {
		usr := models.User{}
		user, err := usr.FindResposeUserByID(requests[i].UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
			return
		}
		requests[i].User = user
	}

	for i, _ := range requests {
		balanceİnfo := models.BalanceInfo{}
		balance, _ := balanceİnfo.FindBalanceInfoByUserId(requests[i].UserID)
		requests[i].User.BalanceInfo = balance
	}
	c.JSON(http.StatusOK, requests)
}

func GetHelpRequestsByUser(c *gin.Context) {
	helpReq := models.HelpRequest{}
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	requests, err := helpReq.FindByUserId(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}
