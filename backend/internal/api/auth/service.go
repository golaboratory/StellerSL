package auth

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/user/stellersl/backend/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	conn    *sql.DB
	queries *db.Queries
	jwtKey  []byte
}

func NewService(conn *sql.DB, queries *db.Queries, jwtKey []byte) *Service {
	return &Service{conn: conn, queries: queries, jwtKey: jwtKey}
}

type AuthClaims struct {
	jwt.RegisteredClaims
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

func (s *Service) Register(ctx context.Context, input RegisterInput) error {
	return s.queries.WithTenant(ctx, s.conn, input.Body.TenantID, func(q *db.Queries) error {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Body.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user, err := q.CreateUser(ctx, input.Body.TenantID, input.Body.Email, string(hash), input.Body.Name)
		if err != nil {
			return err
		}

		if input.Body.InviteTeamID != "" {
			_ = q.AddTeamMember(ctx, input.Body.InviteTeamID, user.ID, "member")
		}

		return nil
	})
}

func (s *Service) Login(ctx context.Context, input LoginInput, tenantID string) (*LoginOutput, error) {
	var user *db.User
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		user, err = q.GetUserByEmail(ctx, input.Body.Email)
		return err
	})
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
