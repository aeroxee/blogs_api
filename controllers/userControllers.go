package controllers

import (
	"github.com/aZ4ziL/blogs_api/handlers"
	"github.com/gin-gonic/gin"
)

func UserControllerNoAuth(group *gin.RouterGroup) {
	group.POST("/get-token", handlers.UserHandlerGetToken)

	group.POST("/register", handlers.UserHandlerRegister)
}

func UserControllerWithAuth(group *gin.RouterGroup) {
	group.GET("/auth", handlers.UserHandlerAuth)
}
