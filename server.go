package main

import (
	"github.com/ardhanagusti/learn-gin/config"
	"github.com/ardhanagusti/learn-gin/middleware"
	"github.com/ardhanagusti/learn-gin/routes"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	config.InitDB()
	gotenv.Load()

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/auth/check", middleware.IsAuth(), routes.CheckToken)
		v1.GET("/auth/:provider", routes.RedirectHandler)
		v1.GET("/auth/:provider/callback", routes.CallbackHandler)

		v1.GET("/auth/profile", middleware.IsAuth(), routes.GetProfile)

		v1.GET("/", middleware.IsAuth(), routes.Home)
		v1.GET("/article/:slug", routes.GetArticle)
		v1.GET("/articles/:tag", middleware.IsAuth(), routes.GetArticleTag)
		v1.POST("/", middleware.IsAuth(), routes.PostArticle)
		v1.PUT("/update/:slug", middleware.IsAuth(), routes.UpdateArticle)
		v1.DELETE("/delete/:slug", middleware.IsAdmin(), routes.DeleteArticle)
	}

	r.Run()
}
