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
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	router.Use(cors.New(config))

	url := router.Group("url")
	url.POST("/add", controllers.ShortLink)
	url.PUT("/update", controllers.UpdateUrl)
	url.GET("/getall", controllers.GetAll)
	url.GET("/getbyid/:id", controllers.GetByID)
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
	user.POST("/forgotpassword", controllers.ForgotPassword)
	user.POST("/resetpassword/:token", controllers.NewPassword)

	balance := router.Group("balance")
	balance.POST("/add", controllers.NewBalanceRequests)
	balance.GET("/getbystatus/:status", controllers.GetBalanceRequests)
	balance.PUT("/updatestatus/:status", controllers.UpdateBalanceRequest)
	balance.GET("/getbyuser", controllers.GetBalanceRequestsById)
	balance.PUT("/info/updateinfo", controllers.UpdateBalanceInfo)
	balance.GET("/info/getbyuserId", controllers.GetBalanceInfo)
	balance.GET("/info/getall", controllers.FindAllBalanceInfo)
	balance.POST("/info/new", controllers.NewBalanceInfo)
	balance.DELETE("/info/delete/:id", controllers.DeleteBalanceInfo)
	balance.DELETE("/papara/delete/:id", controllers.DeletePaparaNo)
	balance.PUT("/papara/updateinfo", controllers.UpdatePaparaNo)
	balance.POST("/papara/new", controllers.NewPaparaNo)

	help := router.Group("help")
	help.POST("/new", controllers.NewHelpRequest)
	help.GET("/getbystatus/:stats", controllers.GetHelpRequestsByStatus)
	help.GET("/getbyuser", controllers.GetHelpRequestsByUser)
	help.PUT("/updatestatus/", controllers.ChangeHelpRequestStatus)

	faq := router.Group("faq")
	faq.POST("/new", controllers.CreateFaq)
	faq.GET("/getall", controllers.GetFaqs)
	faq.GET("/getbyid/:id", controllers.GetFaq)
	faq.PUT("/update", controllers.UpdateFaq)
	faq.DELETE("/delete/:id", controllers.DeleteFaq)

	seen := router.Group("seen")
	seen.GET("/userseen/:days", controllers.GetUserSeenData)

	sitesett := router.Group("sitesettings")
	sitesett.GET("/getbysite/:siteName", controllers.GetBySiteName)
	sitesett.PUT("/update/:siteName", controllers.UpdateSiteSettings)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8180"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
