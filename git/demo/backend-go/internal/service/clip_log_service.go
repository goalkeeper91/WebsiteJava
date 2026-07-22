package service

import (
	"context"
	"fmt"

	"demo/backend-go/internal/domain"
	"demo/backend-go/internal/repository"
)

type ClipLogService struct {
	clipRepo repository.ClipLogRepository
}

func NewClipLogService(clipRepo repository.ClipLogRepository) *ClipLogService {
	return &ClipLogService{clipRepo: clipRepo}
}

// GetClips lädt die Clip-Historie eines Users (paginiert).
func (s *ClipLogService) GetClips(ctx context.Context, userID string, limit, offset int) ([]*domain.ClipLog, error) {
	clips, err := s.clipRepo.GetByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Clips: %w", err)
	}
	return clips, nil
}

// GetStats liefert Aggregatszahlen (total/pending/completed/failed, success rate) für einen User.
func (s *ClipLogService) GetStats(ctx context.Context, userID string) (*repository.ClipLogStats, error) {
	stats, err := s.clipRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Clip-Statistiken: %w", err)
	}
	return stats, nil
}
