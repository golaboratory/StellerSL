package gamification

type GrowthOutput struct {
	Body struct {
		Level         int    `json:"level"`
		Exp           int    `json:"exp"`
		CharacterType string `json:"character_type"`
	}
}

type BadgeItem struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	IconSlug        string `json:"icon_slug"`
	RequirementType string `json:"requirement_type"`
}

type BadgeListOutput struct {
	Body struct {
		Items []BadgeItem `json:"items"`
	}
}
