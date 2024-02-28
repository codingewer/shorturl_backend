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
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	//config.AllowOrigins = []string{"http://localhost:3001"}
	router.Use(cors.New(config))

	url := router.Group("url")
	url.POST("/add", controllers.ShortLink)
	url.GET("/getall", controllers.GetAll)
	url.DELETE("/delete/:id", controllers.DeleteByID)
	url.GET("/get/:shortenedurl", controllers.GetByUrl)
	url.GET("/getbycreatedby", controllers.GetByCreatedBy)

	//Kullanıcılar için api linkleri
	user := router.Group("user")
	user.GET("/getbyId", controllers.GetUserByID)
	user.GET("/getall", controllers.GetAllUsers)
	user.POST("/new", controllers.CreateUser)
	user.POST("/login", controllers.Login)
	user.PUT("/update", controllers.UpdateUser)
	user.PUT("/updatepassword", controllers.UpdatePassword)

	balance := router.Group("balance")
	balance.POST("/add", controllers.NewBalanceRequests)
	balance.GET("/getbystatus/:status", controllers.GetBalanceRequests)
	balance.PUT("/updatestatus/:status", controllers.UpdateBalanceRequest)
	balance.GET("/getbyuser", controllers.GetBalanceRequestsById)
	balance.PUT("/info/updateinfo", controllers.UpdateBalanceInfo)
	balance.GET("/info/getbyuserId", controllers.GetBalanceInfo)

	help := router.Group("help")
	help.POST("/new", controllers.NewHelpRequest)
	help.GET("/getbystatus/:status", controllers.GetHelpRequestsByStatus)
	help.GET("/getbyuser", controllers.GetHelpRequestsByUser)
	help.PUT("/updatestatus/:status", controllers.ChangeHelpRequestStatus)

	seen := router.Group("seen")
	seen.GET("/userseen/:days", controllers.GetUserSeenData)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8180"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

}
