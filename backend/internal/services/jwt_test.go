package services_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain/mocks"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/services"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	test_private_key = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDcxhNwFwJtXSdI
QE+6Ioe82HyYTZpn6KadsNI5lzcdpgvj9ngPFUbqGtSJarMc93qznwk2tuZkKC8k
cOOkrJMBwZ9bLK8W+euRMHgk9vZEuGB+/v6n8N3mvW+vsn6g0fg7LKa2gILhs1uK
eiPXfneqL7bCk3xjU7/Q8CDxlE9OLEwkICphD37ll2y4PRcrchzhO579+M6YrBMI
G9DTPWt2DvueZQyXZdWudio9UmzgIqHLjoeyROlSIO1KX10U3CGv2UTT6YsUoeyp
8wjozLSKXR7RlzZ53NXrCXdLwE6V6NTeQEqXdKb4yPS5oDyqPFTi39Hfx7heGfBo
lf3FNJ4bAgMBAAECggEAYoyep6X1xOjUvKlMjYuVaOSANaJKfwC4w2JnbSLFjRwO
abufGyiFx8GjRxYUjyUfpiejPsPFM0dGx+8Ghv8r/hg2sMXRAKIeF+j5cJK3GrTt
CjN8bG4WN8YvIVA9uz8PHicP4h6ajfJ4tedQsYR4GUWEQPYCC/qaAMP4CK6J+hvy
v2vMMbzN1PHBKXIeax5qWHpXMvOVGMOWdsZ8Rc7UbgY80dKYNs+ahgLOxiVHxsXs
aoFJnIe81p24RLHtabH1N8cGhZOyINjnPvoeNE/HW7LZFJ/J5k9kuIrGtOAAKZjh
bit2qtnR5fVO0pVPsdrxYVTL/2M5Hd3oLAf4mFeF3QKBgQDsOtUWDKeDxma4Twkq
h43xJtiWfIFatQVELwho5NodEfh9ZfL06wlxbsJ8t2jJpcQKQ3zDWej3MKi9W4Qu
FVnaUfj1sL50EVH2Syxh5HPRrXX4R/hni0jdJVjCahHHE2u1dYnoziogDbPdMWc/
xzEqAj6/ocSBjb7dHoM/DXG+FQKBgQDvQBnJAufcKer/10xmylf+AZGKyPHqSw2+
9SzPIY6o8OiTG6x/0ldHOLGnZeAiivQq4v9A76YFvPhz9UKSv00VXxKr86Uu+gKO
c4VbVzwe7/Fed++m/PP3uQKo8JY/+/nBbfJCvFlL0l9rLo+TEQYAoXAsX/a0b2EZ
lOrgErYnbwKBgH7Ef48KkWZstLjZaQDSp4AuqXHwNHZZyA6z8p5fmRCakS+x4vQ9
oN6nYmUNA4WamB4t4yjt+c+U5ChhkQgt2v8GmEQ4aavdk49I/fM2ZlSx8imfbZUb
MKnEHeKOiyW6rUU+Yxh0cjSrRcdAeLjICwERHV020T34s+DzO9k9PLmVAoGBAKbZ
BSJxrFCVyxTwiI+GvSae4WjwCgViogtx3/XzaRHYL9mnivz5K3S3zOz41v4/+VeP
RoN6nUWTK5FykSLV1mP5EYRpPeEs6Wt+lJnGlF7e5m0DJ1ZFQb6Yf4phfebRSrPi
gPiZcYy3AWQ17FqbnJwD+b54jgv3QLgeak4pvm5xAoGAUmiDI0Jbqi+5UdxgiLxT
pOXxK31rp4OBsLXCF2pMJteWGF4nqRjhawB5si8Qp6AVVlfVK05CFmuIAIDZzUMb
Kd3u0fmLDiDKxWHieyfKirJ5lF0FD194zaY0ndn1gR2AztbUkEWDLK6heVA39AJS
K1YJybLpwnAmqgy1hTfLLMg=
-----END PRIVATE KEY-----`
	test_public_key = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3MYTcBcCbV0nSEBPuiKH
vNh8mE2aZ+imnbDSOZc3HaYL4/Z4DxVG6hrUiWqzHPd6s58JNrbmZCgvJHDjpKyT
AcGfWyyvFvnrkTB4JPb2RLhgfv7+p/Dd5r1vr7J+oNH4OyymtoCC4bNbinoj1353
qi+2wpN8Y1O/0PAg8ZRPTixMJCAqYQ9+5ZdsuD0XK3Ic4Tue/fjOmKwTCBvQ0z1r
dg77nmUMl2XVrnYqPVJs4CKhy46HskTpUiDtSl9dFNwhr9lE0+mLFKHsqfMI6My0
il0e0Zc2edzV6wl3S8BOlejU3kBKl3Sm+Mj0uaA8qjxU4t/R38e4XhnwaJX9xTSe
GwIDAQAB
-----END PUBLIC KEY-----`
	test_token = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYzNTg4NDcsInBlcm1pc3Npb25zIjpbIlVTRVJfR0VUIl0sInNlcmlhbCI6NDAzLCJzdWIiOiIyNjMzOWMyYS03Zjc3LTQ3ODctYWY2Yy1kYWYxNGQ2MzdkYzIifQ.IRma1ab0cWIAHt0TqJzoM1DHtZB80KhiDN3dqiO7kpNCZ1G-GQ24zC1llEmsNq88EH0Q8nHxGfbVRolqOqM27w-eJaUVBQHTOQQZMiP-Bd3loz6-6239hq6_eWIwDB4Gtryh6LFIv3h8cCN0HJgCCxNkXHnpgkpZ_xSes7pJvOIf8Ha-3udcgRQzPz49UKG0pMhyW1m_99xKNLwJKoDv2M-aLQSLZObOQ8J1ALY0mMhykxf_gTMuMvDaUEZwUZUNwwF4A_FHFdGV67TF_EZ6aOlv7BLE7DV3cVxrzD-GVzO3haflq2SGgNamhrn1BFuGrj2KhzDgMakiiC97O6Lx-g`
)

func TestJWTService_GenerateToken_Success(t *testing.T) {
	mockBL := mocks.NewBlacklistRepository(t)
	swtService, err := services.NewJWTService(&config.JWTConfig{PrivateKey: test_private_key, PublicKey: test_public_key, TokenExpiry: 86400 * time.Second}, mockBL)
	if err != nil {
		t.Fatalf("Error creating JWTService: %v", err)
	}
	mockBL.On("NewSerial").Return("123", nil)
	token, err := swtService.GenerateToken("26339c2a-7f77-4787-af6c-daf14d637dc2", []domain.Permission{domain.PermissionUserGet})
	assert.NotEqual(t, domain.TokenString(""), token)
	assert.NoError(t, err)
	mockBL.AssertExpectations(t)
}
func TestJWTService_GenerateToken_DatabaseError(t *testing.T) {
	mockBL := mocks.NewBlacklistRepository(t)
	swtService, err := services.NewJWTService(&config.JWTConfig{PrivateKey: test_private_key, PublicKey: test_public_key, TokenExpiry: 86400 * time.Second}, mockBL)
	if err != nil {
		t.Fatalf("Error creating JWTService: %v", err)
	}
	mockBL.On("NewSerial").Return("", errors.New("Database error"))
	token, err := swtService.GenerateToken("26339c2a-7f77-4787-af6c-daf14d637dc2", []domain.Permission{domain.PermissionUserGet})
	assert.Equal(t, domain.TokenString(""), token)
	assert.Error(t, err)
	mockBL.AssertExpectations(t)
}

func TestJWTService_ValidateToken_Success(t *testing.T) {
	// Мокаем time.Now() с помощью gomonkey
	patches := gomonkey.ApplyFunc(time.Now, func() time.Time {
		return time.Unix(1726358840, 0) // замена текущего времени на нужное значение
	})
	defer patches.Reset() // восстанавливаем оригинальное поведение после теста

	fmt.Printf("Now is %s\n", time.Now().Format("2006-01-02 15:04:05"))

	mockBL := mocks.NewBlacklistRepository(t)
	swtService, err := services.NewJWTService(&config.JWTConfig{PrivateKey: test_private_key, PublicKey: test_public_key, TokenExpiry: 86400 * time.Second}, mockBL)
	if err != nil {
		t.Fatalf("Error creating JWTService: %v", err)
	}
	mockBL.On("IsBlacklisted", mock.Anything).Return(false, nil)
	token, err := swtService.ValidateToken(test_token)
	assert.NotNil(t, token)
	assert.NoError(t, err)
	mockBL.AssertExpectations(t)
}

func TestJWTService_ValidateToken_Expired(t *testing.T) {
	// Мокаем time.Now() с помощью gomonkey
	patches := gomonkey.ApplyFunc(time.Now, func() time.Time {
		return time.Unix(1826358848, 0) // замена текущего времени на нужное значение
	})
	defer patches.Reset() // восстанавливаем оригинальное поведение после теста

	fmt.Printf("Now is %s\n", time.Now().Format("2006-01-02 15:04:05"))

	mockBL := mocks.NewBlacklistRepository(t)
	swtService, err := services.NewJWTService(&config.JWTConfig{PrivateKey: test_private_key, PublicKey: test_public_key, TokenExpiry: 86400 * time.Second}, mockBL)
	if err != nil {
		t.Fatalf("Error creating JWTService: %v", err)
	}
	// mockBL.On("IsBlacklisted", mock.Anything).Return(false, nil)
	token, err := swtService.ValidateToken(test_token)
	assert.Nil(t, token)
	assert.Error(t, err)
	mockBL.AssertExpectations(t)
}

func TestJWTService_ValidateToken_Blacklisted(t *testing.T) {
	// Мокаем time.Now() с помощью gomonkey
	patches := gomonkey.ApplyFunc(time.Now, func() time.Time {
		return time.Unix(1726358840, 0) // замена текущего времени на нужное значение
	})
	defer patches.Reset() // восстанавливаем оригинальное поведение после теста

	fmt.Printf("Now is %s\n", time.Now().Format("2006-01-02 15:04:05"))

	mockBL := mocks.NewBlacklistRepository(t)
	swtService, err := services.NewJWTService(&config.JWTConfig{PrivateKey: test_private_key, PublicKey: test_public_key, TokenExpiry: 86400 * time.Second}, mockBL)
	if err != nil {
		t.Fatalf("Error creating JWTService: %v", err)
	}
	mockBL.On("IsBlacklisted", mock.Anything).Return(true, nil)
	token, err := swtService.ValidateToken(test_token)
	assert.Nil(t, token)
	assert.Error(t, err)
	mockBL.AssertExpectations(t)
}

func TestJWTService_RevokeToken_Success(t *testing.T) {
	// Мокаем time.Now() с помощью gomonkey
	patches := gomonkey.ApplyFunc(time.Now, func() time.Time {
		return time.Unix(1726358840, 0) // замена текущего времени на нужное значение
	})
	defer patches.Reset() // восстанавливаем оригинальное поведение после теста

	fmt.Printf("Now is %s\n", time.Now().Format("2006-01-02 15:04:05"))

	mockBL := mocks.NewBlacklistRepository(t)
	swtService, err := services.NewJWTService(&config.JWTConfig{PrivateKey: test_private_key, PublicKey: test_public_key, TokenExpiry: 86400 * time.Second}, mockBL)
	if err != nil {
		t.Fatalf("Error creating JWTService: %v", err)
	}
	mockBL.On("AddToBlacklist", mock.Anything, mock.Anything).Return(nil)

	mockClaims := &domain.UserClaims{
		Permissions: []domain.Permission{domain.PermissionUserGet}, // Пример кастомных claims
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "123",                                             // ID пользователя
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Действителен на 1 час
			ID:        "400",                                             // jti
		},
	}
	mockToken := jwt.Token{Claims: mockClaims}

	err = swtService.RevokeToken(&mockToken)
	assert.NoError(t, err)
	mockBL.AssertExpectations(t)
}

func TestJWTService_RevokeToken_Error(t *testing.T) {
	// Мокаем time.Now() с помощью gomonkey
	patches := gomonkey.ApplyFunc(time.Now, func() time.Time {
		return time.Unix(1726358840, 0) // замена текущего времени на нужное значение
	})
	defer patches.Reset() // восстанавливаем оригинальное поведение после теста

	fmt.Printf("Now is %s\n", time.Now().Format("2006-01-02 15:04:05"))

	mockBL := mocks.NewBlacklistRepository(t)
	swtService, err := services.NewJWTService(&config.JWTConfig{PrivateKey: test_private_key, PublicKey: test_public_key, TokenExpiry: 86400 * time.Second}, mockBL)
	if err != nil {
		t.Fatalf("Error creating JWTService: %v", err)
	}
	mockBL.On("AddToBlacklist", mock.Anything, mock.Anything).Return(errors.New("Database error"))

	mockClaims := &domain.UserClaims{
		Permissions: []domain.Permission{domain.PermissionUserGet}, // Пример кастомных claims
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "123",                                             // ID пользователя
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Действителен на 1 час
			ID:        "400",                                             // jti
		},
	}
	mockToken := jwt.Token{Claims: mockClaims}

	err = swtService.RevokeToken(&mockToken)
	assert.Error(t, err)
	mockBL.AssertExpectations(t)
}
