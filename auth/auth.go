package auth

import (
	"errors"
	"fmt"
	"os"
	"short-link/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserClaims struct {
	ID       primitive.ObjectID `json:"user_id"`
	Role     string             `json:"role"`
	Mail     string             `json:"mail"`
	UserName string             `json:"username"`
	ExpDate  time.Time          `json:"expDate"`
	jwt.StandardClaims
}

type UserClaimsResponse struct {
	ID       primitive.ObjectID `json:"user_id"`
	Role     string             `json:"role"`
	Mail     string             `json:"mail"`
	UserName string             `json:"username"`
	ExpDate  time.Time          `json:"expDate"`
}

func GenerateTokenForUser(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = UserClaims{
		ID:       user.ID,
		Role:     user.Role,
		Mail:     user.Mail,
		UserName: user.UserName,
	}
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GENERATE FOR forgot paswword
func GenerateTokenForForgotPassword(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = UserClaims{
		ID:       user.ID,
		Role:     user.Role,
		ExpDate:  time.Now().Add(30 * time.Minute),
		UserName: user.UserName,
	}
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// valitate for forgot password
func ValidateForgotPasswordToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil

	} else {
		return nil, errors.New("Invalid token claims")
	}
}

func ValidateUseToken(c *gin.Context) (jwt.MapClaims, error) {
	bearerToken := c.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		bearerToken = strings.Split(bearerToken, " ")[1]
	}
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil

	} else {
		return nil, errors.New("Invalid token claims")
	}
}

func ClaimsToUser(claims jwt.MapClaims) UserClaimsResponse {
	ID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		return UserClaimsResponse{}
	}
	layout := time.RFC3339Nano
	expDateStr := claims["expDate"].(string)
	t, err := time.Parse(layout, expDateStr)
	if err != nil {
		fmt.Println("Tarih ayrıştırılamadı:", err)
		return UserClaimsResponse{} // Tarih ayrıştırılamazsa hata durumu
	}
	fmt.Println(claims["expDate"].(string))
	fmt.Println(t)

	return UserClaimsResponse{
		ID:       ID,
		UserName: claims["username"].(string),
		Role:     claims["role"].(string),
		Mail:     claims["mail"].(string),
		ExpDate:  t,
	}
}

func CheckIsAdmin(c *gin.Context) bool {
	user := models.User{}
	claims, err := ValidateUseToken(c)
	if err != nil {
		return false
	}
	tokenUser, err := user.FindUserByUserName(claims["username"].(string))
	if err != nil {
		return false
	}
	if !tokenUser.Admin {
		return false
	}
	return true
}
