package db

import (
	"context"
)

type UserGrowth struct {
	UserID        string
	Level         int
	Exp           int
	CharacterType string
	UpdatedAt     string
}

type Badge struct {
	ID              string
	Name            string
	Description     string
	IconSlug        string
	RequirementType string
}

func (q *Queries) GetUserGrowth(ctx context.Context, userID string) (*UserGrowth, error) {
	row := q.db.QueryRowContext(ctx, `
		SELECT user_id, level, exp, character_type, updated_at
		FROM user_growth
		WHERE user_id = $1`, userID)
	var g UserGrowth
	err := row.Scan(&g.UserID, &g.Level, &g.Exp, &g.CharacterType, &g.UpdatedAt)
	if err != nil {
		// If not found, create default
		_, err = q.db.ExecContext(ctx, "INSERT INTO user_growth (user_id) VALUES ($1) ON CONFLICT DO NOTHING", userID)
		return &UserGrowth{UserID: userID, Level: 1, Exp: 0, CharacterType: "default"}, nil
	}
	return &g, nil
}

func (q *Queries) AddExp(ctx context.Context, userID string, amount int) error {
	_, err := q.db.ExecContext(ctx, `
		UPDATE user_growth
		SET exp = exp + $2, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1`, userID, amount)
	if err != nil {
		return err
	}

	// Simple Level Up logic: level = exp / 100 + 1
	_, err = q.db.ExecContext(ctx, `
		UPDATE user_growth
		SET level = (exp / 100) + 1
		WHERE user_id = $1`, userID)
	return err
}

func (q *Queries) ListUserBadges(ctx context.Context, userID string) ([]Badge, error) {
	rows, err := q.db.QueryContext(ctx, `
		SELECT b.id, b.name, b.description, b.icon_slug, b.requirement_type
		FROM badges b
		JOIN user_badges ub ON b.id = ub.badge_id
		WHERE ub.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []Badge
	for rows.Next() {
		var b Badge
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.IconSlug, &b.RequirementType); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}
