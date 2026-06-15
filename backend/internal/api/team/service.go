package team

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/user/stellersl/backend/internal/db"
)

type Service struct {
	conn    *sql.DB
	queries *db.Queries
}

func NewService(conn *sql.DB, queries *db.Queries) *Service {
	return &Service{conn: conn, queries: queries}
}

func (s *Service) List(ctx context.Context, tenantID string, limit, offset int32) (*TeamListOutput, error) {
	var teams []db.Team
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		teams, err = q.ListTeams(ctx, db.ListTeamsParams{
			Limit:  limit,
			Offset: offset,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &TeamListOutput{}
	for _, t := range teams {
		resp.Body.Items = append(resp.Body.Items, TeamItem{ID: t.ID.String(), Name: t.Name})
	}
	return resp, nil
}

func (s *Service) Create(ctx context.Context, tenantID, callerID, name string) (*TeamItem, error) {
	var t db.Team
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		t, err = q.CreateTeam(ctx, db.CreateTeamParams{
			TenantID: db.ParseUUID(tenantID),
			Name:     name,
		})
		if err != nil {
			return err
		}
		// The creator becomes the team owner; without this nobody could
		// pass the role checks on team management operations.
		return q.AddTeamMember(ctx, db.AddTeamMemberParams{
			TeamID: t.ID,
			UserID: db.ParseUUID(callerID),
			Role:   "owner",
		})
	})
	if err != nil {
		return nil, err
	}
	return &TeamItem{ID: t.ID.String(), Name: t.Name}, nil
}

func (s *Service) GetTeam(ctx context.Context, tenantID, teamID string) (*TeamItem, error) {
	var t db.Team
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		t, err = q.GetTeam(ctx, db.ParseUUID(teamID))
		return err
	})
	if err != nil {
		return nil, err
	}
	return &TeamItem{ID: t.ID.String(), Name: t.Name}, nil
}

func (s *Service) Update(ctx context.Context, tenantID, callerID, teamID, name string) (*TeamItem, error) {
	var t db.Team
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		if err := s.requireRole(ctx, q, teamID, callerID, "owner", "admin"); err != nil {
			return err
		}
		var err error
		t, err = q.UpdateTeam(ctx, db.UpdateTeamParams{
			ID:   db.ParseUUID(teamID),
			Name: name,
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return &TeamItem{ID: t.ID.String(), Name: t.Name}, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, callerID, teamID string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		if err := s.requireRole(ctx, q, teamID, callerID, "owner"); err != nil {
			return err
		}
		return q.DeleteTeam(ctx, db.ParseUUID(teamID))
	})
}

// requireRole checks if the caller has one of the allowed roles.
func (s *Service) requireRole(ctx context.Context, q *db.Queries, teamID, callerID string, allowedRoles ...string) error {
	role, err := q.GetTeamMemberRole(ctx, db.GetTeamMemberRoleParams{
		TeamID: db.ParseUUID(teamID),
		UserID: db.ParseUUID(callerID),
	})
	if err != nil {
		return fmt.Errorf("permission denied: not a team member")
	}
	for _, allowed := range allowedRoles {
		if role == allowed {
			return nil
		}
	}
	return fmt.Errorf("permission denied: requires role %v", allowedRoles)
}

// AddMember adds a member or, because AddTeamMember upserts, changes the role
// of an existing member.
func (s *Service) AddMember(ctx context.Context, tenantID, callerID, teamID, userID, role string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		if err := s.requireRole(ctx, q, teamID, callerID, "owner", "admin"); err != nil {
			return err
		}
		// Prevent demoting the owner via the upsert path
		targetRole, err := q.GetTeamMemberRole(ctx, db.GetTeamMemberRoleParams{
			TeamID: db.ParseUUID(teamID),
			UserID: db.ParseUUID(userID),
		})
		if err == nil && targetRole == "owner" && role != "owner" {
			return fmt.Errorf("cannot change the team owner's role")
		}
		return q.AddTeamMember(ctx, db.AddTeamMemberParams{
			TeamID: db.ParseUUID(teamID),
			UserID: db.ParseUUID(userID),
			Role:   role,
		})
	})
}

func (s *Service) RemoveMember(ctx context.Context, tenantID, callerID, teamID, userID string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		if err := s.requireRole(ctx, q, teamID, callerID, "owner", "admin"); err != nil {
			return err
		}
		// Prevent removing the owner
		targetRole, err := q.GetTeamMemberRole(ctx, db.GetTeamMemberRoleParams{
			TeamID: db.ParseUUID(teamID),
			UserID: db.ParseUUID(userID),
		})
		if err == nil && targetRole == "owner" {
			return fmt.Errorf("cannot remove the team owner")
		}
		return q.RemoveTeamMember(ctx, db.RemoveTeamMemberParams{
			TeamID: db.ParseUUID(teamID),
			UserID: db.ParseUUID(userID),
		})
	})
}

func (s *Service) ListMembers(ctx context.Context, tenantID, teamID string) (*TeamMemberListOutput, error) {
	var members []db.ListTeamMembersRow
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		members, err = q.ListTeamMembers(ctx, db.ParseUUID(teamID))
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &TeamMemberListOutput{}
	for _, m := range members {
		resp.Body.Items = append(resp.Body.Items, TeamMemberUser{
			ID:    m.ID.String(),
			Email: m.Email,
			Name:  m.Name,
			Role:  m.Role,
		})
	}
	return resp, nil
}
