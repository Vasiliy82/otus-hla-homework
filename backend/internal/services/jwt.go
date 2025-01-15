package services

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	cfg        *config.JWTConfig
	blacklist  domain.BlacklistRepository
}

// NewJWTService - конструктор для JWT сервиса
func NewJWTService(cfg *config.JWTConfig, blacklist domain.BlacklistRepository) (domain.JWTService, error) {
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

// GenerateToken создает JWT токен для пользователя
func (s *jwtService) GenerateToken(userID domain.UserKey, permissions []domain.Permission) (domain.TokenString, error) {
	serial, err := s.blacklist.NewSerial()
	if err != nil {
		return "", fmt.Errorf("jwtService.GenerateToken: s.blacklist.NewSerial returned error: %w", err)
	}

	claims := domain.UserClaims{
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        serial,
			Subject:   string(userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.TokenExpiry)),
			Issuer:    "myApp", // Это можно настроить через cfg, если нужно
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("jwtService.GenerateToken: token.SignedString returned error: %w", err)
	}

	return domain.TokenString(signedToken), nil
}

// ValidateToken проверяет JWT токен и возвращает *jwt.Token
func (s *jwtService) ValidateToken(tokenString domain.TokenString) (*jwt.Token, error) {
	jwtToken, err := jwt.ParseWithClaims(string(tokenString), &domain.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Извлекаем claims и проверяем blacklist
	claims, err := s.ExtractClaims(jwtToken)
	if err != nil {
		return nil, err
	}

	blacklisted, err := s.blacklist.IsBlacklisted(claims.ID)
	if err != nil {
		return nil, fmt.Errorf("jwtService.ValidateToken: s.blacklist.IsBlacklisted returned error: %w", err)
	}
	if blacklisted {
		return nil, fmt.Errorf("jwtService.ValidateToken: token is blacklisted")
	}

	return jwtToken, nil
}

// ExtractClaims извлекает кастомные claims из *jwt.Token
func (s *jwtService) ExtractClaims(token *jwt.Token) (*domain.UserClaims, error) {
	claims, ok := token.Claims.(*domain.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims structure: expected MyCustomClaims, got %T", token.Claims)
	}
	return claims, nil
}

// RevokeToken помещает токен в черный список
func (s *jwtService) RevokeToken(token *jwt.Token) error {
	claims, err := s.ExtractClaims(token)
	if err != nil {
		return err
	}

	if err := s.blacklist.AddToBlacklist(claims.ID, claims.ExpiresAt.Time); err != nil {
		return apperrors.NewInternalServerError("Internal server error", err)
	}
	return nil
}
