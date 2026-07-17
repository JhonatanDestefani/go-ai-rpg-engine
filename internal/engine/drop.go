package engine

import (
	"go-inventory-management/internal/domain"
)

func RollDrops(dropTable domain.DropTable) ([]domain.Item, int) {
	drops := []domain.Item{}

	roll := RollD20().Total

	for _, drop := range dropTable.Drops {
		if roll < drop.DiceMin || roll > drop.DiceMax {
			continue
		}

		item := drop.Item

		if drop.QuantityMax > drop.QuantityMin {
			quantityRoll := RollDice(1, drop.QuantityMax-drop.QuantityMin+1).Total
			item.Quantity = drop.QuantityMin + quantityRoll - 1
		} else {
			item.Quantity = drop.QuantityMin
		}

		drops = append(drops, item)
	}

	return drops, roll
}
