package engine

import (
	"go-inventory-management/internal/domain"
	"fmt"
)
func GetAvailableMonstersForRegion(
	region string,
	sceneNumber int,
) []domain.AvailableMonster {
	fmt.Printf(
		"GET MONSTERS -> region=%q sceneNumber=%d sanctuary=%q\n",
		region,
		sceneNumber,
		RegionSanctuary,
	)
	switch region {
	case RegionDarkForest:
		return []domain.AvailableMonster{
			{
				ID:          "goblin",
				Name:        "Goblin",
				Description: "A small and aggressive forest raider armed with crude weapons.",
			},
			{
				ID:          "orc",
				Name:        "Orc",
				Description: "A brutal warrior roaming the deeper parts of the forest.",
			},
		}

	case RegionAncientRuins:
		return []domain.AvailableMonster{
			{
				ID:          "orc",
				Name:        "Orc",
				Description: "A heavily armed orc searching the ancient ruins.",
			},
			{
				ID:          "guardian",
				Name:        "Guardian",
				Description: "An ancient guardian awakened within the ruined structures.",
			},
		}

	case RegionSanctuary:
		if sceneNumber == 15 {
			return []domain.AvailableMonster{
				{
					ID:   "dragon",
					Name: "Dragon",
					Description: "An enormous ancient dragon whose flames consume " +
						"everything that enters the sanctuary.",
				},
			}
		}

		if sceneNumber >= 11 && sceneNumber <= 14 {
			return []domain.AvailableMonster{
				{
					ID:   "dragon_statue",
					Name: "Dragon Statue",
					Description: "An ancient stone statue infused with draconic magic. " +
						"As intruders approach, glowing runes ignite across its body and awakens. ",
				},
			}
		}

		return nil

	default:
		return nil
	}
}