package controllers

import (
	"net/http"
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
