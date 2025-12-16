package auth

import (
	"crypto/rsa"
	_ "crypto/x509"
	_ "encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSigningKey *rsa.PrivateKey
	jwtVerifyKey  *rsa.PublicKey
	jwtExpiration = time.Hour * 24
)

// JwtVerifyKey возвращает публичный ключ для проверки JWT
func JwtVerifyKey() *rsa.PublicKey {
	return jwtVerifyKey
}

func InitJWT(privateKeyPath, publicKeyPath string) error {
	// 1. Чтение приватного ключа
	signBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %w", err) // 💡 ОШИБКА ЗДЕСЬ
	}

	// 2. Парсинг приватного ключа
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err) // 💡 ОШИБКА ЗДЕСЬ
	}
	jwtSigningKey = signingKey // <--- Эта переменная должна быть инициализирована!

	// 3. Чтение и парсинг публичного ключа (для верификации)
	verifyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err) // 💡 ОШИБКА ЗДЕСЬ
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err) // 💡 ОШИБКА ЗДЕСЬ
	}
	jwtVerifyKey = verifyKey

	return nil
}
func GenerateJWT(userID uint, isModerator bool) (string, error) {
	if jwtSigningKey == nil {
		return "", errors.New("jwt not initialized")
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"mod": isModerator,
		"exp": time.Now().Add(jwtExpiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(jwtSigningKey)
}
