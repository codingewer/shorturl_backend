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

	if url.OrginalUrl[0:8] != "https://" {
		if url.OrginalUrl[0:7] == "http://" {
			fmt.Println("değişmiyor")
		} else {
			url.OrginalUrl = "https://" + url.OrginalUrl
		}
	}
	if url.ShortenedUrl == "" {
		url.ShortenedUrl = models.GenerateString(10)
	}

	urlsFromDB, _ := url.FindByCreatedBy(url.UserID)
	for i, _ := range urlsFromDB {
		if urlsFromDB[i].ShortenedUrl == url.ShortenedUrl {
			c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Bu link zaten mevcut"})
			return
		}
	}
	urls, err := url.FindByCreatedBy(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Kullanıcı bulunamadı"})
		return
	}
	for i, _ := range urls {
		if urls[i].OrginalUrl == url.OrginalUrl {
			c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Bu link zaten mevcut"})
			return
		}
	}
	//kullanıcını olup olmadığı ve kullanıcı link oluştruma seviyesi artırılır
	err = user.NewLinkCount(tokenUser.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Yeni link seviyesi eklenirken hata"})
		return
	}
	userFromDB, err := user.FindUserByID(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Kullanıcı bulunamadı"})
		return
	}
	if userFromDB.Blocked {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Kullanıcı engellenmiş bunu gerçekleştiremez"})
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
	user := models.User{}
	siteSettings := models.Settings{}
	setdata, err := siteSettings.FindBySiteName("short-url")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	visitor := models.GetVisitor(c.ClientIP())
	// Eğer 10 dakikadan fazla zaman geçtiyse sayaç ve zamanı sıfırla.
	if time.Since(visitor.LastReset) > 30*time.Minute {
		visitor.Requests = 0
		visitor.LastReset = time.Now()
	}
	// İstek sayısını arttır ve kontrol et.
	visitor.Requests++
	if visitor.Requests > 10 {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"ERROR": "10 dakika içinde çok fazla istek gönderdiniz. Lütfen daha sonra tekrar deneyiniz.",
		})
		return
	}

	title := c.Param("shortenedurl")
	owner := c.Param("createdby")
	result, err := url.FindByUrl(owner, title)
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
	err = user.AddBalance(result.CreatedBy, setdata.RevenuePerClick)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
	}
	c.JSON(http.StatusOK, result)
}

// Link İdsine göre veri tabanından linki  silen fonksiyonu http üzerinden bağlanmamızı sağlayan fonksiyon
func DeleteUrlByID(c *gin.Context) {
	url := models.Url{}
	user := models.User{}
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	id := c.Param("id")
	idd, _ := primitive.ObjectIDFromHex(id)

	urlll, err := url.FindByID(idd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	err = user.DownLinkCount(urlll.CreatedBy)
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
	c.JSON(http.StatusOK, id)
}

// kullanıcı adına göre veri tabanında linkleri çeken fonksiyonu http üzerinden bağlanmamızı sağlayan fonksiyon
func GetByCreatedBy(c *gin.Context) {
	url := models.Url{}
	id := c.Param("id")
	claims, _ := auth.ValidateUseToken(c)
	idd, _ := primitive.ObjectIDFromHex(id)
	result, err := url.FindByCreatedBy(idd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	if claims["user_id"] != id {
		for i, _ := range result {
			result[i].OrginalUrl = ""
		}
	}
	c.JSON(http.StatusOK, result)
}

func UpdateUrl(c *gin.Context) {
	url := models.Url{}
	c.BindJSON(&url)
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	urll, err := url.FindByID(url.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	if urll.UserID != tokenUser.ID {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}
	err = url.Update(url.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "link başarıyla güncellendi"})
}

// get by id
func GetByID(c *gin.Context) {
	url := models.Url{}
	id := c.Param("id")
	idd, _ := primitive.ObjectIDFromHex(id)
	result, err := url.FindByID(idd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// delete url by admin
func DeleteUrlByAdmin(c *gin.Context) {
	url := models.Url{}
	user := models.User{}
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}
	id := c.Param("id")
	idd, _ := primitive.ObjectIDFromHex(id)

	urlll, err := url.FindByID(idd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	err = user.DownLinkCount(urlll.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	err = url.DeleteByID(idd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, id)
}
