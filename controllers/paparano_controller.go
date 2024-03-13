package controllers

import (
	"net/http"
	"short-link/auth"
	"short-link/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func DeletePaparaNo(c *gin.Context) {
	paparano := models.PaparaNo{}
	id := c.Param("id")

	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)

	iid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	paparanoInfo, err := paparano.FindPaparaNoById(iid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	if tokenUser.ID != paparanoInfo.UserId {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Yetkiniz yok!"})
		return
	}
	err = paparano.DeletePaparaNoById(paparanoInfo.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, id)
}

func UpdatePaparaNo(c *gin.Context) {
	paparano := models.PaparaNo{}
	c.BindJSON(&paparano)
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)

	balancInfofromDB, err := paparano.FindPaparaNoById(paparano.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	if balancInfofromDB.UserId != tokenUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Yetkiniz yok!"})
		return
	}
	paparanoInfoUpdated, err := paparano.UpdatePaparaNo(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paparanoInfoUpdated)
}
