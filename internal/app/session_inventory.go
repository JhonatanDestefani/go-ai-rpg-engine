package app

import (
	"context"
	"fmt"

	"go-inventory-management/internal/domain"
	"go-inventory-management/internal/engine"
)

func (a *App) inventoryReply(player *domain.Player) *Reply {
	text := "INVENTORY\n\n" + formatInventoryText(player)
	return &Reply{
		Text: text,
		Buttons: []Button{
			{ID: CallbackInvEquip, Label: "Equip item"},
			{ID: CallbackInvUse, Label: "Use item"},
			{ID: CallbackInvBack, Label: "Back to adventure"},
		},
	}
}

func formatInventoryText(player *domain.Player) string {
	if len(player.Inventory.Items) == 0 {
		return "Your inventory is empty."
	}

	text := ""
	for _, item := range player.Inventory.Items {
		equipped := ""
		switch item.Type {
		case domain.ItemWeapon:
			if item.ID == player.EquippedWeaponID {
				equipped = " [Equipped]"
			}
		case domain.ItemArmor:
			if item.ID == player.EquippedArmorID {
				equipped = " [Equipped]"
			}
		}

		text += fmt.Sprintf("- %s x%d%s\n", item.Name, item.Quantity, equipped)

		switch item.Type {
		case domain.ItemWeapon:
			text += fmt.Sprintf("⚔️ Attack: +%d\n", item.Attack)
			text += fmt.Sprintf("🪄 Magic: +%d\n", item.Magic)
		case domain.ItemArmor:
			text += fmt.Sprintf("🛡️  Defense: +%d\n", item.Defense)
		case domain.ItemConsumable:
			text += fmt.Sprintf("🧪  Heal: %d |💎 MP Heal: %d\n", item.Heal, item.MPHeal)
		}
	}
	return text
}

func (a *App) handleInventoryMenu(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
	callback string,
) (*Reply, error) {
	switch callback {
	case CallbackInvBack:
		player.StoryState.InteractionState = domain.StateChoosingScene
		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}
		return a.sceneReply(player), nil

	case CallbackInvEquip:
		player.StoryState.InteractionState = domain.StateChoosingEquip
		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}
		return a.equipMenuReply(player), nil

	case CallbackInvUse:
		return a.useConsumableFromInventory(ctx, player, telegramID)

	default:
		return a.inventoryReply(player), nil
	}
}

func (a *App) equipMenuReply(player *domain.Player) *Reply {
	items := player.Inventory.GetEquippableItems()
	if len(items) == 0 {
		return &Reply{
			Text: "No equippable items.",
			Buttons: []Button{
				{ID: CallbackEquipBack, Label: "Back"},
			},
		}
	}

	buttons := make([]Button, 0, len(items)+1)
	text := "Choose an item to equip:"
	for _, item := range items {
		label := item.Name
		if item.Type == domain.ItemWeapon && item.ID == player.EquippedWeaponID {
			label += " [Equipped]"
		}
		if item.Type == domain.ItemArmor && item.ID == player.EquippedArmorID {
			label += " [Equipped]"
		}
		buttons = append(buttons, Button{
			ID:    PrefixEquip + item.ID,
			Label: label,
		})
	}
	buttons = append(buttons, Button{ID: CallbackEquipBack, Label: "Back"})
	return &Reply{Text: text, Buttons: buttons}
}

func (a *App) handleEquipPick(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
	callback string,
) (*Reply, error) {
	if callback == CallbackEquipBack {
		player.StoryState.InteractionState = domain.StateChoosingInventory
		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}
		return a.inventoryReply(player), nil
	}

	itemID, ok := parsePrefixedID(callback, PrefixEquip)
	if !ok {
		return a.equipMenuReply(player), nil
	}

	if err := a.Game.EquipItem(ctx, player, itemID); err != nil {
		reply := a.equipMenuReply(player)
		reply.Text = "Failed: " + err.Error() + "\n\n" + reply.Text
		return reply, nil
	}

	player.StoryState.InteractionState = domain.StateChoosingInventory
	_ = a.Game.SaveGame(ctx, player)
	reply := a.inventoryReply(player)
	reply.Text = "Equipped!\n\n" + reply.Text
	return reply, nil
}

func (a *App) useConsumableFromInventory(
	ctx context.Context,
	player *domain.Player,
	_ int64,
) (*Reply, error) {
	items := player.Inventory.GetUsableItems()
	if len(items) == 0 {
		reply := a.inventoryReply(player)
		reply.Text = "No consumable items.\n\n" + reply.Text
		return reply, nil
	}

	// Use first potion-like item for simplicity from inventory menu;
	// detailed pick uses choosing_item-style buttons.
	player.StoryState.InteractionState = domain.StateChoosingItem
	if err := a.Game.SaveGame(ctx, player); err != nil {
		return nil, err
	}

	buttons := make([]Button, 0, len(items)+1)
	text := "Choose an item to use:"
	for _, item := range items {
		text += fmt.Sprintf("\n- %s x%d", item.Name, item.Quantity)
		buttons = append(buttons, Button{
			ID:    PrefixItem + item.ID,
			Label: item.Name,
		})
	}
	buttons = append(buttons, Button{ID: CallbackInvBack, Label: "Back"})

	// Reuse item pick but route back to inventory when not in combat.
	return &Reply{Text: text, Buttons: buttons}, nil
}

// Override item pick when not in combat: use from inventory.
func (a *App) useInventoryItemDirect(
	ctx context.Context,
	player *domain.Player,
	itemID string,
) (*Reply, error) {
	message := engine.UseConsumableItem(
		&player.Inventory,
		&player.Stats,
		itemID,
	)
	player.StoryState.InteractionState = domain.StateChoosingInventory
	if err := a.Game.SaveGame(ctx, player); err != nil {
		return nil, err
	}
	reply := a.inventoryReply(player)
	reply.Text = message + "\n\n" + reply.Text
	return reply, nil
}
