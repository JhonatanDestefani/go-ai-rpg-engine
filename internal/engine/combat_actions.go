package engine

import (
	"fmt"
	"go-inventory-management/internal/domain"
)

func (g *GameEngine) executePlayerBasicAttack(
	combat *domain.CombatState,
) string {
	damage, attackRoll, defenseRoll := BasicAttackContestedDamage(
		combat.Player.Stats,
		combat.Monster.Stats,
	)

	combat.Monster.Stats.HP -= damage

	if combat.Monster.Stats.HP < 0 {
		combat.Monster.Stats.HP = 0
	}

	return fmt.Sprintf(
		"%s\nYou dealt 💥 %d damage .",
		FormatContestedRolls(attackRoll, defenseRoll),
		damage,
	)
}

func (g *GameEngine) executePlayerSkill(
	combat *domain.CombatState,
	skillID string,
) string {
	skill, found := FindSkillByID(combat.Player.Skills, skillID)
	if !found {
		return "Skill not found."
	}

	if !SpendMP(combat.Player, skill.MPCost) {
		return "Not enough MP."
	}

	switch skill.ID {
	case "heal":
		return g.executeHealSkill(combat, *skill)
	default:
		return g.executeDamageSkill(combat, *skill)
	}
}

func (g *GameEngine) executePlayerItem(combat *domain.CombatState, itemID string) string {
	return UseConsumableItem(
		&combat.Player.Inventory,
		&combat.Player.Stats,
		itemID,
	)
}

func (g *GameEngine) executeMonsterBasicAttack(
	combat *domain.CombatState,
) string {
	damage, attackRoll, defenseRoll := BasicAttackContestedDamage(
		combat.Monster.Stats,
		combat.Player.Stats,
	)

	combat.Player.Stats.HP -= damage

	if combat.Player.Stats.HP < 0 {
		combat.Player.Stats.HP = 0
	}

	return fmt.Sprintf(
		"%s\n%s dealt 💥 %d damage.",
		FormatContestedRolls(attackRoll, defenseRoll),
		combat.Monster.Name,
		damage,
	)
}

func (g *GameEngine) executeMonsterSkill(combat *domain.CombatState) string {
	if len(combat.Monster.Skills) == 0 {
	return g.executeMonsterBasicAttack(combat)
	}

	index := RollDice(1, len(combat.Monster.Skills)).Total - 1
	skill := combat.Monster.Skills[index]

	if combat.Monster.Stats.MP < skill.MPCost {
	return g.executeMonsterBasicAttack(combat)
	}

	combat.Monster.Stats.MP -= skill.MPCost

	damage, attackRoll, defenseRoll := ResolveContestedDamage(combat.Monster.Stats, combat.Player.Stats, skill)

	combat.Player.Stats.HP -= damage
	ApplySkillEffects(combat, skill, false)

	return fmt.Sprintf(
		"%s\n%s used %s and dealt 💥 %d damage.",
		FormatContestedRolls(attackRoll, defenseRoll),
		combat.Monster.Name,
		skill.Name,
		damage,
	)
}

func (g *GameEngine) executeDamageSkill(
	combat *domain.CombatState,
	skill domain.Skill,
) string {
	damage, attackRoll, defenseRoll := ResolveContestedDamage(
		combat.Player.Stats,
		combat.Monster.Stats,
		skill,
	)

	combat.Monster.Stats.HP -= damage
	if combat.Monster.Stats.HP < 0 {
		combat.Monster.Stats.HP = 0
	}

	ApplySkillEffects(combat, skill, true)

	return fmt.Sprintf(
		"%s\nYou used %s and dealt 💥 %d damage.",
		FormatContestedRolls(attackRoll, defenseRoll),
		skill.Name,
		damage,
	)
}
func (g *GameEngine) executeHealSkill(
	combat *domain.CombatState,
	skill domain.Skill,
) string {
	hpBefore := combat.Player.Stats.HP

	healAmount := skill.Power + combat.Player.Stats.Magic

	combat.Player.Stats.HP += healAmount
	if combat.Player.Stats.HP > combat.Player.Stats.MaxHP {
		combat.Player.Stats.HP = combat.Player.Stats.MaxHP
	}

	actualHeal := combat.Player.Stats.HP - hpBefore

	return fmt.Sprintf(
		"You used %s and restored %d HP.",
		skill.Name,
		actualHeal,
	)
}
