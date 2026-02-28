package gamification

import (
	"context"

	"github.com/user/stellersl/backend/internal/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) GetGrowth(ctx context.Context, userID string) (*GrowthOutput, error) {
	g, err := s.queries.GetUserGrowth(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &GrowthOutput{}
	resp.Body.Level = g.Level
	resp.Body.Exp = g.Exp
	resp.Body.CharacterType = g.CharacterType
	return resp, nil
}

func (s *Service) ListBadges(ctx context.Context, userID string) (*BadgeListOutput, error) {
	badges, err := s.queries.ListUserBadges(ctx, userID)
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
