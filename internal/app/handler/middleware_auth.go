package handler

import (
	"context"
	"net/http"
	"strings"

	"LAB1/internal/app/auth"
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type CurrentUser struct {
	ID         uint
	Login      string
	IsLinguist bool
}

const CtxUserKey = "user"

func AuthMiddleware(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Пытаемся авторизовать через JWT (Header)
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			tokenStr := strings.TrimSpace(authHeader[7:])

			// Используем нашу новую функцию
			claims, err := auth.ValidateToken(tokenStr)

			if err == nil {
				// Если токен валиден, берем UserID прямо из нашей структуры
				uid := claims.UserID
				if user, err := repo.GetUserByID(uid); err == nil {
					c.Set(CtxUserKey, &CurrentUser{
						ID:         user.ID,
						Login:      user.Login,
						IsLinguist: user.IsLinguist,
					})
					c.Next()
					return
				}
			}
		}

		// 2. Если JWT нет или он невалиден, пытаемся через сессию (Cookie)
		if sid, err := c.Cookie("session_id"); err == nil && sid != "" {
			uid, err := auth.GetUserIDBySession(context.Background(), sid)
			if err == nil && uid > 0 {
				if user, err := repo.GetUserByID(uid); err == nil {
					c.Set(CtxUserKey, &CurrentUser{
						ID:         user.ID,
						Login:      user.Login,
						IsLinguist: user.IsLinguist,
					})
					c.Next()
					return
				}
			}
		}

		// Если ничего не подошло, просто идем дальше (пользователь будет анонимным)
		c.Next()
	}
}

/*func AuthMiddleware(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			tokenStr := strings.TrimSpace(authHeader[7:])
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return auth.JwtVerifyKey(), nil
			})
			if err == nil && token != nil && token.Valid {
				if sub, ok := claims["sub"].(float64); ok {
					uid := uint(sub)
					if user, err := repo.GetUserByID(uid); err == nil {
						c.Set(CtxUserKey, &CurrentUser{
							ID:         user.ID,
							Login:      user.Login,
							IsLinguist: user.IsLinguist,
						})
						c.Next()
						return
					}
				}
			}
		}

		if sid, err := c.Cookie("session_id"); err == nil && sid != "" {
			uid, err := auth.GetUserIDBySession(context.Background(), sid)
			if err == nil && uid > 0 {
				if user, err := repo.GetUserByID(uid); err == nil {
					c.Set(CtxUserKey, &CurrentUser{
						ID:         user.ID,
						Login:      user.Login,
						IsLinguist: user.IsLinguist,
					})
					c.Next()
					return
				}
			}
		}

		c.Next()
	}
}*/

// RequireAuth — это middleware для проверки JWT токена
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		// 2. Проверяем, что заголовок имеет формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		// 3. Валидируем JWT токен
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 4. Сохраняем информацию о пользователе в контекст для следующих обработчиков
		c.Set(CtxUserKey, &CurrentUser{
			ID:         claims.UserID,
			IsLinguist: claims.IsLinguist,
		})

		// 5. Передаем управление дальше
		c.Next()
	}
}

func RequireLinguist() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(CtxUserKey)
		if !ok || v == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}
		u := v.(*CurrentUser)
		if !u.IsLinguist {
			c.JSON(http.StatusForbidden, gin.H{"error": "linguist role required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
