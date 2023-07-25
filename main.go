package main

import (
	"log"
	"os"
	"short-link/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//Url için api linkleri

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	//config.AllowOrigins = []string{"http://localhost:3001"}
	router.Use(cors.New(config))

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
