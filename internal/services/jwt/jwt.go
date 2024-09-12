package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type BlacklistRepository interface {
	NewSerial() (int64, error)
	AddToBlacklist(serial int64, expireDate time.Time) error
	IsBlacklisted(serial int64) (bool, error)
}

type jwtService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	cfg        *config.JWTConfig
	blacklist  BlacklistRepository
}

// NewJWTService - конструктор для JWT сервиса
func NewJWTService(cfg *config.JWTConfig, blacklist BlacklistRepository) (domain.JWTService, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &jwtService{
		privateKey: privateKey,
		publicKey:  publicKey,
		cfg:        cfg,
		blacklist:  blacklist,
	}, nil
}

func (s *jwtService) GenerateToken(userID string, permissions []string) (string, error) {
	serial, err := s.blacklist.NewSerial()
	if err != nil {
		return "", fmt.Errorf("failed to generate serial: %w", err)
	}

	// Создание JWT токена с использованием приватного ключа RSA
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"serial":               serial,
		"sub":                  userID,
		"exp":                  time.Now().Add(s.cfg.TokenExpiry).Unix(), // срок действия токена
		s.cfg.PermissionsClaim: permissions,                              // права доступа
	})

	return token.SignedString(s.privateKey)
}

func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	serial, err := extractSerialFromToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to extract serial from token: %w", err)
	}

	blacklisted, err := s.blacklist.IsBlacklisted(serial)
	if err != nil {
		return nil, fmt.Errorf("failed to check blacklist: %w", err)
	}
	if blacklisted {
		return nil, errors.New("token is blacklisted")
	}

	// Валидация токена с использованием публичного ключа RSA
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return token, nil
}

// Вспомогательная функция для извлечения серийного номера из токена
func extractSerialFromToken(tokenString string) (int64, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	serial, ok := claims["serial"].(int64)
	if !ok {
		return 0, errors.New("serial not found in token")
	}

	return serial, nil
}
