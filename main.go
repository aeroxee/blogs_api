package main

import (
	"github.com/aZ4ziL/blogs_api/controllers"
	"github.com/aZ4ziL/blogs_api/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.Static("/media", "./media")

	userGroupNoAuth := r.Group("/user")
	controllers.UserControllerNoAuth(userGroupNoAuth)

	userGroupWithAuth := r.Group("/user")
	userGroupWithAuth.Use(middlewares.CORS())
	userGroupWithAuth.Use(middlewares.Authentication())
	controllers.UserControllerWithAuth(userGroupWithAuth)

	// tag no auth
	tagGroupNoAuth := r.Group("/tags")
	tagGroupNoAuth.Use(middlewares.CORS())
	controllers.TagControllerNoAuth(tagGroupNoAuth)

	// tag with auth
	tagGroupWithAuth := r.Group("/tags")
	tagGroupWithAuth.Use(middlewares.CORS())
	tagGroupWithAuth.Use(middlewares.Authentication())
	controllers.TagControllerWithAuth(tagGroupWithAuth)

	// article no auth
	articleGroupNoAuth := r.Group("/articles")
	articleGroupNoAuth.Use(middlewares.CORS())
	controllers.ArticleControllerNoAuth(articleGroupNoAuth)

	// article with auth
	articleGroupWithAuth := r.Group("/articles")
	articleGroupWithAuth.Use(middlewares.CORS())
	articleGroupWithAuth.Use(middlewares.Authentication())
	controllers.ArticleControllerWithAuth(articleGroupWithAuth)

	r.Run(":8000")
}
