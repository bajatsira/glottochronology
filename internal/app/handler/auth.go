package handler

import (
	"context"
	"net/http"
	_"fmt"
	_"time"
	

	"LAB1/internal/app/auth"
	"LAB1/internal/app/ds"

	"github.com/gin-gonic/gin"
	_"github.com/golang-jwt/jwt/v5"

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

// ApiLogin godoc
// @Summary Аутентификация пользователя
// @Description Вход в систему по логину и паролю. Возвращает JWT-токен и устанавливает сессионную куку.
// @Tags Аутентификация
// @Accept  json
// @Produce  json
// @Param   credentials  body   object{login=string,password=string} true "Учетные данные пользователя"
// @Success 200 {object} object "Успешный вход"
// @Failure 400 {object} object "Неверные данные"
// @Failure 401 {object} object "Неверный логин или пароль"
// @Router /auth/login [post]
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

	//token, _ := auth.GenerateJWT(user.ID, user.IsLinguist)

	token, err := auth.GenerateJWT(user.ID, user.IsLinguist)
	if err != nil {
		// Временно выведите ошибку, чтобы увидеть, что происходит с JWT:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT generation failed: " + err.Error()})
		return
	}
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


type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
/*
func (h *Handler) Login(gCtx *gin.Context) {
	cfg := h.Config
	req := &loginReq{}

	// Декодируем тело запроса
	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Ищем пользователя в БД по Email
	user, err := h.Repository.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":      "error",
				"description": "invalid email or password",
			})
		} else {
			gCtx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	// Проверяем email и хеш пароля (bcrypt)
	if req.Email == user.Email && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) == nil {
		
		// Создаем Claims на основе вашей структуры ds.JWTClaims
		claims := &ds.JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.ExpiresIn)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "bitop-admin", // можно заменить на ваше название приложения
			},
			UserUUID:   user.ID,         // Предполагаем, что user.ID имеет тип uuid.UUID
			IsLinguist: user.IsLinguist, // Заменили IsProfessor на IsLinguist
		}

		// Создаем токен (HS256)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		if token == nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
			return
		}

		// Подписываем токен секретной строкой из конфига
		strToken, err := token.SignedString([]byte(cfg.JWT.Token))
		if err != nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
			return
		}

		// Возвращаем ответ в требуемом формате
		gCtx.JSON(http.StatusOK, loginResp{
			ExpiresIn:   int64(cfg.JWT.ExpiresIn.Seconds()),
			AccessToken: strToken,
			TokenType:   "Bearer",
		})
		return
	}

	// Если пароль не подошел
	gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":      "error",
		"description": "invalid email or password",
	})
}*/

func (h *Handler) ApiLogout(c *gin.Context) {
	if sid, err := c.Cookie("session_id"); err == nil {
		_ = auth.DeleteSession(context.Background(), sid)
	}

	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout ok"})
}
