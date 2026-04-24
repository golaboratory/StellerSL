package notification

import (
	"context"
	"database/sql"
	"time"

	"github.com/user/stellersl/backend/internal/db"
)

type Service struct {
	conn    *sql.DB
	queries *db.Queries
}

func NewService(conn *sql.DB, queries *db.Queries) *Service {
	return &Service{conn: conn, queries: queries}
}

func (s *Service) List(ctx context.Context, tenantID, userID string, limit, offset int32) (*NotificationListOutput, error) {
	var items []db.Notification
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		items, err = q.ListNotifications(ctx, db.ListNotificationsParams{
			UserID: db.ParseUUID(userID),
			Limit:  limit,
			Offset: offset,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	resp := &NotificationListOutput{}
	for _, n := range items {
		item := NotificationItem{
			ID:        n.ID.String(),
			Type:      n.Type,
			Title:     n.Title,
			CreatedAt: n.CreatedAt.Time.Format(time.RFC3339),
		}
		if n.Message.Valid {
			item.Message = &n.Message.String
		}
		if n.ReadAt.Valid {
			s := n.ReadAt.Time.Format(time.RFC3339)
			item.ReadAt = &s
		}
		resp.Body.Items = append(resp.Body.Items, item)
	}
	return resp, nil
}

func (s *Service) UnreadCount(ctx context.Context, tenantID, userID string) (*UnreadCountOutput, error) {
	var count int32
	err := s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		var err error
		count, err = q.CountUnreadNotifications(ctx, db.ParseUUID(userID))
		return err
	})
	if err != nil {
		return nil, err
	}
	resp := &UnreadCountOutput{}
	resp.Body.Count = count
	return resp, nil
}

func (s *Service) MarkRead(ctx context.Context, tenantID, userID, notifID string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		return q.MarkNotificationRead(ctx, db.MarkNotificationReadParams{
			ID:     db.ParseUUID(notifID),
			UserID: db.ParseUUID(userID),
		})
	})
}

func (s *Service) MarkAllRead(ctx context.Context, tenantID, userID string) error {
	return s.queries.WithTenant(ctx, s.conn, tenantID, func(q *db.Queries) error {
		return q.MarkAllNotificationsRead(ctx, db.ParseUUID(userID))
	})
}
