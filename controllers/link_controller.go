package controllers

import (
	"fmt"
	"net/http"
	"short-link/auth"
	"short-link/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Veritabanına kayeden fonsiyonu çağırıp http ile bağlanmamızı sağlayan fonksiyon
func ShortLink(c *gin.Context) {
	url := models.Url{}
	user := models.User{}
	c.BindJSON(&url)
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	//kullanıcı isim vermezse rastgele 10 karakterlik benzersiz bir isim oluştrulur
	urls, err := url.FindAllUrl()
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
	//kullanıcını olup olmadığı ve kullanıcı link oluştruma seviyesi artırılır
	err = user.NewLinkCount(tokenUser.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Yeni link seviyesi eklenirken hata"})
		return
	}
	url.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	url.CreatedBy = tokenUser.UserName
	url.UserID = tokenUser.ID
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
	seen := models.Seen{}
	err = seen.NewSeen(result.UserID, result.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Link İdsine göre veri tabanından linki  silen fonksiyonu http üzerinden bağlanmamızı sağlayan fonksiyon
func DeleteByID(c *gin.Context) {
	url := models.Url{}
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	c.BindJSON(url)
	id := c.Param("id")
	idd, _ := primitive.ObjectIDFromHex(id)

	urlll, err := url.FindByID(idd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	if urlll.CreatedBy != tokenUser.UserName && urlll.UserID != tokenUser.ID {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}

	err = url.DeleteByID(idd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

// kullanıcı adına göre veri tabanında linkleri çeken fonksiyonu http üzerinden bağlanmamızı sağlayan fonksiyon
func GetByCreatedBy(c *gin.Context) {
	url := models.Url{}
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	result, err := url.FindByCreatedBy(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
