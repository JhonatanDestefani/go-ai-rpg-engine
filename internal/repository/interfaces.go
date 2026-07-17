package repository

import (
	"context"

	"go-inventory-management/internal/domain"
)

type PlayerRepository interface {
	Create(ctx context.Context, player *domain.Player) error
	GetByID(ctx context.Context, id int64) (*domain.Player, error)
	GetLatest(ctx context.Context) (*domain.Player, error)
	Update(ctx context.Context, player *domain.Player) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.Player, error)
}
