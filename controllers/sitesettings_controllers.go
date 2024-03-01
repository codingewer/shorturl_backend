package controllers

import (
	"net/http"
	"short-link/auth"
	"short-link/models"

	"github.com/gin-gonic/gin"
)

func GetBySiteName(c *gin.Context) {
	setting := models.Settings{}
	siteName := c.Param("siteName")
	stgns, err := setting.FindBySiteName(siteName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stgns)
}

// update
func UpdateSiteSettings(c *gin.Context) {
	setting := models.Settings{}
	err := c.BindJSON(&setting)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	siteName := c.Param("siteName")
	updated, err := setting.UpdateSettings(siteName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update Success",
		"data": updated})
}
