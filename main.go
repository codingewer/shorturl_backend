package main

import (
	"log"
	"os"
	"short-link/controllers"
	"short-link/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	go GetByTime()
	router := gin.Default()
	//Url için api linkleri

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	//config.AllowOrigins = []string{"http://localhost:3001"}
	router.Use(cors.Default())

	url := router.Group("url")
	url.POST("/add", controllers.ShortLink)
	url.GET("/getall", controllers.GetAll)
	url.DELETE("/delete/:id", controllers.DeleteByID)
	url.GET("/get/:shortenedurl", controllers.GetByUrl)
	url.GET("/getbycreatedby/:username", controllers.GetByCreatedBy)

	//Kullanıcılar için api linkleri
	user := router.Group("user")
	user.GET("/get/:username", controllers.GetByUserName)
	user.GET("/getall", controllers.GetAllUsers)
	user.POST("/new", controllers.CreateUser)
	user.POST("/login", controllers.Login)

	balance := router.Group("balance")
	balance.POST("/add", controllers.NewBalanceRequests)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

}

func GetByTime() {
	ticker := time.NewTicker(30 * time.Minute)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			url := models.Url{}
			_, _ = url.FindAllUrl()

		}
	}
}
