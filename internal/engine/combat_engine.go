package engine

import "go-inventory-management/internal/domain"

func (g *GameEngine) StartCombat(
	player *domain.Player,
	monster *domain.Monster,
) domain.CombatState {
	return domain.CombatState{
		Player:       player,
		Monster:      monster,
		IsPlayerTurn: true,
		IsActive:     true,
	}
}
