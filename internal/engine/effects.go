package engine

import (
	"fmt"
	"go-inventory-management/internal/domain"
)

var effectIcons = map[domain.EffectType]string{
	domain.EffectBleeding: "🩸",
	domain.EffectPoisoned: "☠️",
	domain.EffectBurning:  "🔥",
}

func ApplyCombatEffect(effects *[]domain.CombatEffect, effect domain.CombatEffect) {
	if effect.Duration <= 0 {
		return
	}

	*effects = append(*effects, effect)
}

func ApplyCombatEffects(
	effects *[]domain.CombatEffect,
	stats *domain.Stats,
) string {
	var message string
	activeEffects := []domain.CombatEffect{}

	for _, effect := range *effects {

		icon, ok := effectIcons[effect.Type]
		if ok {
			stats.HP -= effect.Power

			message += fmt.Sprintf(
				"%s %s dealt %d damage.\n",
				icon,
				effect.Type,
				effect.Power,
			)

			if stats.HP < 0 {
				stats.HP = 0
			}
		}

		effect.Duration--

		if effect.Duration > 0 {
			activeEffects = append(activeEffects, effect)
		}
	}

	*effects = activeEffects
	return message
}

func TickCombatEffects(effects *[]domain.CombatEffect) {
	activeEffects := []domain.CombatEffect{}

	for _, effect := range *effects {
		effect.Duration--

		if effect.Duration > 0 {
			activeEffects = append(activeEffects, effect)
		}
	}

	*effects = activeEffects
}
