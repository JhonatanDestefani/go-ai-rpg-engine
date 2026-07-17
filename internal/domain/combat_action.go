package domain

type CombatActionType string

const (
	ActionBasicAttack CombatActionType = "basic_attack"
	ActionSkill       CombatActionType = "skill"
	ActionItem        CombatActionType = "item"
)

type CombatAction struct {
	Action CombatActionType

	SkillID string
	ItemID  string
}
