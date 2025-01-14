package main

import (
	"github.com/ArdhanaGusti/Golang_api/config"
	"github.com/ArdhanaGusti/Golang_api/middleware"
	"github.com/ArdhanaGusti/Golang_api/routes"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/auth/check", middleware.IsAuth(), routes.CheckToken)
		v1.GET("/auth/:provider", routes.RedirectHandler)
		v1.GET("/auth/:provider/callback", routes.CallbackHandler)

		v1.POST("/auth/register", routes.RegisterUser)
		v1.POST("/auth/login", routes.LoginUser)
		v1.PATCH("/auth/change-role", middleware.IsAuth(), routes.ChangeRole)

		v1.GET("/auth/profile", middleware.IsAuth(), routes.GetProfile)

		v1.GET("/article", middleware.IsAuth(), routes.Home)
		v1.GET("/article/:slug", routes.GetArticle)
		v1.POST("/article", middleware.IsAuth(), routes.PostArticle)
		v1.PUT("/article/:slug", middleware.IsAuth(), routes.UpdateArticle)
		v1.DELETE("/article/:slug", middleware.IsAdmin(), routes.DeleteArticle)
	}

	return r
}

func main() {
	gotenv.Load()
	config.InitDB()
	config.InitRedis()

	r := setupRouter()
	r.Run(":8080")
}
