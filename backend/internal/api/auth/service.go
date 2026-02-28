package auth

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/user/stellersl/backend/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries *db.Queries
	jwtKey  []byte
}

func NewService(queries *db.Queries, jwtKey []byte) *Service {
	return &Service{queries: queries, jwtKey: jwtKey}
}

type AuthClaims struct {
	jwt.RegisteredClaims
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

func (s *Service) Register(ctx context.Context, input RegisterInput) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := s.queries.CreateUser(ctx, input.Body.TenantID, input.Body.Email, string(hash), input.Body.Name)
	if err != nil {
		return err
	}

	if input.Body.InviteTeamID != "" {
		_ = s.queries.AddTeamMember(ctx, input.Body.InviteTeamID, user.ID, "member")
	}

	return nil
}

func (s *Service) Login(ctx context.Context, input LoginInput, tenantID string) (*LoginOutput, error) {
	user, err := s.queries.GetUserByEmail(ctx, tenantID, input.Body.Email)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(user.PasswordHash, "$2a$") {
		if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Body.Password)); err != nil {
			return nil, err
		}
	} else if user.PasswordHash != input.Body.Password {
		return nil, err
	}

	claims := &AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))},
		TenantID:         user.TenantID,
		UserID:           user.ID,
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	resp := &LoginOutput{}
	resp.Body.Token = tokenString
	resp.Body.User.ID = user.ID
	resp.Body.User.Name = user.Name
	resp.Body.User.Email = user.Email
	return resp, nil
}
