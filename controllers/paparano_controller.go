package controllers

import (
	"short-link/auth"
	"short-link/models"

	"github.com/gin-gonic/gin"
)

func NewPaparaNo(c *gin.Context) {
	paparaNo := models.PaparaNo{}
	c.BindJSON(&paparaNo)
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(401, gin.H{"ERROR": "Unauthorized"})
		return
	}
	tokeUser := auth.ClaimsToUser(claims)
	noFroDb, _ := paparaNo.FindPaparaNoByUserId(tokeUser.ID)
	if noFroDb.UserId == tokeUser.ID {
		c.JSON(400, gin.H{"ERROR": "Bu zaten var"})
		return
	}
	noSaved, err := paparaNo.CreatePaparaNo(tokeUser.ID)
	if err != nil {
		c.JSON(400, gin.H{"ERROR": "Bad Request"})
		return
	}
	c.JSON(201, noSaved)
}
