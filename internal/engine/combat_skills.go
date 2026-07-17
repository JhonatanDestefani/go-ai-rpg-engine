package engine

import (
	"go-inventory-management/internal/domain"
)

func SpendMP(player *domain.Player, amount int) bool {
	if player.Stats.MP < amount {
		return false
	}
	player.Stats.MP -= amount
	return true
}

func FindSkillByID(skills []domain.Skill, skillID string) (*domain.Skill, bool) {
	for i := range skills {
		if skills[i].ID == skillID {
			return &skills[i], true
		}
	}

	return nil, false
}

func ApplySkillEffects(
	combat *domain.CombatState,
	skill domain.Skill,
	casterIsPlayer bool,
) {
	for _, effect := range skill.Effects {
		if effect.IsBuff {
			if casterIsPlayer {
				AddOrRefreshEffect(
					&combat.PlayerEffects,
					effect,
				)
			} else {
				AddOrRefreshEffect(
					&combat.MonsterEffects,
					effect,
				)
			}

			continue
		}

		if effect.IsDebuff {
			if casterIsPlayer {
				AddOrRefreshEffect(
					&combat.MonsterEffects,
					effect,
				)
			} else {
				AddOrRefreshEffect(
					&combat.PlayerEffects,
					effect,
				)
			}
		}
	}
}

func AddOrRefreshEffect(
	effects *[]domain.CombatEffect,
	newEffect domain.CombatEffect,
) {
	for i := range *effects {
		if (*effects)[i].Type == newEffect.Type {
			(*effects)[i].Duration = newEffect.Duration

			if newEffect.Power > (*effects)[i].Power {
				(*effects)[i].Power = newEffect.Power
			}

			return
		}
	}

	*effects = append(*effects, newEffect)
}
