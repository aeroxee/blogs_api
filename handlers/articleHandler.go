package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aZ4ziL/blogs_api/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
)

// ArticleHandlerGET is function to handling request GET for article.
func ArticleHandlerGET(ctx *gin.Context) {
	offset := getQueryInt(ctx.Request, "offset", 0)
	limit := getQueryInt(ctx.Request, "limit", 10)
	sorted := getQueryString(ctx.Request, "sorted", "desc")
	slug := getQueryString(ctx.Request, "slug", "")

	if slug != "" {
		article, err := models.GetArticleBySlug(slug)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Halaman yang anda tuju tidak dapat ditemukan.",
			})
			return
		}
		ctx.JSON(http.StatusOK, article)
		return
	}

	articles := models.GetAllArticles(offset, limit, sorted)
	ctx.JSON(http.StatusOK, articles)
}

// ArticleHandlerPOST is function handler to handling request to create new user.
func ArticleHandlerPOST(ctx *gin.Context) {
	payloads := struct {
		Tags        string                `form:"tags"`
		UserID      int                   `form:"user_id" validate:"required"`
		Title       string                `form:"title" validate:"required,max=50"`
		Logo        *multipart.FileHeader `form:"logo" validate:"required"`
		Description string                `form:"description" validate:"required,max=255"`
		Content     string                `form:"content" validate:"required"`
		Status      string                `form:"status" validate:"required,max=9"`
	}{}
	if err := ctx.ShouldBindWith(&payloads, binding.FormMultipart); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload yang anda gunakan tidak diijinkan.",
		})
		return
	}

	validate = validator.New()
	if err := validate.Struct(&payloads); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Error pada field `%s` dengan kode error `%s`", err.Field(), err.ActualTag()))
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kesalahan pada validasi field.",
			"errors":  errorMessages,
		})
		return
	}

	slug := slug.MakeLang(payloads.Title, "id")

	var tags []*models.Tag
	tagSplit := strings.Split(payloads.Tags, ",")
	for _, tag := range tagSplit {
		t, err := models.GetTagByTitle(tag)
		if err != nil {
			newTag := models.Tag{
				Title: tag,
			}
			err := models.CreateNewTags(&newTag)
			if err != nil {
				continue
			}
			tags = append(tags, &newTag)
			continue
		}
		tags = append(tags, &t)
	}

	// try to save
	article := models.Article{
		UserID:      payloads.UserID,
		Tags:        tags,
		Title:       payloads.Title,
		Slug:        slug,
		Description: payloads.Description,
		Content:     payloads.Description,
		Status:      payloads.Status,
	}

	// save
	if err := models.CreateNewArticle(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("/media/articles/%s/%s", slug, payloads.Logo.Filename)

	if err := ctx.SaveUploadedFile(payloads.Logo, "."+filename); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	article.Logo = filename

	models.GetDB().Save(&article)

	ctx.JSON(http.StatusCreated, article)
}
