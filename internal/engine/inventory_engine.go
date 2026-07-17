package engine

import (
	"context"
	"errors"
	"fmt"

	"go-inventory-management/internal/domain"
)

func (g *GameEngine) EquipItem(
	ctx context.Context,
	player *domain.Player,
	itemID string,
) error {
	item, found := player.Inventory.FindItem(itemID)
	if !found {
		return errors.New("item not found in inventory")
	}

	switch item.Type {
	case domain.ItemWeapon:
		player.EquippedWeaponID = item.ID

	case domain.ItemArmor:
		player.EquippedArmorID = item.ID

	default:
		return fmt.Errorf(
			"%s cannot be equipped",
			item.Name,
		)
	}

	if err := g.playerRepository.Update(ctx,player); err != nil {
		return fmt.Errorf(
			"failed to save equipped item: %w",
			err,
		)
	}

	return nil
}
func (g *GameEngine) GetEquippedAttackBonus(
	player *domain.Player,
) int {
	if player.EquippedWeaponID == "" {
		return 0
	}

	item, found := player.Inventory.FindItem(
		player.EquippedWeaponID,
	)
	if !found {
		return 0
	}

	return item.Attack
}

func (g *GameEngine) GetEquippedDefenseBonus(
	player *domain.Player,
) int {
	if player.EquippedArmorID == "" {
		return 0
	}

	item, found := player.Inventory.FindItem(
		player.EquippedArmorID,
	)
	if !found {
		return 0
	}

	return item.Defense
}