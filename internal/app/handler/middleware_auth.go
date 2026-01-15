package handler

import (
	"LAB1/internal/app/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserKey = "user"

type CurrentUser struct {
	ID         uint
	IsLinguist bool
}

// RequireAuth (остается без изменений)
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			return
		}

		c.Set(CtxUserKey, &CurrentUser{
			ID:         claims.UserID,
			IsLinguist: claims.IsLinguist,
		})
		c.Next()
	}
}

// RequireLinguist - НОВАЯ, ПРАВИЛЬНАЯ, УПРОЩЕННАЯ ВЕРСИЯ
// Этот middleware предполагает, что RequireAuth уже отработал в цепочке до него.
func RequireLinguist() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Безопасно получаем пользователя из контекста.
		v, exists := c.Get(CtxUserKey)
		if !exists {
			// Этого не должно случиться, если RequireAuth был вызван до, но это хорошая защита.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user context not found, authentication may have failed"})
			return
		}

		// 2. Безопасно приводим тип.
		user, ok := v.(*CurrentUser)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid user type in context"})
			return
		}

		// 3. Выполняем главную задачу: проверку прав модератора.
		if !user.IsLinguist {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: linguist rights required"})
			return
		}

		// 4. Все в порядке, передаем управление дальше по цепочке.
		c.Next()
	}
}
