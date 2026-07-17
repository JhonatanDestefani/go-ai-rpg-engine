package engine

import (
	"context"

	"go-inventory-management/internal/domain"
)

func (g *GameEngine) SaveGame(
	ctx context.Context,
	player *domain.Player,
) error {
	return g.playerRepository.Update(ctx, player)
}