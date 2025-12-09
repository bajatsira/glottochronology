package handler

import (
	"context"
	"net/http"

	"LAB1/internal/app/auth"
	"LAB1/internal/app/ds"

	"github.com/gin-gonic/gin"
)

/*
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
		if err != nil || !auth.ComparePassword(user.Password, body.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		sid, err := auth.CreateSession(context.Background(), user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "session create failed"})
			return
		}

		c.SetCookie("session_id", sid, auth.SessionTTLSeconds(), "/", "", false, true)

		token, _ := auth.GenerateJWT(user.ID, user.IsLinguist)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"jwt":     token,
			"user": gin.H{
				"id":          user.ID,
				"login":       user.Login,
				"is_linguist": user.IsLinguist,
			},
		})
	}

	func (h *Handler) ApiLogout(c *gin.Context) {
		if sid, err := c.Cookie("session_id"); err == nil {
			_ = auth.DeleteSession(context.Background(), sid)
		}
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "logout ok"})
	}
*/
func (h *Handler) ApiRegisterUser(c *gin.Context) {
	var body struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// plain password
	user := &ds.Users{
		Login:    body.Login,
		Password: body.Password,
	}

	if err := h.Repository.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"login": user.Login,
	})
}

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

	// plain password check
	if user.Password != body.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Create session (как в оригинале)
	sid, err := auth.CreateSession(context.Background(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session create failed"})
		return
	}

	// Set cookie
	c.SetCookie(
		"session_id",
		sid,
		auth.SessionTTLSeconds(),
		"/",
		"",
		false,
		true,
	)

	token, _ := auth.GenerateJWT(user.ID, user.IsLinguist)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"jwt":     token,
		"user": gin.H{
			"id":          user.ID,
			"login":       user.Login,
			"is_linguist": user.IsLinguist,
		},
	})
}

func (h *Handler) ApiLogout(c *gin.Context) {
	if sid, err := c.Cookie("session_id"); err == nil {
		_ = auth.DeleteSession(context.Background(), sid)
	}

	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout ok"})
}
