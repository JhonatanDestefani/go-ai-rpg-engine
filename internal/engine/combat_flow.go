package engine

import (
	"fmt"
	"go-inventory-management/internal/domain"
)

func (g *GameEngine) FinishCombatVictory(combat *domain.CombatState) string {
	combat.Monster.Stats.HP = 0
	combat.IsFinished = true
	combat.PlayerWon = true

	drops, dropRoll, leveledUp, newSkills := g.GiveCombatRewards(
		combat.Player,
		combat.Monster,
	)

	message := fmt.Sprintf(
		"\nYou defeated %s!\n✨ +%dXP\n🪙 +%d Gold\n🎲 Drop roll: %d",
		combat.Monster.Name,
		combat.Monster.XPReward,
		combat.Monster.GoldReward,
		dropRoll,
	)

	if len(drops) > 0 {
		message += "\n\nDrops:"
		for _, item := range drops {
			message += fmt.Sprintf("\n- %s x%d", item.Name, item.Quantity)
		}
	}

	if leveledUp {
		message += fmt.Sprintf(
			"\n\n⬆️ You reached level %d!",
			combat.Player.Level,
		)
	}

	for _, skill := range newSkills {
		message += fmt.Sprintf(
			"\n⚔️ You learned a new skill: %s!",
			skill.Name,
		)
	}

	
	return message
}

func (g *GameEngine) FinishCombatDefeat(combat *domain.CombatState) string {
	combat.Player.Stats.HP = 0
	combat.IsFinished = true
	combat.PlayerWon = false

	return "You were defeated."
}

func EndPlayerTurn(combat *domain.CombatState) {
	TickCombatEffects(&combat.PlayerEffects)
	combat.IsPlayerTurn = false
}

func EndMonsterTurn(combat *domain.CombatState) {
	TickCombatEffects(&combat.MonsterEffects)
	combat.IsPlayerTurn = true
}
