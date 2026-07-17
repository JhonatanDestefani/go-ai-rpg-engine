package domain

type CombatMode string

const (
	CombatForbidden CombatMode = "forbidden"
	CombatOptional  CombatMode = "optional"
	CombatRequired  CombatMode = "required"
)