package domain

type DamageType string

const (
	DamagePhysical DamageType = "physical"
	DamageMagical  DamageType = "magical"
)

type SkillType string
const (
	SkillAttack  SkillType = "attack"
	SkillSupport SkillType = "support"
)


type Skill struct {
	ID            string
	Name          string
	MPCost        int
	Power         int
	SkillType SkillType
	DamageType    DamageType
	Description   string
	RequiredLevel int

	Effects []CombatEffect
}
