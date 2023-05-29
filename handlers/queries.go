package handlers

import (
	"net/http"
	"strconv"
)

// get query string by key, if key is undefined return default value.
func getQueryString(r *http.Request, key, defaultValue string) string {
	result := r.URL.Query().Get(key)
	if result == "" {
		return defaultValue
	}

	return result
}

// get query int by key
func getQueryInt(r *http.Request, key string, defaultValue int) int {
	result := r.URL.Query().Get(key)
	if result == "" {
		return defaultValue
	}

	resultInt, err := strconv.Atoi(result)
	if err != nil {
		return defaultValue
	}

	return resultInt
}
