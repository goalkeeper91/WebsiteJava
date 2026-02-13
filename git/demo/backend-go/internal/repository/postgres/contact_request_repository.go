package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type contactRequestRepository struct {
	db *sql.DB
}

// NewContactRequestRepository erstellt ein neues PostgreSQL Contact Request Repository
func NewContactRequestRepository(db *sql.DB) repository.ContactRequestRepository {
	return &contactRequestRepository{db: db}
}

func (r *contactRequestRepository) Create(ctx context.Context, request *domain.ContactRequest) error {
	query := `
		INSERT INTO contact_requests (name, email, phone, subject, message, consent_given, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		request.Name,
		request.Email,
		request.Phone,
		request.Subject,
		request.Message,
		request.ConsentGiven,
		request.CreatedAt,
	).Scan(&request.ID)

	if err != nil {
		return fmt.Errorf("fehler beim Erstellen der Contact Request: %w", err)
	}

	return nil
}

func (r *contactRequestRepository) GetAll(ctx context.Context) ([]*domain.ContactRequest, error) {
	query := `
		SELECT id, name, email, phone, subject, message, consent_given, created_at
		FROM contact_requests
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Requests: %w", err)
	}
	defer rows.Close()

	requests := make([]*domain.ContactRequest, 0)
	for rows.Next() {
		request := &domain.ContactRequest{}
		err := rows.Scan(
			&request.ID,
			&request.Name,
			&request.Email,
			&request.Phone,
			&request.Subject,
			&request.Message,
			&request.ConsentGiven,
			&request.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen: %w", err)
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (r *contactRequestRepository) DeleteOlderThan(ctx context.Context, date time.Time) error {
	query := `DELETE FROM contact_requests WHERE created_at < $1`

	result, err := r.db.ExecContext(ctx, query, date)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen alter Requests: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen betroffener Zeilen: %w", err)
	}

	if rowsAffected > 0 {
		fmt.Printf("✅ %d alte Contact Requests gelöscht\n", rowsAffected)
	}

	return nil
}