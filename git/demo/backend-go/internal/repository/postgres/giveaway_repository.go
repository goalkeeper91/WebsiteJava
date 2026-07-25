package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type giveawayRepository struct {
	db *sql.DB
}

func NewGiveawayRepository(db *sql.DB) repository.GiveawayRepository {
	return &giveawayRepository{db: db}
}

const giveawayColumns = `
	id, user_twitch_id, status, keyword, sub_bonus, winner_twitch_id, winner_login,
	started_at, ended_at, created_at, updated_at
`

func scanGiveaway(row interface{ Scan(dest ...interface{}) error }) (*domain.Giveaway, error) {
	g := &domain.Giveaway{}
	var winnerTwitchID, winnerLogin sql.NullString
	var endedAt sql.NullTime

	err := row.Scan(
		&g.ID, &g.UserTwitchID, &g.Status, &g.Keyword, &g.SubBonus, &winnerTwitchID, &winnerLogin,
		&g.StartedAt, &endedAt, &g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if winnerTwitchID.Valid {
		g.WinnerTwitchID = &winnerTwitchID.String
	}
	if winnerLogin.Valid {
		g.WinnerLogin = &winnerLogin.String
	}
	if endedAt.Valid {
		g.EndedAt = &endedAt.Time
	}

	return g, nil
}

func (r *giveawayRepository) CreateGiveaway(ctx context.Context, g *domain.Giveaway) error {
	query := `
		INSERT INTO giveaways (user_twitch_id, status, keyword, sub_bonus, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx, query, g.UserTwitchID, g.Status, g.Keyword, g.SubBonus, g.StartedAt, g.CreatedAt, g.UpdatedAt,
	).Scan(&g.ID)
	if err != nil {
		return fmt.Errorf("fehler beim Erstellen des Giveaways: %w", err)
	}

	return nil
}

func (r *giveawayRepository) GetOpenGiveaway(ctx context.Context, userTwitchID string) (*domain.Giveaway, error) {
	query := `
		SELECT ` + giveawayColumns + `
		FROM giveaways
		WHERE user_twitch_id = $1 AND status = 'open'
		ORDER BY started_at DESC
		LIMIT 1
	`

	g, err := scanGiveaway(r.db.QueryRowContext(ctx, query, userTwitchID))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden des offenen Giveaways: %w", err)
	}

	return g, nil
}

// GetAllOpenGiveaways returns every currently-open giveaway across all
// channels - backs the bot's startup cache warmup (a bot restart mid-giveaway
// otherwise wouldn't know which keywords are currently live).
func (r *giveawayRepository) GetAllOpenGiveaways(ctx context.Context) ([]*domain.Giveaway, error) {
	query := `SELECT ` + giveawayColumns + ` FROM giveaways WHERE status = 'open'`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden offener Giveaways: %w", err)
	}
	defer rows.Close()

	giveaways := make([]*domain.Giveaway, 0)
	for rows.Next() {
		g, err := scanGiveaway(rows)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen offener Giveaways: %w", err)
		}
		giveaways = append(giveaways, g)
	}

	return giveaways, nil
}

func (r *giveawayRepository) AddEntry(ctx context.Context, giveawayID int64, viewerTwitchID, viewerLogin string, entries int) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO giveaway_entries (giveaway_id, viewer_twitch_id, viewer_login, entries, entered_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (giveaway_id, viewer_twitch_id) DO NOTHING
	`, giveawayID, viewerTwitchID, viewerLogin, entries)
	if err != nil {
		return false, fmt.Errorf("fehler beim Eintragen ins Giveaway: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}

	return rowsAffected > 0, nil
}

func (r *giveawayRepository) GetEntries(ctx context.Context, giveawayID int64) ([]*domain.GiveawayEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, giveaway_id, viewer_twitch_id, viewer_login, entries, entered_at
		FROM giveaway_entries
		WHERE giveaway_id = $1
	`, giveawayID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Giveaway-Teilnehmer: %w", err)
	}
	defer rows.Close()

	entries := make([]*domain.GiveawayEntry, 0)
	for rows.Next() {
		e := &domain.GiveawayEntry{}
		if err := rows.Scan(&e.ID, &e.GiveawayID, &e.ViewerTwitchID, &e.ViewerLogin, &e.Entries, &e.EnteredAt); err != nil {
			return nil, fmt.Errorf("fehler beim Scannen des Giveaway-Teilnehmers: %w", err)
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (r *giveawayRepository) GetEntryCount(ctx context.Context, giveawayID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM giveaway_entries WHERE giveaway_id = $1", giveawayID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Zählen der Giveaway-Teilnehmer: %w", err)
	}
	return count, nil
}

func (r *giveawayRepository) CloseGiveaway(ctx context.Context, giveawayID int64, winnerTwitchID, winnerLogin string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE giveaways
		SET status = 'closed', winner_twitch_id = $2, winner_login = $3, ended_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, giveawayID, winnerTwitchID, winnerLogin)
	if err != nil {
		return fmt.Errorf("fehler beim Schließen des Giveaways: %w", err)
	}
	return nil
}

func (r *giveawayRepository) CancelGiveaway(ctx context.Context, giveawayID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE giveaways
		SET status = 'closed', ended_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, giveawayID)
	if err != nil {
		return fmt.Errorf("fehler beim Abbrechen des Giveaways: %w", err)
	}
	return nil
}

func (r *giveawayRepository) GetHistory(ctx context.Context, userTwitchID string, limit, offset int) ([]*domain.Giveaway, int64, error) {
	var total int64
	err := r.db.QueryRowContext(
		ctx, "SELECT COUNT(*) FROM giveaways WHERE user_twitch_id = $1", userTwitchID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Zählen der Giveaways: %w", err)
	}

	query := `
		SELECT ` + giveawayColumns + `,
		       (SELECT COUNT(*) FROM giveaway_entries WHERE giveaway_entries.giveaway_id = giveaways.id) AS entry_count
		FROM giveaways
		WHERE user_twitch_id = $1
		ORDER BY started_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userTwitchID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler beim Laden der Giveaway-Historie: %w", err)
	}
	defer rows.Close()

	giveaways := make([]*domain.Giveaway, 0)
	for rows.Next() {
		g, err := scanGiveawayWithEntryCount(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("fehler beim Scannen des Giveaways: %w", err)
		}
		giveaways = append(giveaways, g)
	}

	return giveaways, total, nil
}

func scanGiveawayWithEntryCount(row interface{ Scan(dest ...interface{}) error }) (*domain.Giveaway, error) {
	g := &domain.Giveaway{}
	var winnerTwitchID, winnerLogin sql.NullString
	var endedAt sql.NullTime

	err := row.Scan(
		&g.ID, &g.UserTwitchID, &g.Status, &g.Keyword, &g.SubBonus, &winnerTwitchID, &winnerLogin,
		&g.StartedAt, &endedAt, &g.CreatedAt, &g.UpdatedAt, &g.EntryCount,
	)
	if err != nil {
		return nil, err
	}

	if winnerTwitchID.Valid {
		g.WinnerTwitchID = &winnerTwitchID.String
	}
	if winnerLogin.Valid {
		g.WinnerLogin = &winnerLogin.String
	}
	if endedAt.Valid {
		g.EndedAt = &endedAt.Time
	}

	return g, nil
}
