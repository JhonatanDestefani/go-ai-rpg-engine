package domain

type EffectType string

const (
	EffectDefending   EffectType = "defending"
	EffectBleeding    EffectType = "bleeding"
	EffectFrozen      EffectType = "frozen"
	EffectPoisoned    EffectType = "poisoned"
	EffectBurning     EffectType = "burning"
	EffectStunned     EffectType = "stunned"
	EffectAttackUp    EffectType = "attack_up"
	EffectDefenseUp   EffectType = "defense_up"
	EffectAttackDown  EffectType = "attack_down"
	EffectDefenseDown EffectType = "defense_down"
)

type CombatEffect struct {
	Type     EffectType
	Duration int
	Power    int
	IsDebuff bool
	IsBuff   bool
}
