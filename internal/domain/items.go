package domain

var Potion = Item{
	ID:          "potion",
	Name:        "Potion",
	Description: "Restores 30 HP.",
	Type:        ItemConsumable,
	Rarity:      RarityCommon,
	Value:       5,
	Quantity:    1,
	Heal:        30,
}

var BronzeSword = Item{
	ID:          "bronze_sword",
	Name:        "Bronze Sword",
	Description: "A simple bronze sword.",
	Type:        ItemWeapon,
	Rarity:      RarityCommon,
	Value:       15,
	Quantity:    1,
	Attack:      5,
}

var LeatherArmor = Item{
	ID:          "leather_armor",
	Name:        "Leather Armor",
	Description: "Basic leather protection.",
	Type:        ItemArmor,
	Rarity:      RarityCommon,
	Value:       20,
	Quantity:    1,
	Defense:     2,
}

var GoblinSword = Item{
	ID:          "goblin_sword",
	Name:        "Goblin Sword",
	Description: "A rare sword taken from a goblin.",
	Type:        ItemWeapon,
	Rarity:      RarityRare,
	Value:       50,
	Quantity:    1,
	Attack:      15,
}

var IronSword = Item{
	ID:          "iron_sword",
	Name:        "Iron Sword",
	Description: "A sturdy sword forged from iron.",
	Type:        ItemWeapon,
	Rarity:      RarityCommon,
	Value:       35,
	Quantity:    1,
	Attack:      8,
}

var OrcAxe = Item{
	ID:          "orc_axe",
	Name:        "Orc Axe",
	Description: "A brutal axe wielded by Orc warriors.",
	Type:        ItemWeapon,
	Rarity:      RarityRare,
	Value:       90,
	Quantity:    1,
	Attack:      22,
}

var ShamanStaff = Item{
	ID:          "shamanstaff",
	Name:        "ShamanStaff",
	Description: "An old Shaman staff but still powerful",
	Type:        ItemWeapon,
	Rarity:      RarityUnique,
	Value:       0,
	Quantity:    1,
	Attack:      5,
	Magic:       15,
}

var GuardianArmor = Item{
	ID:          "guardian_armor",
	Name:        "GuardianArmor",
	Description: "Guardian armor made of ancient steel",
	Type:        ItemArmor,
	Rarity:      RarityRare,
	Value:       0,
	Quantity:    1,
	Defense:     15,
}

var GuardianSword = Item{
	ID:          "guardian_sword",
	Name:        "GuardianSword",
	Description: "Guardian sword made of ancient steel",
	Type:        ItemWeapon,
	Rarity:      RarityRare,
	Value:       0,
	Quantity:    1,
	Attack:      30,
}

var DragonStaff = Item{
	ID:          "dragon_staff",
	Name:        "DragonStaff",
	Description: "DragonStaff forged from dragon bones",
	Type:        ItemArmor,
	Rarity:      RarityRare,
	Value:       0,
	Quantity:    1,
	Attack:     15,
	Magic:      30,
}

var DragonRobe = Item{
	ID:          "dragon_robe",
	Name:        "DragonRobe",
	Description: "Dragon Robe made of dragon scales",
	Type:        ItemArmor,
	Rarity:      RarityRare,
	Value:       0,
	Quantity:    1,
	Defense:     15,
}
