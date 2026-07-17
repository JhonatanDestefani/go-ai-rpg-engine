package domain

func GetClassSkills(class ClassType) []Skill {
	switch class {
	case ClassWarrior:
		return []Skill{

			Slash(),
			GuardBreak(),
			BattleCry(),
			EarthShatter(),
		}
	case ClassMage:
		return []Skill{

			Fireball(),
			Heal(),
			IceLance(),
			LightningBolt(),
			Meteor(),
		}
	}
	return nil
}

func Slash() Skill {
	return Skill{
		ID:            "slash",
		Name:          "Slash",
		Description:   "A Slash that deals physical damage and make them bleed.",
		MPCost:        2,
		Power:         6,
		DamageType:    DamagePhysical,
		SkillType: 	   SkillAttack,	
		RequiredLevel: 1,
		Effects: []CombatEffect{
			{
				Type:     EffectBleeding,
				Duration: 3,
				Power:    5,
				IsDebuff: true,
				IsBuff:   false,
			},
		},
	}
}

func GuardBreak() Skill {
	return Skill{
		ID:            "guard-break",
		Name:          "Guard Break",
		Description:   "Lowers the enemy's defense and deals physical damage.",
		MPCost:        5,
		Power:         8,
		DamageType:    DamagePhysical,
		SkillType: 	   SkillAttack,	
		RequiredLevel: 3,
	}
}

func BattleCry() Skill {
	return Skill{
		ID:            "battle-cry",
		Name:          "Battle Cry",
		Description:   "Give a very high shout damaging all its surroundings.",
		MPCost:        4,
		Power:         5,
		SkillType: 	   SkillAttack,	
		RequiredLevel: 5,
	}
}

func EarthShatter() Skill {
	return Skill{
		ID:            "earth-shatter",
		Name:          "Earth Shatter",
		Description:   "Power surge attack that destroy the earth beneath the enemy, dealing physical damage and stunning them.",
		MPCost:        14,
		Power:         30,
		SkillType: 	   SkillAttack,	
		DamageType:    DamagePhysical,
		RequiredLevel: 8,
	}
}

func Heal() Skill {
	return Skill{
		ID:            "heal",
		Name:          "Heal",
		Description:   "A healing spell that recovers some health.",
		SkillType:     SkillSupport,
		MPCost:        6,
		Power:         20,
		RequiredLevel: 2,
	}
}

func Fireball() Skill {
	return Skill{
		ID:            "fireball",
		Name:          "Fireball",
		Description:   "A powerful fire-based spell that deals magical damage.",
		MPCost:        6,
		Power:         10,
		DamageType:    DamageMagical,
		SkillType: 	   SkillAttack,	
		RequiredLevel: 1,
		Effects: []CombatEffect{
			{
				Type:     EffectBurning,
				Duration: 5,
				Power:    3,
				IsDebuff: true,
				IsBuff:   false,
			},
		},
	}
}

func IceLance() Skill {
	return Skill{
		ID:            "ice-lance",
		Name:          "Ice Lance",
		Description:   "A sharp lance of ice that deals magical damage and freezes the target.",
		MPCost:        6,
		Power:         15,
		DamageType:    DamageMagical,
		SkillType: 	   SkillAttack,	
		RequiredLevel: 2,
	}
}

func LightningBolt() Skill {
	return Skill{
		ID:            "lightning-bolt",
		Name:          "Lightning Bolt",
		Description:   "A powerful lightning-based spell that deals magical damage.",
		MPCost:        8,
		Power:         20,
		DamageType:    DamageMagical,
		SkillType: 	   SkillAttack,	
		RequiredLevel: 5,
	}
}

func Meteor() Skill {
	return Skill{
		ID:            "meteor",
		Name:          "Meteor",
		Description:   "A devastating meteor strike that deals massive magical damage.",
		MPCost:        20,
		Power:         40,
		DamageType:    DamageMagical,
		SkillType: 	   SkillAttack,	
		RequiredLevel: 8,
	}
}
