package engine

import "go-inventory-management/internal/domain"

func (g *GameEngine) IsCombatOver(combat *domain.CombatState) bool {
	return combat.IsFinished ||
		combat.Player.Stats.HP <= 0 ||
		combat.Monster.Stats.HP <= 0
}

func (g *GameEngine) PlayerWon(combat *domain.CombatState) bool {
	return combat.PlayerWon ||
		(combat.Monster.Stats.HP <= 0 && combat.Player.Stats.HP > 0)
}
