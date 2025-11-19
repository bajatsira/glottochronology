package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
	if privateKeyPath == "" || publicKeyPath == "" {
		return nil
	}
	// private
	privPem, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(privPem)
	if block == nil {
		return errors.New("invalid private key pem")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// try PKCS8
		k, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return err
		}
		if rsaKey, ok := k.(*rsa.PrivateKey); ok {
			privKey = rsaKey
		} else {
			return errors.New("private key is not RSA")
		}
	}
	jwtSigningKey = privKey

	// public
	pubPem, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}
	blockPub, _ := pem.Decode(pubPem)
	if blockPub == nil {
		return errors.New("invalid public key pem")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(blockPub.Bytes)
	if err != nil {
		// try ParsePKCS1PublicKey
		pubKey, err2 := x509.ParsePKCS1PublicKey(blockPub.Bytes)
		if err2 != nil {
			return err
		}
		jwtVerifyKey = pubKey
	} else {
		if rsaPub, ok := pubInterface.(*rsa.PublicKey); ok {
			jwtVerifyKey = rsaPub
		} else {
			return errors.New("public key is not RSA")
		}
	}
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
