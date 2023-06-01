package main

import (
	"github.com/aZ4ziL/blogs_api/controllers"
	"github.com/aZ4ziL/blogs_api/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.Static("/media", "./media")

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AddAllowMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	config.AddAllowHeaders("Content-Type", "Authorization", "Accept", "Content-Disposition")

	r.Use(cors.New(config))

	userGroupNoAuth := r.Group("/user")
	controllers.UserControllerNoAuth(userGroupNoAuth)

	userGroupWithAuth := r.Group("/user")
	userGroupWithAuth.Use(middlewares.Authentication())
	controllers.UserControllerWithAuth(userGroupWithAuth)

	// tag no auth
	tagGroupNoAuth := r.Group("/tags")
	controllers.TagControllerNoAuth(tagGroupNoAuth)

	// tag with auth
	tagGroupWithAuth := r.Group("/tags")
	tagGroupWithAuth.Use(middlewares.Authentication())
	controllers.TagControllerWithAuth(tagGroupWithAuth)

	// article no auth
	articleGroupNoAuth := r.Group("/articles")
	controllers.ArticleControllerNoAuth(articleGroupNoAuth)

	// article with auth
	articleGroupWithAuth := r.Group("/articles")
	articleGroupWithAuth.Use(middlewares.Authentication())
	controllers.ArticleControllerWithAuth(articleGroupWithAuth)

	r.Run(":8000")
}
