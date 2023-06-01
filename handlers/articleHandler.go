package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/aZ4ziL/blogs_api/auth"
	"github.com/aZ4ziL/blogs_api/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ArticleHandlerGET is function to handling request GET for article.
func ArticleHandlerGET(ctx *gin.Context) {
	offset := getQueryInt(ctx.Request, "offset", 0)
	limit := getQueryInt(ctx.Request, "limit", 10)
	sorted := getQueryString(ctx.Request, "sorted", "desc")
	slug := getQueryString(ctx.Request, "slug", "")
	status := getQueryString(ctx.Request, "status", "PUBLISHED")
	q := getQueryString(ctx.Request, "q", "")
	authorId := getQueryString(ctx.Request, "authorId", "")
	token := getQueryString(ctx.Request, "token", "")

	if authorId != "" && token != "" {
		claims, err := auth.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan.",
			})
			return
		}
		user, err := models.GetUserByID(claims.Credential.UserID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan.",
			})
			return
		}

		var articlesByUser []models.Article
		models.GetDB().Model(&models.Article{}).Where("user_id = ?", user.ID).
			Limit(limit).Offset(offset).Order(fmt.Sprintf("created_at %s", sorted)).Preload("Tags").
			Find(&articlesByUser)
		ctx.JSON(http.StatusOK, articlesByUser)
		return
	}

	if q != "" {
		articles := models.GetArticleFilterByTitle(cases.Title(language.Indonesian).String(q))
		ctx.JSON(http.StatusOK, articles)
		return
	}

	if slug != "" {
		article, err := models.GetArticleBySlug(slug)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Halaman yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		article.Views = article.Views + 1

		models.GetDB().Save(&article)

		ctx.JSON(http.StatusOK, article)
		return
	}

	articles := models.GetAllArticles(offset, limit, sorted, status)
	ctx.JSON(http.StatusOK, articles)
}

// ArticleHandlerPOST is function handler to handling request to create new user.
func ArticleHandlerPOST(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(1024)
	if err != nil {
		fmt.Println(err)
	}

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
			"message": err.Error(),
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

	slugString := slug.MakeLang(payloads.Title, "id")

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
		Title:       cases.Title(language.Indonesian).String(payloads.Title),
		Slug:        slugString,
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

	filename := fmt.Sprintf("/media/articles/%s/%s", slugString, payloads.Logo.Filename)

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

// ArticleHandlerPUT is function handling request to edit article.
func ArticleHandlerPUT(ctx *gin.Context) {
	slugQuery := getQueryString(ctx.Request, "slug", "")
	if slugQuery == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Halaman yang anda tuju tidak dapat ditemukan.",
		})
		return
	}

	article, err := models.GetArticleBySlug(slugQuery)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Halaman yang anda tuju tidak dapat ditemukan.",
		})
		return
	}

	userContext := getUserFromContext(ctx.Request)
	user, err := models.GetUserByID(userContext.Credential.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Autentikasi diperlukan.",
		})
		return
	}

	// check is user is author for this article
	if user.ID != article.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Anda tidak mempunyai ijin untuk mengakses metode ini.",
		})
		return
	}

	payloads := struct {
		Tags        string                `form:"tags"`
		Title       string                `form:"title" validate:"max=50"`
		Logo        *multipart.FileHeader `form:"logo"`
		Description string                `form:"description" validate:"max=255"`
		Content     string                `form:"content"`
		Status      string                `form:"status" validate:"max=9"`
	}{}
	if err := ctx.ShouldBindWith(&payloads, binding.FormMultipart); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload yang anda gunakan tidak di ijinkan.",
		})
		return
	}

	if payloads.Tags != "" {
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
		article.Tags = tags
	}
	if payloads.Title != "" {
		article.Title = cases.Title(language.Indonesian).String(payloads.Title)
		slugString := slug.MakeLang(payloads.Title, "id")
		article.Slug = slugString
	}

	var newFile string
	if payloads.Logo != nil {

		oldFile := article.Logo
		_ = os.RemoveAll("." + oldFile)

		newFile = fmt.Sprintf("/media/articles/%s/%s", article.Slug, payloads.Logo.Filename)
	}

	if payloads.Description != "" {
		article.Description = payloads.Description
	}
	if payloads.Content != "" {
		article.Content = payloads.Content
	}
	if payloads.Status != "" {
		article.Status = payloads.Status
	}

	if newFile != "" {
		article.Logo = newFile

		if err := ctx.SaveUploadedFile(payloads.Logo, "."+newFile); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}

	models.GetDB().Save(&article)

	ctx.JSON(http.StatusOK, article)
}
