package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aZ4ziL/blogs_api/auth"
	"github.com/aZ4ziL/blogs_api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var validate *validator.Validate

// UserHandlerGetToken is function to handle request to get new token.
func UserHandlerGetToken(ctx *gin.Context) {
	payloads := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := ctx.ShouldBindJSON(&payloads)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload yang anda gunakan tidak dijinkan.",
		})
		return
	}

	user, err := models.GetUserByUsername(payloads.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Username atau kata sandi yang anda masukkan salah.",
		})
		return
	}

	if !auth.DecryptionPassword(user.Password, payloads.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Username atau kata sandi yang anda masukkan salah.",
		})
		return
	}

	credential := auth.Credential{UserID: user.ID}
	token, err := auth.GetToken(credential)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"mesage": "Jangan bagikan akses token ini kepada siapapun.",
		"token":  token,
	})
}

// UserHandlerRegister is function to handler request to registration user.
func UserHandlerRegister(ctx *gin.Context) {
	payloads := struct {
		FirstName string `json:"first_name" validate:"required,max=50"`
		LastName  string `json:"last_name" validate:"required,max=50"`
		Username  string `json:"username" validate:"required,max=20"`
		Email     string `json:"email" validate:"required,email,max=50"`
		Password  string `json:"password" validate:"required,max=20"`
	}{}
	if err := ctx.ShouldBindJSON(&payloads); err != nil {
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

	user := models.User{
		FirstName: cases.Title(language.Indonesian).String(payloads.FirstName),
		LastName:  cases.Title(language.Indonesian).String(payloads.LastName),
		Username:  strings.TrimSpace(payloads.Username),
		Email:     strings.TrimSpace(payloads.Email),
		Password:  strings.TrimSpace(payloads.Password),
	}
	if err := models.CreateNewUser(&user); err != nil {
		// FIXME: Create response error if username or email is already taken by another user.
		if strings.Contains(err.Error(), "users_username_key") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Username yang anda gunakan telah dipakai oleh pengguna lain.",
			})
			return
		} else if strings.Contains(err.Error(), "users_email_key") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Alamat email yang anda gunakan telah terdaftar oleh pengguna lain.",
			})
			return
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, user)
}

// UserHandlerAuth is function check authentication from request context.
func UserHandlerAuth(ctx *gin.Context) {
	userContext := getUserFromContext(ctx.Request)
	user, err := models.GetUserByID(userContext.Credential.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"mesage": "Autentikasi dibutuhkan.",
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
