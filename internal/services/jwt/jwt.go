package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/config"
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

func (s *jwtService) GenerateToken(userID string, permissions []domain.Permission) (domain.TokenString, error) {
	serial, err := s.blacklist.NewSerial()
	if err != nil {
		return "", fmt.Errorf("jwtService.GenerateToken: s.blacklist.NewSerial returned error: %w", err)
	}

	exp := time.Now().Add(s.cfg.TokenExpiry)
	// Создание JWT токена с использованием приватного ключа RSA
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":                  userID,
		"exp":                  exp.Unix(),  // срок действия токена
		s.cfg.PermissionsClaim: permissions, // права доступа
		s.cfg.SerialClaim:      serial,
	})

	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("jwtService.GenerateToken: token.SignedString returned error: %w", err)
	}

	return domain.TokenString(signedToken), nil
}

func (s *jwtService) toToken(token *jwt.Token) (*domain.Token, error) {
	var ok bool
	result := domain.Token{}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("toToken: invalid token claims")
	}

	strSubject, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("toToken: claim not found: sub")
	}
	result.Subject = strSubject

	// Приведение к типу int64
	serialFloat, ok := claims[s.cfg.SerialClaim].(float64)
	if !ok {
		return nil, errors.New(fmt.Sprintf("toToken: claim not found: %s", s.cfg.SerialClaim))
	}
	result.Serial = int64(serialFloat)

	// Приведение времени истечения токена
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("toToken: claim not found: exp")
	}
	result.Expire = time.Unix(int64(expFloat), 0)

	// Дальнейшая обработка прав доступа и других полей
	permIf, ok := claims[s.cfg.PermissionsClaim].([]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("toToken: claim not found: %s", s.cfg.PermissionsClaim))
	}

	var permissions []domain.Permission
	for _, perm := range permIf {
		permStr, ok := perm.(string)
		if !ok {
			return nil, errors.New("toToken: invalid permission type")
		}
		permissions = append(permissions, domain.Permission(permStr))
	}
	result.Permissions = permissions

	return &result, nil
}

func (s *jwtService) ValidateToken(tokenString domain.TokenString) (*domain.Token, error) {

	// Валидация токена с использованием публичного ключа RSA
	jwtToken, err := jwt.Parse(string(tokenString), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.publicKey, nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, err
	}

	token, err := s.toToken(jwtToken)
	if err != nil {
		return nil, fmt.Errorf("jwtService.ValidateToken: s.toToken returned error: %w", err)
	}

	blacklisted, err := s.blacklist.IsBlacklisted(token.Serial)
	if err != nil {
		return nil, fmt.Errorf("jwtService.ValidateToken: s.blacklist.IsBlacklisted returned error: %w", err)
	}
	if blacklisted {
		return nil, errors.New("jwtService.ValidateToken: token is blacklisted")
	}

	return token, nil
}

func (s *jwtService) RevokeToken(tokenString domain.TokenString) error {
	var err error

	err = errors.New("Metod not implemented")
	// if err = s.blacklist.AddToBlacklist(serial); err != nil {

	// }
	return err
}

func parsePermissions(permStr string) ([]domain.Permission, error) {
	if permStr == "" {
		return nil, fmt.Errorf("permission string is empty")
	}

	// Переменная для хранения распарсенных данных
	var permissions []domain.Permission

	// Парсим JSON-массив
	err := json.Unmarshal([]byte(permStr), &permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse permissions: %w", err)
	}

	return permissions, nil
}
