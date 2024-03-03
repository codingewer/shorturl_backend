package controllers

import (
	"net/http"
	"short-link/auth"
	"short-link/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateFaq(c *gin.Context) {
	faq := models.Faq{}
	c.BindJSON(&faq)
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"ERROR": "Unauthorized"})
		return
	}
	faqSaved, err := faq.NewFaq()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, faqSaved)
}

func GetFaqs(c *gin.Context) {
	faq := models.Faq{}
	faqs, err := faq.FindAllFaqs()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, faqs)
}

func GetFaq(c *gin.Context) {
	faq := models.Faq{}
	id := c.Param("id")
	faq.ID, _ = primitive.ObjectIDFromHex(id)
	faq, err := faq.FindFaqByID()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, faq)
}

// update faq
func UpdateFaq(c *gin.Context) {
	faq := models.Faq{}
	c.BindJSON(&faq)
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}
	err := faq.UpdateFaqByID()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "güncellendi"})
}

// delete by id
func DeleteFaq(c *gin.Context) {
	faq := models.Faq{}
	id := c.Param("id")
	faq.ID, _ = primitive.ObjectIDFromHex(id)
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}
	err := faq.DeleteFaqByID()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": faq.ID})
}
