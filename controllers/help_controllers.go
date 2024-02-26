package controllers

import (
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
	param := c.Param("status")
	status, err := strconv.ParseBool(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	err = helpReq.ChangeStatus(status)
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
	param := c.Param("status")
	status, err := strconv.ParseBool(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	requests, err := helpReq.FindByStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}
