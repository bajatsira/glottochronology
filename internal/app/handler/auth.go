package handler

import (
	"context"
	"net/http"
	_ "time"

	"LAB1/internal/app/auth"
	"LAB1/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// ApiRegisterUser - POST /api/auth/register
func (h *Handler) ApiRegisterUser(c *gin.Context) {
	var body struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, err := auth.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}
	user := &ds.Users{
		Login:    body.Login,
		Password: hash,
	}
	if err := h.Repository.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "login": user.Login})
}

// ApiLogin - POST /api/auth/login
func (h *Handler) ApiLogin(c *gin.Context) {
	var body struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Repository.GetUserByLogin(body.Login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !auth.ComparePassword(user.Password, body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Создаем сессию в Redis
	sid, err := auth.CreateSession(context.Background(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session create failed"})
		return
	}

	// Установим cookie session_id
	c.SetCookie("session_id", sid, int(auth.SessionTTLSeconds()), "/", "", false, true) // HttpOnly. в prod: Secure=true

	// Опционально: создаем JWT и возвращаем в теле
	token, _ := auth.GenerateJWT(user.ID, user.IsLinguist)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"jwt":     token,
		"user": gin.H{
			"id":           user.ID,
			"login":        user.Login,
			"is_moderator": user.IsLinguist,
		},
	})
}

// ApiLogout - POST /api/auth/logout
func (h *Handler) ApiLogout(c *gin.Context) {
	// удалить session cookie и удалить в Redis
	cookie, err := c.Cookie("session_id")
	if err == nil && cookie != "" {
		_ = auth.DeleteSession(context.Background(), cookie)
		// удалить cookie у клиента: установить maxAge < 0
		c.SetCookie("session_id", "", -1, "/", "", false, true)
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
