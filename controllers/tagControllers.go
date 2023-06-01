package controllers

import (
	"github.com/aZ4ziL/blogs_api/handlers"
	"github.com/gin-gonic/gin"
)

func TagControllerNoAuth(group *gin.RouterGroup) {
	group.GET("", handlers.TagHandlerGET)
}

func TagControllerWithAuth(group *gin.RouterGroup) {
	group.POST("", handlers.TagHandlerPOST)
}
