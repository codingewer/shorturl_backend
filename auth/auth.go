package auth

import (
	"errors"
	"fmt"
	"os"
	"short-link/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserClaims struct {
	UserID   primitive.ObjectID `json:"user_id"`
	Role     string             `json:"role"`
	Username string             `json:"username"`
	ExpDate  int64              `json:"expDate"`
	jwt.StandardClaims
}

func GenerateTokenForUser(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = UserClaims{
		UserID:   user.ID,
		Role:     user.Role,
		Username: user.UserName,
	}
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
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

func ClaimsToUser(claims jwt.MapClaims) models.ResponseUser {
	ID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		return models.ResponseUser{}
	}
	return models.ResponseUser{
		ID:       ID,
		UserName: claims["username"].(string),
		Role:     claims["role"].(string),
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
