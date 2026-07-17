package engine

import "go-inventory-management/internal/domain"

const (
	RegionDarkForest    = "Dark Forest"
	RegionAncientRuins  = "Ancient Ruins"
	RegionSanctuary = "Sanctuary"
)

func GetRegionForScene(sceneNumber int) string {
	switch {
	case sceneNumber >= 1 && sceneNumber <= 5:
		return RegionDarkForest

	case sceneNumber >= 6 && sceneNumber <= 10:
		return RegionAncientRuins

	case sceneNumber >= 11 && sceneNumber <= 15:
	    return RegionSanctuary

	default:
		return RegionSanctuary
	}
}

func GetCombatModeForScene(sceneNumber int) domain.CombatMode {
	switch {
	case sceneNumber == 1:
		return domain.CombatForbidden

	case sceneNumber%3 == 0:
		return domain.CombatRequired

	default:
		return domain.CombatOptional
	}
}