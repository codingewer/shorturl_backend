package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"short-link/models"

	"github.com/gin-gonic/gin"
)

// Veritabanına kayeden fonsiyonu çağırıp http ile bağlanmamızı sağlayan fonksiyon
func CreateUser(c *gin.Context) {
	user := models.User{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": http.StatusUnprocessableEntity})
		return
	}

	//Gelen veriyi User nsnesine dönüştürdük
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": http.StatusUnprocessableEntity})
		return
	}
	userr := models.User{}
	//Kullanıcı adı benzersiz mi diye kontrol edilir
	users, err := userr.FindAllUsers()
	for i, _ := range users {
		if users[i].UserName == user.UserName {
			c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Bu kullancı adı kullanılıyor."})
			return
		}
	}
	//Veri tabanına kaydetme
	userll, err := user.CreateUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userll)
}

// Bütün kulanıcıları çekmemizi sağlayan fonksiyona http üzerinden erişmeyi sağlayan fonksiyon
func GetAllUsers(c *gin.Context) {
	user := models.User{}
	users, err := user.FindAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		fmt.Println("3")
		return
	}
	c.JSON(http.StatusOK, users)
}

// Kullanıcı adına göre kullanıcıyı çekmemizi sağlayan fonksiyona http üzerinden erişmeyi sağlayan fonksiyon
func GetByUserName(c *gin.Context) {
	user := models.User{}
	username := c.Param("username")
	result, err := user.FindByUserName(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Giriş yapmamızı sağlayan fonkiyonu http üzerinden çağıran fonksiyon
func Login(c *gin.Context) {
	user := models.User{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": http.StatusUnprocessableEntity})
		return
	}

	err = json.Unmarshal(body, &user)

	//Kullanıcı adı ve şifre kontrol edilir
	result, err := user.FindByUserName(user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Kullanıcı bulunamadı"})
		return
	}

	if user.Password != result.Password {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Şifre Hatalı"})
		return
	}
	c.JSON(http.StatusOK, result)
}
