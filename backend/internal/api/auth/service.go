package auth

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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

		user, err := q.CreateUser(ctx, db.CreateUserParams{
			TenantID:     db.ParseUUID(input.Body.TenantID),
			Email:        input.Body.Email,
			PasswordHash: string(hash),
			Name:         input.Body.Name,
		})
		if err != nil {
			return err
		}

		if err := q.CreateUserGrowth(ctx, user.ID); err != nil {
			return err
		}

		if input.Body.InviteTeamID != "" {
			_ = q.AddTeamMember(ctx, db.AddTeamMemberParams{
				TeamID: db.ParseUUID(input.Body.InviteTeamID),
				UserID: user.ID,
				Role:   "member",
			})
		}

		return nil
	})
}

func (s *Service) Login(ctx context.Context, input LoginInput, tenantID string) (*LoginOutput, error) {
	var user db.User
	fmt.Println("tenantID", tenantID)
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		user, err = q.GetUserByEmail(ctx, input.Body.Email)
		return err
	})
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Body.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	claims := &AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))},
		TenantID:         user.TenantID.String(),
		UserID:           user.ID.String(),
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	resp := &LoginOutput{}
	resp.Body.Token = tokenString
	resp.Body.User.ID = user.ID.String()
	resp.Body.User.Name = user.Name
	resp.Body.User.Email = user.Email
	return resp, nil
}

func (s *Service) UpdateProfile(ctx context.Context, tenantID, userID string, input UpdateProfileInput) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		_, err := q.UpdateUser(ctx, db.UpdateUserParams{
			ID:        db.ParseUUID(userID),
			Name:      input.Body.Name,
			AvatarUrl: sql.NullString{String: input.Body.AvatarUrl, Valid: input.Body.AvatarUrl != ""},
		})
		return err
	})
}

func (s *Service) GetMe(ctx context.Context, tenantID, userID string) (*MeOutput, error) {
	var user db.User
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		user, err = q.GetUserByID(ctx, db.ParseUUID(userID))
		return err
	})
	if err != nil {
		return nil, err
	}
	resp := &MeOutput{}
	resp.Body.ID = user.ID.String()
	resp.Body.Email = user.Email
	resp.Body.Name = user.Name
	resp.Body.AvatarUrl = user.AvatarUrl.String
	return resp, nil
}

func (s *Service) ChangePassword(ctx context.Context, tenantID, userID string, input ChangePasswordInput) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		user, err := q.GetUserByID(ctx, db.ParseUUID(userID))
		if err != nil {
			return err
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Body.CurrentPassword)); err != nil {
			return fmt.Errorf("current password is incorrect")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Body.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		return q.UpdatePassword(ctx, db.UpdatePasswordParams{
			ID:           db.ParseUUID(userID),
			PasswordHash: string(hash),
		})
	})
}

func (s *Service) ResetPassword(ctx context.Context, tenantID string, input ResetPasswordInput) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Body.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		return q.UpdatePassword(ctx, db.UpdatePasswordParams{
			ID:           db.ParseUUID(input.Body.UserID),
			PasswordHash: string(hash),
		})
	})
}

func (s *Service) UpdateAvatar(ctx context.Context, tenantID, userID string, filename string, data []byte) (string, error) {
	// 1. Create uploads directory if not exists
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		_ = os.Mkdir(uploadDir, 0755)
	}

	// 2. Generate unique filename
	ext := "png"
	if strings.Contains(filename, ".") {
		parts := strings.Split(filename, ".")
		ext = parts[len(parts)-1]
	}
	newFilename := fmt.Sprintf("%s_%d.%s", userID, time.Now().Unix(), ext)
	filepath := fmt.Sprintf("%s/%s", uploadDir, newFilename)

	// 3. Save file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", err
	}

	url := fmt.Sprintf("/uploads/%s", newFilename)

	// 4. Update user record
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		user, err := q.GetUserByID(ctx, db.ParseUUID(userID))
		if err != nil {
			return err
		}
		_, err = q.UpdateUser(ctx, db.UpdateUserParams{
			ID:        user.ID,
			Name:      user.Name,
			AvatarUrl: sql.NullString{String: url, Valid: true},
		})
		return err
	})

	return url, err
}
