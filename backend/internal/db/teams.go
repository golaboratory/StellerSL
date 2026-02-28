package db

import (
	"context"
)

type Team struct {
	ID        string
	TenantID  string
	Name      string
	CreatedAt string
	UpdatedAt string
}

type TeamMember struct {
	TeamID string
	UserID string
	Role   string
}

func (q *Queries) CreateTeam(ctx context.Context, tenantID, name string) (*Team, error) {
	row := q.db.QueryRowContext(ctx, `
		INSERT INTO teams (tenant_id, name)
		VALUES ($1, $2)
		RETURNING id, tenant_id, name, created_at, updated_at`,
		tenantID, name)
	var t Team
	err := row.Scan(&t.ID, &t.TenantID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

func (q *Queries) ListTeams(ctx context.Context, tenantID string) ([]Team, error) {
	rows, err := q.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, created_at, updated_at
		FROM teams
		WHERE tenant_id = $1
		ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.TenantID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func (q *Queries) AddTeamMember(ctx context.Context, teamID, userID, role string) error {
	_, err := q.db.ExecContext(ctx, `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (team_id, user_id) DO UPDATE SET role = EXCLUDED.role`,
		teamID, userID, role)
	return err
}

func (q *Queries) ListTeamMembers(ctx context.Context, teamID string) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, `
		SELECT u.id, u.tenant_id, u.email, u.password_hash, u.name
		FROM users u
		JOIN team_members tm ON u.id = tm.user_id
		WHERE tm.team_id = $1`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
