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
func CreateUser(c *gin.Context) {
	user := models.User{}
	c.BindJSON(&user)

	userfromdb, _ := user.FindUserByUserName(user.UserName)
	if userfromdb.UserName != "" {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Kullanıcı adı zaten kullanılıyor"})
		return
	}
	user.Role = "user"
	user.Admin = false
	//Veri tabanına kaydetme
	userll, err := user.CreateUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	token, _ := auth.GenerateTokenForUser(userll)
	responseUser := models.ResponseUser{
		ID:       userll.ID,
		UserName: userll.UserName,
		Admin:    userll.Admin,
		Role:     userll.Role,
		Balance:  userll.Balance,
		UrlCount: userll.UrlCount,
	}
	c.JSON(http.StatusOK, gin.H{"token": token,
		"user": responseUser,
	})
}

// Bütün kulanıcıları çekmemizi sağlayan fonksiyona http üzerinden erişmeyi sağlayan fonksiyon
func GetAllUsers(c *gin.Context) {
	user := models.User{}
	if !auth.CheckIsAdmin(c) {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}
	users, err := user.FindAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Kullanıcı adına göre kullanıcıyı çekmemizi sağlayan fonksiyona http üzerinden erişmeyi sağlayan fonksiyon
func GetUserByID(c *gin.Context) {
	user := models.User{}
	balanceInfo := models.BalanceInfo{}
	pprano := models.PaparaNo{}
	id := c.Param("id")
	idd, _ := primitive.ObjectIDFromHex(id)
	claims, _ := auth.ValidateUseToken(c)
	//tokenUser := auth.ClaimsToUser(claims)

	result, err := user.FindUserByID(idd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Kullanıcı bulunamadı"})
		return
	}
	userinfo, _ := balanceInfo.FindBalanceInfoByUserId(result.ID)
	paparano, _ := pprano.FindPaparaNoByUserId(result.ID)

	result.Password = ""
	if auth.CheckIsAdmin(c) || claims["user_id"] == id {
		fmt.Println("userid:", claims["user_id"])
		result.PaparaNo = paparano
		result.BalanceInfo = userinfo
	}
	c.JSON(http.StatusOK, result)
}

// Giriş yapmamızı sağlayan fonkiyonu http üzerinden çağıran fonksiyon
func Login(c *gin.Context) {
	user := models.User{}
	c.BindJSON(&user)

	//Kullanıcı adı ve şifre kontrol edilir
	result, err := user.FindUserByUserName(user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Kullanıcı bulunamadı"})
		return
	}
	err = models.ComparePasswords(result.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Şifre Hatalı"})
		return
	}
	token, _ := auth.GenerateTokenForUser(result)
	responseUser := models.ResponseUser{
		ID:       result.ID,
		UserName: result.UserName,
		Admin:    result.Admin,
		Role:     result.Role,
		Balance:  result.Balance,
		UrlCount: result.UrlCount,
	}
	c.JSON(http.StatusOK, gin.H{"token": token,
		"user": responseUser,
	})
}

func UpdatePassword(c *gin.Context) {
	user := models.User{}
	updateUser := models.UpdatePasswordUser{}
	c.BindJSON(&updateUser)
	if updateUser.NewPassword == "" && updateUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Lütfen şifre giriniz"})
		return
	}
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	result, err := user.FindUserByID(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Kullanıcı bulunamadı"})
		return
	}
	err = models.ComparePasswords(result.Password, updateUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Şifre Hatalı"})
		return
	}
	err = user.UpdatePassword(result.ID, updateUser.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"SUCCESS": "Şifre Güncellendi"})
}

// update user
func UpdateUser(c *gin.Context) {
	user := models.User{}
	c.BindJSON(&user)
	claims, err := auth.ValidateUseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)

	result, err := user.FindUserByID(tokenUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Kullanıcı bulunamadı"})
		return
	}
	users, err := user.FindAllUsers()
	for err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	if user.UserName != result.UserName {
		for i, _ := range users {
			if users[i].UserName == user.UserName {
				c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Bu kullancı adı kullanılıyor."})
				return
			}
		}
	}
	err = user.UpdateUser(result.ID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"SUCCESS": "Kullanıcı Güncellendi"})
}

func ForgotPassword(c *gin.Context) {
	user := models.User{}
	data := models.ForgotPassword{}
	c.BindJSON(&data)
	userFromDb, err := user.FindUserByUserMail(data.Mail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Böyle bir kullanıcı yok!"})
		return
	}

	token, err := auth.GenerateTokenForForgotPassword(userFromDb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	url := data.Domain + token
	err = auth.SendForgotPasswordEmail(userFromDb.Mail, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"SUCCESS": "Mail gönderildi"})
}

func NewPassword(c *gin.Context) {
	user := models.User{}
	updatepass := models.UpdatePasswordUser{}
	c.BindJSON(&updatepass)
	token := c.Param("token")
	claims, err := auth.ValidateForgotPasswordToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}
	tokenUser := auth.ClaimsToUser(claims)
	if tokenUser.ExpDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Token süresi doldu"})
		return
	}

	err = user.UpdatePassword(tokenUser.ID, updatepass.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": tokenUser.ID})
		return
	}
	c.JSON(http.StatusOK, gin.H{"SUCCESS": "Şifre Güncellendi",
		"req": tokenUser})
}

// update blocked by Admin
func UpdateBlocked(c *gin.Context) {
	user := models.User{}
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Geçersiz ID"})
		return
	}
	if auth.CheckIsAdmin(c) {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}
	userfromDB, err := user.FindUserByID(objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	if userfromDB.Admin {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Admin kullanıcılar engellenemez"})
		return
	}
	err = user.UpdateBlocked(objectID, !userfromDB.Blocked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"SUCCESS": "Kullanıcı Güncellendi"})
}

func DeleteUserByAdmin(c *gin.Context) {
	user := models.User{}
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Geçersiz ID"})
		return
	}
	if auth.CheckIsAdmin(c) {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Yetkiniz yok"})
		return
	}
	userfromDB, err := user.FindUserByID(objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	if userfromDB.Admin {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Admin kullanıcılar silinemez"})
		return
	}
	err = user.DeleteUser(objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"SUCCESS": "Kullanıcı Silindi"})
}
