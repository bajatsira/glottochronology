package handler

import (
	"context"
	"net/http"
	"strings"

	"LAB1/internal/app/auth"
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CurrentUser struct {
	ID         uint
	Login      string
	IsLinguist bool
}

const CtxUserKey = "current_user"

func AuthMiddleware(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		// -------- 1) JWT Bearer --------
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

		// -------- 2) Cookie session --------
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
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if v, ok := c.Get(CtxUserKey); !ok || v == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}
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
