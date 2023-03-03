package main

import (

	"github.com/gin-gonic/gin"
	"github.com/mhmdKhasawneh/musicrecommendationapp/controllers"
	"github.com/mhmdKhasawneh/musicrecommendationapp/initializers"
	"github.com/mhmdKhasawneh/musicrecommendationapp/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDb()
	initializers.MigrateDatabase()
}

func main() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.POST("/api/user/signup", controllers.Signup)
	r.POST("/api/user/login", controllers.Login)
	r.GET("/api/user/tokenlogin", controllers.TokenLogin)

	r.POST("/api/recommendation/recommend", middleware.ExtractUserFromLocalStorage, controllers.Recommend)
	r.GET("/api/recommendation/get", middleware.ExtractUserFromLocalStorage, controllers.GetRecommendations)

	r.Run()
}
