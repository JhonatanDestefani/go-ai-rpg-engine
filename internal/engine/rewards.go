package engine

import "go-inventory-management/internal/domain"

func (g *GameEngine) AddGold(player *domain.Player, amount int) {
	if player == nil || amount <= 0 {
		return
	}

	player.Gold += amount
}

func (g *GameEngine) GiveCombatRewards(
	player *domain.Player,
	monster *domain.Monster,
) ([]domain.Item, int, bool, []domain.Skill) {
	if player == nil || monster == nil {
		return nil, 0, false, nil
	}

	leveledUp, newSkills := g.AddXP(
		player,
		monster.XPReward,
	)

	g.AddGold(
		player,
		monster.GoldReward,
	)

	drops, dropRoll := RollDrops(
		monster.DropTable,
	)

	for _, item := range drops {
		player.Inventory.AddItem(item)
	}

	return drops, dropRoll, leveledUp, newSkills
}
