package controllers

import (
	"net/http"
	"short-link/auth"
	"short-link/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUrlFaq(c *gin.Context) {
	UrlFaq := models.UrlFaq{}
	c.BindJSON(&UrlFaq)
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"ERROR": "Unauthorized"})
		return
	}
	UrlFaqSaved, err := UrlFaq.NewUrlFaq()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, UrlFaqSaved)
}

func GetUrlFaqs(c *gin.Context) {
	UrlFaq := models.UrlFaq{}
	UrlFaqs, err := UrlFaq.FindAllUrlFaqs()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, UrlFaqs)
}

func GetUrlFaq(c *gin.Context) {
	UrlFaq := models.UrlFaq{}
	id := c.Param("id")
	UrlFaq.ID, _ = primitive.ObjectIDFromHex(id)
	UrlFaq, err := UrlFaq.FindUrlFaqByID()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, UrlFaq)
}

// update UrlFaq
func UpdateUrlFaq(c *gin.Context) {
	UrlFaq := models.UrlFaq{}
	c.BindJSON(&UrlFaq)
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}
	err := UrlFaq.UpdateUrlFaqByID()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, UrlFaq)
}

// delete by id
func DeleteUrlFaq(c *gin.Context) {
	UrlFaq := models.UrlFaq{}
	id := c.Param("id")
	UrlFaq.ID, _ = primitive.ObjectIDFromHex(id)
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}
	err := UrlFaq.DeleteUrlFaqByID()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": UrlFaq.ID})
}
