package engine

import (
	"go-inventory-management/internal/domain"
)

func (g *GameEngine) ExecutePlayerTurn(combat *domain.CombatState, action domain.CombatActionType, actionID string) string {
	if !combat.IsPlayerTurn || combat.IsFinished {
		return "It is not the player's turn."
	}

	effectMessage := ApplyCombatEffects(&combat.PlayerEffects, &combat.Player.Stats)
	if combat.Player.Stats.HP <= 0 {
		combat.Player.Stats.HP = 0
		combat.IsFinished = true
		combat.PlayerWon = false
		return "You were defeated."
	}

	var message string

	switch action {
	case domain.ActionBasicAttack:
		message = g.executePlayerBasicAttack(combat)

	case domain.ActionSkill:
		message = g.executePlayerSkill(combat, actionID)

	case domain.ActionItem:
		message = g.executePlayerItem(combat, actionID)

	default:
		return "Invalid action."
	}

	if combat.Monster.Stats.HP <= 0 {
		return message + "\n" + g.FinishCombatVictory(combat)
	}

	EndPlayerTurn(combat)
	return effectMessage + message + "\nMonster turn."
}

func (g *GameEngine) ExecuteMonsterTurn(combat *domain.CombatState) string {
	if combat.IsPlayerTurn || combat.IsFinished {
		return "It is not the monster's turn."
	}

	var message string

	effectMessage := ApplyCombatEffects(&combat.MonsterEffects, &combat.Monster.Stats)

	if combat.Monster.Stats.HP <= 0 {
		return effectMessage + g.FinishCombatVictory(combat)
	}

	action := RollDice(1, 2).Total

	if action == 1 || len(combat.Monster.Skills) == 0 {
		message = g.executeMonsterBasicAttack(combat)
	} else {
		message = g.executeMonsterSkill(combat)
	}

	if combat.Player.Stats.HP <= 0 {
		return message + "\n" + g.FinishCombatDefeat(combat)
	}

	EndMonsterTurn(combat)

	return effectMessage + message + "\nYour turn."
}
