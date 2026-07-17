package engine

import (
	"go-inventory-management/internal/domain"
)

func CalculateDamage(attacker domain.Stats, defender domain.Stats, skill domain.Skill) int {
	var damage int

	switch skill.DamageType {
	case domain.DamagePhysical:
		damage = attacker.Attack + skill.Power - defender.Defense

	case domain.DamageMagical:
		damage = attacker.Magic + skill.Power - defender.Magic

	default:
		damage = attacker.Attack + skill.Power - defender.Defense
	}

	if damage < 1 {
		return 1
	}

	return damage
}

func BasicAttackDamage(attacker domain.Stats, defender domain.Stats) int {

	damage := attacker.Attack - defender.Defense

	if damage < 1 {
		return 1
	}

	return damage
}

func ResolveContestedDamage(
	attacker domain.Stats,
	defender domain.Stats,
	skill domain.Skill,
) (int, DiceRoll, DiceRoll) {
	attackRoll := Roll2D6()
	defenseRoll := Roll2D6()

	damage := CalculateDamage(attacker, defender, skill)

	if attackRoll.Total < defenseRoll.Total {
		damage /= 2
	}

	if damage < 1 {
		damage = 1
	}

	return damage, attackRoll, defenseRoll
}

func BasicAttackContestedDamage(
	attacker domain.Stats,
	defender domain.Stats,
) (int, DiceRoll, DiceRoll) {
	attackRoll := Roll2D6()
	defenseRoll := Roll2D6()

	damage := BasicAttackDamage(attacker, defender)

	if attackRoll.Total < defenseRoll.Total {
		damage /= 2
	}

	if damage < 1 {
		damage = 1
	}

	return damage, attackRoll, defenseRoll
}
