/*package auth

import (
	"fmt"
	"time"



	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"LAB1/internal/app/config"
	"LAB1/internal/app/ds"
)

type JWTManager struct {
	config *config.Config
	secret []byte
}

func NewJWTManager(cfg *config.Config) (*JWTManager, error) {
	if cfg.JWT.Token == "" {
		return nil, fmt.Errorf("JWT secret (Token) is empty")
	}

	return &JWTManager{
		config: cfg,
		secret: []byte(cfg.JWT.Token),
	}, nil
}

// GenerateToken
func (m *JWTManager) GenerateToken(userUUID uuid.UUID, isLinguist bool) (string, error) {
	// Исправление: используем значение напрямую, так как оно уже time.Duration
	// Если это int64, приводим его: time.Duration(m.config.JWT.ExpiresIn)
	duration := m.config.JWT.ExpiresIn

	now := time.Now()

	claims := ds.JWTClaims{
		UserUUID:   userUUID,
		IsLinguist: isLinguist,
		RegisteredClaims: jwt.RegisteredClaims{
			// Прибавляем duration к текущему времени
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(duration))),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userUUID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken
func (m *JWTManager) ValidateToken(tokenString string) (*ds.JWTClaims, error) {
	claims := &ds.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, m.keyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (m *JWTManager) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return m.secret, nil
}

func GenerateJWT(userUUID uuid.UUID, isLinguist bool, cfg *config.Config) (string, error) {
	// 1. Подготовка данных (Claims)
	// Убедитесь, что тип userUUID совпадает с типом в ds.JWTClaims
	claims := &ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserUUID:   userUUID,
		IsLinguist: isLinguist,
	}

	// 2. Создание объекта токена с алгоритмом HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Подписание токена секретной строкой из конфига
	// Метод SignedString ожидает []byte для HS256
	tokenString, err := token.SignedString([]byte(cfg.JWT.Token))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании строкового токена: %w", err)
	}

	return tokenString, nil
}
*/

package auth

import (
	"fmt"
	"time"

	"LAB1/internal/app/config"
	"LAB1/internal/app/ds"
	"github.com/golang-jwt/jwt/v5"
)

var manager *JWTManager

type JWTManager struct {
	config *config.Config
}

func Init(cfg *config.Config) {
	manager = &JWTManager{config: cfg}
}
func GenerateJWT(userID uint, isLinguist bool) (string, error) {
	if manager == nil {
		return "", fmt.Errorf("JWT manager is not initialized")
	}

	cfg := manager.config
	claims := &ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:     userID,
		IsLinguist: isLinguist,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Token))
}

// ValidateToken проверяет строку токена и возвращает данные пользователя (claims)
func ValidateToken(tokenString string) (*ds.JWTClaims, error) {
	if manager == nil {
		return nil, fmt.Errorf("JWT manager is not initialized")
	}

	claims := &ds.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи (защита от подмены алгоритма)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Возвращаем тот же секрет, что и для генерации
		return []byte(manager.config.JWT.Token), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}
