package db

import (
	"testing"

	"github.com/google/uuid"
)

func TestParseUUID(t *testing.T) {
	validUUIDStr := "123e4567-e89b-12d3-a456-426614174000"
	validUUID, _ := uuid.Parse(validUUIDStr)

	tests := []struct {
		name     string
		input    string
		expected uuid.UUID
	}{
		{"Valid UUID", validUUIDStr, validUUID},
		{"Invalid UUID", "invalid", uuid.UUID{}},
		{"Empty string", "", uuid.UUID{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseUUID(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestToNullUUID(t *testing.T) {
	validUUIDStr := "123e4567-e89b-12d3-a456-426614174000"
	validUUID, _ := uuid.Parse(validUUIDStr)

	tests := []struct {
		name     string
		input    string
		expected uuid.NullUUID
	}{
		{"Valid UUID", validUUIDStr, uuid.NullUUID{UUID: validUUID, Valid: true}},
		{"Invalid UUID", "invalid", uuid.NullUUID{Valid: false}},
		{"Empty string", "", uuid.NullUUID{Valid: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToNullUUID(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
