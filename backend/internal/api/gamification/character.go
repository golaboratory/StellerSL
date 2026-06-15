package gamification

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/stellersl/backend/internal/db"
)

// CharacterTypeForLevel returns the character type based on level thresholds.
// Lv.1-4: egg, Lv.5-9: chick, Lv.10-19: chicken, Lv.20+: phoenix
func CharacterTypeForLevel(level int32) string {
	switch {
	case level >= 20:
		return "phoenix"
	case level >= 10:
		return "chicken"
	case level >= 5:
		return "chick"
	default:
		return "egg"
	}
}

// UpdateCharacterType checks the user's level and updates character_type if needed.
func UpdateCharacterType(ctx context.Context, queries db.Querier, userID uuid.UUID) error {
	growth, err := queries.GetUserGrowth(ctx, userID)
	if err != nil {
		return err
	}

	expected := CharacterTypeForLevel(growth.Level)
	if growth.CharacterType != expected {
		return queries.UpdateCharacterType(ctx, db.UpdateCharacterTypeParams{
			UserID:        userID,
			CharacterType: expected,
		})
	}
	return nil
}
