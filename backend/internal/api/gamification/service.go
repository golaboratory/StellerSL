package gamification

import (
	"context"
	"database/sql"

	"github.com/user/stellersl/backend/internal/db"
)

type Service struct {
	conn    *sql.DB
	queries *db.Queries
}

func NewService(conn *sql.DB, queries *db.Queries) *Service {
	return &Service{conn: conn, queries: queries}
}

func (s *Service) GetGrowth(ctx context.Context, tenantID, userID string) (*GrowthOutput, error) {
	var g *db.UserGrowth
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		g, err = q.GetUserGrowth(ctx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &GrowthOutput{}
	resp.Body.Level = g.Level
	resp.Body.Exp = g.Exp
	resp.Body.CharacterType = g.CharacterType
	return resp, nil
}

func (s *Service) ListBadges(ctx context.Context, tenantID, userID string) (*BadgeListOutput, error) {
	var badges []db.Badge
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		badges, err = q.ListUserBadges(ctx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &BadgeListOutput{}
	for _, b := range badges {
		resp.Body.Items = append(resp.Body.Items, BadgeItem{
			ID:          b.ID,
			Name:        b.Name,
			Description: b.Description,
			IconSlug:    b.IconSlug,
		})
	}
	return resp, nil
}
