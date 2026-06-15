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
	var g db.UserGrowth
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		uid := db.ParseUUID(userID)
		g, err = q.GetUserGrowth(ctx, uid)
		if err != nil {
			// Auto create if not exists
			_ = q.CreateUserGrowth(ctx, uid)
			g, err = q.GetUserGrowth(ctx, uid)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &GrowthOutput{}
	resp.Body.Level = int(g.Level)
	resp.Body.Exp = int(g.Exp)
	resp.Body.CharacterType = g.CharacterType
	return resp, nil
}

func (s *Service) ListBadges(ctx context.Context, tenantID, userID string) (*BadgeListOutput, error) {
	var badges []db.Badge
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		badges, err = q.ListUserBadges(ctx, db.ParseUUID(userID))
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &BadgeListOutput{}
	for _, b := range badges {
		resp.Body.Items = append(resp.Body.Items, BadgeItem{
			ID:              b.ID.String(),
			Name:            b.Name,
			Description:     b.Description.String,
			IconSlug:        b.IconSlug,
			RequirementType: b.RequirementType,
		})
	}
	return resp, nil
}
