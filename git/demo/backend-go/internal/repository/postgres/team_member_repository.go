package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type teamMemberRepository struct {
	db *sql.DB
}

func NewTeamMemberRepository(db *sql.DB) repository.TeamMemberRepository {
	return &teamMemberRepository{db: db}
}

func (r *teamMemberRepository) Create(ctx context.Context, ownerTwitchID, memberTwitchID, memberLogin string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO channel_team_members (owner_twitch_id, member_twitch_id, member_login, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (owner_twitch_id, member_twitch_id) DO NOTHING
	`, ownerTwitchID, memberTwitchID, memberLogin)
	if err != nil {
		return false, fmt.Errorf("fehler beim Einladen: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}

	return rowsAffected > 0, nil
}

func (r *teamMemberRepository) Delete(ctx context.Context, ownerTwitchID, memberTwitchID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM channel_team_members WHERE owner_twitch_id = $1 AND member_twitch_id = $2
	`, ownerTwitchID, memberTwitchID)
	if err != nil {
		return fmt.Errorf("fehler beim Entfernen des Team-Mitglieds: %w", err)
	}
	return nil
}

func (r *teamMemberRepository) ListByOwner(ctx context.Context, ownerTwitchID string) ([]*domain.TeamMember, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, owner_twitch_id, member_twitch_id, member_login, created_at
		FROM channel_team_members
		WHERE owner_twitch_id = $1
		ORDER BY created_at ASC
	`, ownerTwitchID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Team-Mitglieder: %w", err)
	}
	defer rows.Close()

	return scanTeamMembers(rows)
}

func (r *teamMemberRepository) ListByMember(ctx context.Context, memberTwitchID string) ([]*domain.TeamMember, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, owner_twitch_id, member_twitch_id, member_login, created_at
		FROM channel_team_members
		WHERE member_twitch_id = $1
		ORDER BY created_at ASC
	`, memberTwitchID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der verwaltbaren Kanäle: %w", err)
	}
	defer rows.Close()

	return scanTeamMembers(rows)
}

func (r *teamMemberRepository) Exists(ctx context.Context, ownerTwitchID, memberTwitchID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM channel_team_members WHERE owner_twitch_id = $1 AND member_twitch_id = $2)
	`, ownerTwitchID, memberTwitchID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("fehler beim Prüfen des Team-Zugriffs: %w", err)
	}
	return exists, nil
}

func scanTeamMembers(rows *sql.Rows) ([]*domain.TeamMember, error) {
	members := make([]*domain.TeamMember, 0)
	for rows.Next() {
		m := &domain.TeamMember{}
		if err := rows.Scan(&m.ID, &m.OwnerTwitchID, &m.MemberTwitchID, &m.MemberLogin, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("fehler beim Scannen des Team-Mitglieds: %w", err)
		}
		members = append(members, m)
	}
	return members, nil
}
