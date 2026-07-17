package engine

import "go-inventory-management/internal/domain"

func (g *GameEngine) CreatePlayer( name string, class domain.ClassType) domain.Player {
	player := domain.Player{
		Name:         name,
		Class:        class,
		CurrentScene: "start",
		Level:        1,
		XP:           0,
		Gold:         0,
		Inventory:    domain.NewInventory(32),
	}

		starter := GetClassDefinition(class)

		player.Stats = starter.Stats
		player.Skills = starter.Skills

	for _, item := range starter.Items {
		player.Inventory.AddItem(item)
	}

	player.EquippedWeaponID = starter.EquippedWeaponID
	player.EquippedArmorID = starter.EquippedArmorID

	return player
}

func GetStarterSkills(class domain.ClassType) []domain.Skill {
	switch class {

	case domain.ClassWarrior:
		return []domain.Skill{
			{
				ID:            "slash",
				Name:          "Slash",
				Description:   "A Slash that deals physical damage and make them bleed.",
				MPCost:        2,
				Power:         6,
				DamageType:    domain.DamagePhysical,
				RequiredLevel: 1,
				Effects: []domain.CombatEffect{
					{
						Type:     domain.EffectBleeding,
						Duration: 3,
						Power:    5,
						IsDebuff: true,
						IsBuff:   false,
					},
				},
			},
		}

	case domain.ClassMage:
		return []domain.Skill{
			{
				ID:            "fireball",
				Name:          "Fireball",
				Description:   "A Fireball that deals magical damage and burns the enemy.",
				MPCost:        6,
				Power:         10,
				DamageType:    domain.DamageMagical,
				RequiredLevel: 1,
				Effects: []domain.CombatEffect{
					{
						Type:     domain.EffectBurning,
						Duration: 3,
						Power:    5,
						IsDebuff: true,
						IsBuff:   false,
					},
				},
			},
		}
	}
	return nil
}
