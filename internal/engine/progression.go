package engine

import (
	"go-inventory-management/internal/domain"
)

var XPRequirements = map[int]int{
	1: 15,
	2: 25,
	3: 40,
	4: 70,
	5: 100,
	6: 150,
	7: 220,
	8: 400,
}

type LevelUpBonus struct {
	MaxHP   int
	MaxMP   int
	Attack  int
	Defense int
	Magic   int
}

var ClassLevelBonus = map[domain.ClassType]LevelUpBonus{
	domain.ClassWarrior: {
		MaxHP:   15,
		MaxMP:   4,
		Attack:  4,
		Defense: 2,
		Magic:   3,
	},
	domain.ClassMage: {
		MaxHP:   10,
		MaxMP:   8,
		Attack:  2,
		Defense: 2,
		Magic:   4,
	},
}

func (g *GameEngine) AddXP(
	player *domain.Player,
	amount int,
) (bool, []domain.Skill) {
	if player == nil || amount <= 0 {
		return false, nil
	}

	player.XP += amount

	if player.Level <= 0 {
		player.Level = 1
	}

	leveledUp := false
	var unlockedSkills []domain.Skill

	for {
		needed := XPToNextLevel(player.Level)

		if player.XP < needed {
			break
		}

		player.XP -= needed

		newSkills := g.LevelUp(player)

		unlockedSkills = append(
			unlockedSkills,
			newSkills...,
		)

		leveledUp = true
	}

	return leveledUp, unlockedSkills
}

func (g *GameEngine) LevelUp(
	player *domain.Player,
) []domain.Skill {
	player.Level++

	bonus, found := ClassLevelBonus[player.Class]
	if found {
		player.Stats.MaxHP += bonus.MaxHP
		player.Stats.MaxMP += bonus.MaxMP
		player.Stats.Attack += bonus.Attack
		player.Stats.Defense += bonus.Defense
		player.Stats.Magic += bonus.Magic
	}

	player.Stats.HP = player.Stats.MaxHP
	player.Stats.MP = player.Stats.MaxMP

	return g.UnlockNewSkills(player)
}

func XPToNextLevel(level int) int {
	if level <= 0 {
		return XPRequirements[1]
	}

	if requiredXP, found := XPRequirements[level]; found {
		return requiredXP
	}

	return level * 10000
}

func GetUnlockedSkills(class domain.ClassType, level int) []domain.Skill {
	allSkills := domain.GetClassSkills(class)

	unlocked := []domain.Skill{}

	for _, skill := range allSkills {
		if level >= skill.RequiredLevel {
			unlocked = append(unlocked, skill)
		}
	}

	return unlocked
}

func (g *GameEngine) UnlockNewSkills(player *domain.Player) []domain.Skill {
	allSkills := domain.GetClassSkills(player.Class)

	newSkills := []domain.Skill{}

	for _, skill := range allSkills {
		if player.Level < skill.RequiredLevel {
			continue
		}

		if _, found := FindSkillByID(player.Skills, skill.ID); found {
			continue
		}

		player.Skills = append(player.Skills, skill)

		newSkills = append(newSkills, skill)
	}
	return newSkills
}
