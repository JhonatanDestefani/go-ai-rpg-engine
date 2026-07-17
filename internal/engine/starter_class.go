package engine

import "go-inventory-management/internal/domain"

type ClassTemplate struct {
	Stats             domain.Stats
	Skills            []domain.Skill
	Items             []domain.Item
	EquippedWeaponID  string
	EquippedArmorID   string
}
func GetClassDefinition(class domain.ClassType) ClassTemplate {
	switch class {

	case domain.ClassWarrior:
		return ClassTemplate{
			Stats: domain.Stats{
				MaxHP:   40,
				HP:      40,
				MaxMP:   6,
				MP:      6,
				Attack:  10,
				Defense: 5,
				Magic:   1,
			},

			Skills: GetStarterSkills(class),

			Items: []domain.Item{
			{	
			ID:          "rusty_sword",
			Name:        "Rusty Sword",
			Description: "An old sword.",
			Type:        domain.ItemWeapon,
			Rarity:      domain.RarityCommon,
			Value:       10,
			Quantity:    1,
			Attack:      5,
			},
			{
				ID:          "worn_armor",
				Name:        "Worn Armor",
				Description: "Old armor that still offers basic protection.",
				Type:        domain.ItemArmor,
				Rarity:      domain.RarityCommon,
				Value:       5,
				Quantity:    1,
				Defense:     2,
			},
			{
				ID:          "potion",
				Name:        "Potion",
				Description: "Restores a amount of HP.",
				Type:        domain.ItemConsumable,
				Rarity:      domain.RarityCommon,
				Value:       5,
				Quantity:    2,
				Heal:        30,
			},
		},
			EquippedWeaponID: "rusty_sword",
			EquippedArmorID:  "worn_armor",
		}

	case domain.ClassMage:
		return ClassTemplate{
			Stats: domain.Stats{
				MaxHP:   30,
				HP:      30,
				MaxMP:   28,
				MP:      28,
				Attack:  5,
				Defense: 2,
				Magic:   15,
			},

			Skills: GetStarterSkills(class),

			Items: []domain.Item{
			{
			ID:          "wooden_staff",
			Name:        "Wooden Staff",
			Description: "A apprentice staff.",
			Type:        domain.ItemWeapon,
			Rarity:      domain.RarityCommon,
			Value:       10,
			Quantity:    1,
			Attack: 	 2,
			Magic:       6,
			},
			{
				ID:          "apprentice_robe",
				Name:        "Apprentice Robe",
				Description: "A robe worn by novice spellcasters.",
				Type:        domain.ItemArmor,
				Rarity:      domain.RarityCommon,
				Value:       5,
				Quantity:    1,
				Defense:     1,
			},
			{
				ID:          "ether",
				Name:        "Ether",
				Description: "Restores a small amount of MP.",
				Type:        domain.ItemConsumable,
				Rarity:      domain.RarityCommon,
				Value:       10,
				Quantity:    2,
				MPHeal:      20,
			},
		},
			EquippedWeaponID: "wooden_staff",
			EquippedArmorID:  "apprentice_robe",
		}
	}

	return ClassTemplate{}
}