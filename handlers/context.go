package handlers

import (
	"net/http"

	"github.com/aZ4ziL/blogs_api/auth"
)

// getUserFromContext return user info
func getUserFromContext(r *http.Request) auth.Claims {
	userContext := r.Context().Value(&auth.UserAuth{}).(auth.Claims)

	return userContext
}
