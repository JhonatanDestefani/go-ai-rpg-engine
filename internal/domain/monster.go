package domain

import "errors"

type Monster struct {
	ID        string
	Name      string
	Level     int
	Stats     Stats
	Skills    []Skill
	DropTable DropTable

	XPReward   int
	GoldReward int

	CanUseItems bool
	Potions     int
}

func NewGoblin() Monster {
	return Monster{
		ID:    "goblin",
		Name:  "Goblin",
		Level: 1,
		Stats: Stats{
			HP:      30,
			MaxHP:   30,
			MP:      10,
			MaxMP:   10,
			Attack:  8,
			Defense: 3,
			Magic:   2,
		},
		Skills: []Skill{
			GoblinSlash(),
		},
		DropTable: CreateDropTable(
			"goblin",
			NewDrop(Potion, 5, 9, 1, 2),
			NewDrop(BronzeSword, 10, 14, 1, 1),
			NewDrop(LeatherArmor, 15, 19, 1, 1),
			NewDrop(GoblinSword, 20, 20, 1, 1)),
		XPReward:    20,
		GoldReward:  8,
		CanUseItems: true,
		Potions:     1,
	}
}

func NewOrc() Monster {
	return Monster{
		ID:    "orc",
		Name:  "Orc",
		Level: 2,
		Stats: Stats{
			HP:      55,
			MaxHP:   55,
			MP:      12,
			MaxMP:   12,
			Attack:  12,
			Defense: 3,
			Magic:   1,
		},
		Skills: []Skill{
			HeavySmash(),
		},
		DropTable: CreateDropTable(
			"orc",
			NewDrop(ShamanStaff, 5, 9, 1, 2),
			NewDrop(LeatherArmor, 10, 14, 1, 1),
			NewDrop(IronSword, 15, 19, 1, 1),
			NewDrop(OrcAxe, 20, 20, 1, 1)),
		XPReward:    50,
		GoldReward:  15,
		CanUseItems: false,
		Potions:     0,
	}
}

func NewGuardian() Monster {
	return Monster{
		ID:    "guardian",
		Name:  "Guardian",
		Level: 5,
		Stats: Stats{
			HP:      100,
			MaxHP:   100,
			MP:      40,
			MaxMP:   40,
			Attack:  12,
			Defense: 5,
			Magic:   10,
		},
		Skills: []Skill{
			Earthquake(),
			LightningWave(),
		},
		DropTable: CreateDropTable(
			"guardian",
			NewDrop(GuardianArmor, 10, 15, 1, 1),
			NewDrop(GuardianSword,16,20, 1, 1),
		),
		XPReward:    150,
		GoldReward:  50,
		CanUseItems: false,
		Potions:     0,
	}
}

func NewDragonStatue() Monster {
	return Monster{
		ID:    "dragon_statue",
		Name:  "Dragon Statue",
		Level: 6,
		Stats: Stats{
			HP:      150,
			MaxHP:   150,
			MP:      40,
			MaxMP:   40,
			Attack:  20,
			Defense: 5,
			Magic:   18,
		},
		Skills: []Skill{
			StoneSlam(),
		},
		DropTable: CreateDropTable(
			"guardian",
			NewDrop(DragonStaff, 10, 15, 1, 1),
			NewDrop(DragonRobe,16,20, 1, 1),
		),
		XPReward:    200,
		GoldReward:  75,
		CanUseItems: false,
		Potions:     0,
	}
}


func NewDragon()Monster {
	return Monster{
		ID:    "dragon",
		Name:  "Dragon",
		Level: 8,
		Stats: Stats{
			HP:      300,
			MaxHP:   300,
			MP:      60,
			MaxMP:   60,
			Attack:  25,
			Defense: 8,
			Magic:   18,
		},
		Skills: []Skill{
			DragonBreath(),
			
		},
		XPReward:    500,
		GoldReward:  200,
		CanUseItems: false,
		Potions:     0,
	}
}

func CreateMonster(monsterID string) (Monster, error) {
	switch monsterID {
	case "goblin":
		return NewGoblin(), nil
	case "orc":
		return NewOrc(), nil
	case "guardian":
		return NewGuardian(), nil
	case "dragon_statue":
		return NewDragonStatue(), nil	
	case "dragon":
		return NewDragon(), nil
	default:
		return Monster{}, errors.New("monster not found")
	}
}
