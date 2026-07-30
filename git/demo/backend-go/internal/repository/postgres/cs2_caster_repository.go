package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type CS2CasterRepository struct {
	db *sql.DB
}

func NewCS2CasterRepository(db *sql.DB) repository.CS2CasterRepository {
	return &CS2CasterRepository{db: db}
}

func generateGSIToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *CS2CasterRepository) GetOrCreateSettings(ctx context.Context, userTwitchID string) (*domain.CS2CasterSettings, error) {
	token, err := generateGSIToken()
	if err != nil {
		return nil, fmt.Errorf("fehler beim Generieren des GSI-Tokens: %w", err)
	}

	// INSERT ... ON CONFLICT DO NOTHING + anschliessendes SELECT deckt sowohl
	// den Erstzugriff als auch ein Rennen zweier gleichzeitiger Erstzugriffe
	// robust ab, ohne eine explizite Transaktion/Lock zu brauchen.
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO cs2_caster_settings (user_twitch_id, gsi_token)
		VALUES ($1, $2)
		ON CONFLICT (user_twitch_id) DO NOTHING
	`, userTwitchID, token)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Anlegen der CS2-Caster-Settings: %w", err)
	}

	return r.getByUserTwitchID(ctx, userTwitchID)
}

func (r *CS2CasterRepository) getByUserTwitchID(ctx context.Context, userTwitchID string) (*domain.CS2CasterSettings, error) {
	return r.scanSettings(ctx, `
		SELECT user_twitch_id, gsi_token, predictions_enabled, multikill_announce_enabled,
		       map_end_announce_enabled, title_update_enabled, created_at, updated_at
		FROM cs2_caster_settings
		WHERE user_twitch_id = $1
	`, userTwitchID)
}

func (r *CS2CasterRepository) GetByGSIToken(ctx context.Context, token string) (*domain.CS2CasterSettings, error) {
	return r.scanSettings(ctx, `
		SELECT user_twitch_id, gsi_token, predictions_enabled, multikill_announce_enabled,
		       map_end_announce_enabled, title_update_enabled, created_at, updated_at
		FROM cs2_caster_settings
		WHERE gsi_token = $1
	`, token)
}

func (r *CS2CasterRepository) scanSettings(ctx context.Context, query string, arg string) (*domain.CS2CasterSettings, error) {
	var s domain.CS2CasterSettings
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&s.UserTwitchID,
		&s.GSIToken,
		&s.PredictionsEnabled,
		&s.MultikillAnnounceEnabled,
		&s.MapEndAnnounceEnabled,
		&s.TitleUpdateEnabled,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der CS2-Caster-Settings: %w", err)
	}
	return &s, nil
}

func (r *CS2CasterRepository) UpdateSettings(ctx context.Context, userTwitchID string, input domain.CS2CasterSettingsUpdateInput) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE cs2_caster_settings
		SET predictions_enabled = COALESCE($2, predictions_enabled),
		    multikill_announce_enabled = COALESCE($3, multikill_announce_enabled),
		    map_end_announce_enabled = COALESCE($4, map_end_announce_enabled),
		    title_update_enabled = COALESCE($5, title_update_enabled),
		    updated_at = NOW()
		WHERE user_twitch_id = $1
	`,
		userTwitchID,
		input.PredictionsEnabled,
		input.MultikillAnnounceEnabled,
		input.MapEndAnnounceEnabled,
		input.TitleUpdateEnabled,
	)
	if err != nil {
		return fmt.Errorf("fehler beim Aktualisieren der CS2-Caster-Settings: %w", err)
	}
	return nil
}

func (r *CS2CasterRepository) ListNotes(ctx context.Context, userTwitchID string) ([]*domain.CS2Note, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_twitch_id, subject_type, subject_name, content, created_at, updated_at
		FROM cs2_notes
		WHERE user_twitch_id = $1
		ORDER BY subject_type, subject_name, created_at DESC
	`, userTwitchID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Notizen: %w", err)
	}
	defer rows.Close()

	notes := make([]*domain.CS2Note, 0)
	for rows.Next() {
		var n domain.CS2Note
		if err := rows.Scan(&n.ID, &n.UserTwitchID, &n.SubjectType, &n.SubjectName, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("fehler beim Scan der Notizen: %w", err)
		}
		notes = append(notes, &n)
	}
	return notes, rows.Err()
}

func (r *CS2CasterRepository) CreateNote(ctx context.Context, userTwitchID string, input domain.CS2NoteCreateInput) (*domain.CS2Note, error) {
	var n domain.CS2Note
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO cs2_notes (user_twitch_id, subject_type, subject_name, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_twitch_id, subject_type, subject_name, content, created_at, updated_at
	`, userTwitchID, input.SubjectType, input.SubjectName, input.Content).Scan(
		&n.ID, &n.UserTwitchID, &n.SubjectType, &n.SubjectName, &n.Content, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der Notiz: %w", err)
	}
	return &n, nil
}

func (r *CS2CasterRepository) UpdateNote(ctx context.Context, userTwitchID string, noteID int64, input domain.CS2NoteUpdateInput) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE cs2_notes
		SET subject_name = COALESCE($3, subject_name),
		    content = COALESCE($4, content),
		    updated_at = NOW()
		WHERE id = $1 AND user_twitch_id = $2
	`, noteID, userTwitchID, input.SubjectName, input.Content)
	if err != nil {
		return fmt.Errorf("fehler beim Aktualisieren der Notiz: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CS2CasterRepository) DeleteNote(ctx context.Context, userTwitchID string, noteID int64) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM cs2_notes WHERE id = $1 AND user_twitch_id = $2
	`, noteID, userTwitchID)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen der Notiz: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
