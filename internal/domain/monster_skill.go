package domain

func GoblinSlash() Skill {
	return Skill{
		ID:          "goblin_quick_slash",
		Name:        "Quick Slash",
		Description: "A quick slash with a rusty dagger.",
		DamageType:  DamagePhysical,
		MPCost:      2,
		Power:       4,
	}
}

func HeavySmash() Skill {
	return Skill{
		ID:          "orc_heavy_smash",
		Name:        "Heavy Smash",
		Description: "A powerful smash with a Axe.",
		DamageType:  DamagePhysical,
		MPCost:      4,
		Power:       6,
	}
}

func Earthquake() Skill {
	return Skill{
		ID:          "earthquake",
		Name:        "Earthquake",
		Description: "Shatter the earth with a strong blow",
		DamageType:  DamagePhysical,
		MPCost:      4,
		Power:       8,
	}
}

func LightningWave() Skill {
	return Skill{
		ID:          "lightningwave",
		Name:        "LightningWave",
		Description: "Create a deadly wave of electricity",
		DamageType:  DamagePhysical,
		MPCost:      8,
		Power:       12,
	}
}

func StoneSlam() Skill {
	return Skill{
		ID:          "stoneslam",
		Name:        "Stone Slam",
		Description: "Slam the ground with its huge sword",
		DamageType:  DamagePhysical,
		MPCost:      10,
		Power:       20,
	}
}

func DragonBreath() Skill {
	return Skill{
		ID:          "dragonbreath",
		Name:        "DragonBreath",
		Description: "Dragon most powerfull attack",
		DamageType:  DamagePhysical,
		MPCost:      20,
		Power:       40,
	}
}

