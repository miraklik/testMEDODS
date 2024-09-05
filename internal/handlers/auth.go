package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"main/internal/storage"
	"main/internal/tokens"
)

type AuthHandler struct {
	Storage *storage.Storage
}

func NewAuthHandler(storage *storage.Storage) *AuthHandler {
	return &AuthHandler{Storage: storage}
}

func (h *AuthHandler) TokenHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	ip := r.RemoteAddr

	accessToken, err := tokens.GenerateAccessToken(userID, ip)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := tokens.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	hashedToken, err := tokens.HashAndSalt(refreshToken)
	if err != nil {
		http.Error(w, "Failed to hash refresh token", http.StatusInternalServerError)
		return
	}

	err = h.Storage.StoreRefreshToken(userID, hashedToken)
	if err != nil {
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID := req["user_id"]
	refreshToken := req["refresh_token"]
	ip := r.RemoteAddr

	storedToken, err := h.Storage.GetRefreshToken(userID)
	if err != nil || storedToken == "" {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	err = tokens.CompareHash(storedToken, refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := tokens.GenerateAccessToken(userID, ip)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	if !strings.Contains(r.RemoteAddr, ip) {
		go func() {
			println("Warning: IP address has changed for user:", userID)
		}()
	}

	response := map[string]string{
		"access_token": newAccessToken,
	}
	json.NewEncoder(w).Encode(response)
}
