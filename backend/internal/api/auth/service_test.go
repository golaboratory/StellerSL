package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestTokenGeneration(t *testing.T) {
	// Setup
	jwtKey := []byte("test_secret_key")
	tenantID := "123e4567-e89b-12d3-a456-426614174000"
	userID := "987e6543-e21b-34d5-a654-426614174999"

	claims := &AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		TenantID: tenantID,
		UserID:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if tokenString == "" {
		t.Error("expected token string, got empty string")
	}

	// Verify
	parsedToken, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if claims, ok := parsedToken.Claims.(*AuthClaims); ok && parsedToken.Valid {
		if claims.TenantID != tenantID {
			t.Errorf("expected tenant_id %s, got %s", tenantID, claims.TenantID)
		}
		if claims.UserID != userID {
			t.Errorf("expected user_id %s, got %s", userID, claims.UserID)
		}
	} else {
		t.Error("invalid token claims")
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "my_secure_password"
	
	// Generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	hashStr := string(hash)
	
	if !strings.HasPrefix(hashStr, "$2a$") {
		t.Errorf("expected hash to start with $2a$, got %s", hashStr)
	}

	// Verify correct password
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		t.Errorf("expected password to match hash, but got error: %v", err)
	}

	// Verify incorrect password
	err = bcrypt.CompareHashAndPassword(hash, []byte("wrong_password"))
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}
