package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"short-link/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Veritabanına kayeden fonsiyonu çağırıp http ile bağlanmamızı sağlayan fonksiyon
func ShortLink(c *gin.Context) {
	url := models.Url{}
	c.BindJSON(&url)
	//kullanıcı isim vermezse rastgele 10 karakterlik benzersiz bir isim oluştrulur
	urll := models.Url{}
	urls, err := urll.FindAllUrl()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"ERROR": "Bu zaten var"})
		return
	}

	if url.OrginalUrl[0:8] != "https://" {
		if url.OrginalUrl[0:7] == "http://" {
			fmt.Println("değişmiyor")
		} else {
			url.OrginalUrl = "https://" + url.OrginalUrl
		}
	}
	var shortenrdUrl string
	if url.ShortenedUrl == "" {
		for i, _ := range urls {
			shortenrdUrl = models.GenerateString(10)
			if shortenrdUrl == urls[i].ShortenedUrl {
				url.ShortenedUrl = models.GenerateString(10)
				break
			} else {
				url.ShortenedUrl = shortenrdUrl
			}

		}
	} else {
		for i, _ := range urls {
			if len(url.ShortenedUrl) > 10 {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"ERROR": "10 karakterden uzun olamaz"})
				return
			}
			if urls[i].ShortenedUrl == url.ShortenedUrl {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"ERROR": "Bu zaten var"})
				return
			}
		}
	}
	//kullanıcı adı zorunlu hale getirildi
	if url.CreatedBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Kullanıcı girişi gerekli"})
		return
	}
	//kullanıcını olup olmadığı ve kullanıcı link oluştruma seviyesi artırılır
	user := models.User{}
	err = user.NewLinkCount(url.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Yeni link seviyesi eklenirken hata"})
		return
	}
	url.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	//Veri tabanına kaydetme
	urlll, err := url.ShortLink()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, urlll)
}

// Bütün linkleri önce en çok tıklanan olmak üzere çeken fonksiyonu http üzerinden bağlanmamızı sağlayan fonksiyon
func GetAll(c *gin.Context) {
	url := models.Url{}
	urls, err := url.FindAllUrl()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		fmt.Println("3")
		return
	}
	c.JSON(http.StatusOK, urls)
}

// Link adına göre veri tabanında linki çeken fonksiyonu http üzerinden bağlanmamızı sağlayan fonksiyon
func GetByUrl(c *gin.Context) {
	url := models.Url{}
	title := c.Param("shortenedurl")
	result, err := url.FindByUrl(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Link İdsine göre veri tabanından linki  silen fonksiyonu http üzerinden bağlanmamızı sağlayan fonksiyon
func DeleteByID(c *gin.Context) {
	url := models.Url{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": http.StatusUnprocessableEntity})
		return
	}
	//Gelen veriyi Link yapısına dönüştürdük
	err = json.Unmarshal(body, &url)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": err.Error()})
		return
	}
	id := c.Param("id")
	idd, _ := primitive.ObjectIDFromHex(id)

	urll := models.Url{}
	urlll, err := urll.FindByID(idd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	if urlll.CreatedBy != url.CreatedBy {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}

	err = url.DeleteByID(idd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sdsd": "sfsf"})
}

// kullanıcı adına göre veri tabanında linkleri çeken fonksiyonu http üzerinden bağlanmamızı sağlayan fonksiyon
func GetByCreatedBy(c *gin.Context) {
	url := models.Url{}
	username := c.Param("username")
	result, err := url.FindByCreatedBy(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
