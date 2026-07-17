package engine

import (
	"fmt"

	"go-inventory-management/internal/domain"
)

func UseConsumableItem(
	userInventory *domain.Inventory,
	targetStats *domain.Stats,
	itemID string,
) string {
	item, found := userInventory.FindItem(itemID)
	if !found {
		return "Item not found."
	}

	if item.Type != domain.ItemConsumable {
		return "This item cannot be used."
	}

	hpBefore := targetStats.HP
	mpBefore := targetStats.MP

	if item.Heal > 0 {
		targetStats.HP += item.Heal
		if targetStats.HP > targetStats.MaxHP {
			targetStats.HP = targetStats.MaxHP
		}
	}

	if item.MPHeal > 0 {
		targetStats.MP += item.MPHeal
		if targetStats.MP > targetStats.MaxMP {
			targetStats.MP = targetStats.MaxMP
		}
	}

	userInventory.RemoveItem(itemID, 1)

	hpRestored := targetStats.HP - hpBefore
	mpRestored := targetStats.MP - mpBefore

	message := fmt.Sprintf("Used %s.", item.Name)
	if hpRestored > 0 {
		message += fmt.Sprintf(" Restored %d HP.", hpRestored)
	}
	if mpRestored > 0 {
		message += fmt.Sprintf(" Restored %d MP.", mpRestored)
	}
	if hpRestored == 0 && mpRestored == 0 && (item.Heal > 0 || item.MPHeal > 0) {
		message += " (already at full)"
	}

	return message
}
