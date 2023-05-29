package controllers

import (
	"github.com/aZ4ziL/blogs_api/handlers"
	"github.com/gin-gonic/gin"
)

func ArticleControllerNoAuth(group *gin.RouterGroup) {
	group.GET("", handlers.ArticleHandlerGET)
}

func ArticleControllerWithAuth(group *gin.RouterGroup) {
	group.POST("", handlers.ArticleHandlerPOST)
	group.OPTIONS("", handlers.ArticleHandlerPOST)
}
