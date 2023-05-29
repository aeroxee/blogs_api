package handlers

import (
	"net/http"

	"github.com/aZ4ziL/blogs_api/models"
	"github.com/gin-gonic/gin"
)

// TagHandlerGET is function to handle request tag, to get data from tag.
func TagHandlerGET(ctx *gin.Context) {
	id := getQueryInt(ctx.Request, "id", 0)
	if id != 0 {
		tag, err := models.GetTagByID(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Tag yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		ctx.JSON(http.StatusOK, tag)
		return
	}

	tags := models.GetAllTags()
	ctx.JSON(http.StatusOK, tags)
}

// TagHandlerPOST is function to get handler request to create new tag.
func TagHandlerPOST(ctx *gin.Context) {
	payloads := struct {
		Title string `json:"title" validate:"required"`
	}{}
	if err := ctx.ShouldBindJSON(&payloads); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload yang anda gunakan tidak diijinkan.",
		})
		return
	}

	if payloads.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Pastikan tag memiliki judul.",
		})
		return
	}

	tag := models.Tag{
		Title: payloads.Title,
	}
	if err := models.CreateNewTags(&tag); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, tag)
}
