package handler

import (
	"context"
	_ "net/http"
	"strings"

	"LAB1/internal/app/auth"
	_ "LAB1/internal/app/ds"
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// структура текущего пользователя
type CurrentUser struct {
	ID          uint
	Login       string
	IsModerator bool
}

// ключ в контексте
const CtxUserKey = "current_user"

// AuthMiddleware использует *repository.Repository напрямую
func AuthMiddleware(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		// --------- 1) Пытаемся прочитать JWT ----------
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			tokenStr := strings.TrimSpace(authHeader[7:])
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return auth.JwtVerifyKey(), nil
			})

			if err == nil && token != nil && token.Valid {
				sub, ok := claims["sub"].(float64)
				if ok {
					uid := uint(sub)
					user, err := repo.GetUserByID(uid)
					if err == nil {
						c.Set(CtxUserKey, &CurrentUser{
							ID:          user.ID,
							Login:       user.Login,
							IsModerator: user.IsLinguist,
						})
						c.Next()
						return
					}
				}
			}
		}

		// --------- 2) Пытаемся прочитать cookie session_id ----------
		if sid, err := c.Cookie("session_id"); err == nil && sid != "" {
			uid, err := auth.GetUserIDBySession(context.Background(), sid)
			if err == nil && uid > 0 {
				user, err := repo.GetUserByID(uid)
				if err == nil {
					c.Set(CtxUserKey, &CurrentUser{
						ID:          user.ID,
						Login:       user.Login,
						IsModerator: user.IsLinguist,
					})
					c.Next()
					return
				}
			}
		}

		// --------- 3) Пользователь гость ----------
		c.Next()
	}
}
