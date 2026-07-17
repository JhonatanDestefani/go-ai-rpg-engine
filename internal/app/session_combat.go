package app

import (
	"context"
	"fmt"

	"go-inventory-management/internal/domain"
)

func (a *App) combatReply(combat *domain.CombatState, prefix string) *Reply {
	text := fmt.Sprintf(
		"%s\n%s ❤️ HP: %d/%d |🔷 MP: %d/%d\n%s ❤️ HP: %d/%d\n\nChoose your action:",
		prefix,
		combat.Player.Name,
		combat.Player.Stats.HP,
		combat.Player.Stats.MaxHP,
		combat.Player.Stats.MP,
		combat.Player.Stats.MaxMP,
		combat.Monster.Name,
		combat.Monster.Stats.HP,
		combat.Monster.Stats.MaxHP,
	)

	return &Reply{
		Text: stringsTrim(text),
		Buttons: []Button{
			{ID: CallbackAttack, Label: "Basic Attack"},
			{ID: CallbackSkillMenu, Label: "Skill"},
			{ID: CallbackItemMenu, Label: "Item"},
		},
	}
}

func stringsTrim(s string) string {
	for len(s) > 0 && s[0] == '\n' {
		s = s[1:]
	}
	return s
}

func (a *App) handleCombatAction(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
	callback string,
) (*Reply, error) {
	combat := a.getCombat(telegramID)
	if combat == nil {
		return a.renderCurrent(ctx, player, telegramID)
	}
	combat.Player = player

	switch callback {
	case CallbackAttack:
		return a.resolvePlayerTurn(ctx, player, telegramID, combat, domain.ActionBasicAttack, "")
	case CallbackSkillMenu:
		player.StoryState.InteractionState = domain.StateChoosingSkill
		_ = a.Game.SaveGame(ctx, player)
		return a.skillMenuReply(player), nil
	case CallbackItemMenu:
		player.StoryState.InteractionState = domain.StateChoosingItem
		_ = a.Game.SaveGame(ctx, player)
		return a.itemMenuReply(player), nil
	default:
		return a.combatReply(combat, ""), nil
	}
}

func (a *App) skillMenuReply(player *domain.Player) *Reply {
	if len(player.Skills) == 0 {
		return &Reply{
			Text: "You have no skills.",
			Buttons: []Button{
				{ID: CallbackSkillBack, Label: "Back"},
			},
		}
	}

	buttons := make([]Button, 0, len(player.Skills)+1)
	text := "Choose a skill:\n"
	for _, skill := range player.Skills {
		text += fmt.Sprintf(
			"\n- %s (🔷 MP %d,Power %d)",
			skill.Name,
			skill.MPCost,
			skill.Power,
		)
		buttons = append(buttons, Button{
			ID:    PrefixSkill + skill.ID,
			Label: skill.Name,
		})
	}
	buttons = append(buttons, Button{ID: CallbackSkillBack, Label: "Back"})
	return &Reply{Text: text, Buttons: buttons}
}

func (a *App) itemMenuReply(player *domain.Player) *Reply {
	items := player.Inventory.GetUsableItems()
	if len(items) == 0 {
		return &Reply{
			Text: "No usable items.",
			Buttons: []Button{
				{ID: CallbackItemBack, Label: "Back"},
			},
		}
	}

	buttons := make([]Button, 0, len(items)+1)
	text := "Choose an item:"
	for _, item := range items {
		text += fmt.Sprintf("\n- %s x%d", item.Name, item.Quantity)
		buttons = append(buttons, Button{
			ID:    PrefixItem + item.ID,
			Label: item.Name,
		})
	}
	buttons = append(buttons, Button{ID: CallbackItemBack, Label: "Back"})
	return &Reply{Text: text, Buttons: buttons}
}

func (a *App) handleSkillPick(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
	callback string,
) (*Reply, error) {
	combat := a.getCombat(telegramID)
	if combat == nil {
		return a.renderCurrent(ctx, player, telegramID)
	}
	combat.Player = player

	if callback == CallbackSkillBack {
		player.StoryState.InteractionState = domain.StateChoosingCombatAction
		_ = a.Game.SaveGame(ctx, player)
		return a.combatReply(combat, ""), nil
	}

	skillID, ok := parsePrefixedID(callback, PrefixSkill)
	if !ok {
		return a.skillMenuReply(player), nil
	}

	for _, skill := range player.Skills {
		if skill.ID == skillID && player.Stats.MP < skill.MPCost {
			reply := a.skillMenuReply(player)
			reply.Text = "Not enough 🔷 MP.\n\n" + reply.Text
			return reply, nil
		}
	}

	return a.resolvePlayerTurn(ctx, player, telegramID, combat, domain.ActionSkill, skillID)
}

func (a *App) handleItemPick(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
	callback string,
) (*Reply, error) {
	combat := a.getCombat(telegramID)

	if callback == CallbackItemBack || callback == CallbackInvBack {
		if combat != nil && !combat.IsFinished {
			player.StoryState.InteractionState = domain.StateChoosingCombatAction
			_ = a.Game.SaveGame(ctx, player)
			return a.combatReply(combat, ""), nil
		}
		player.StoryState.InteractionState = domain.StateChoosingInventory
		_ = a.Game.SaveGame(ctx, player)
		return a.inventoryReply(player), nil
	}

	itemID, ok := parsePrefixedID(callback, PrefixItem)
	if !ok {
		if combat != nil && !combat.IsFinished {
			return a.itemMenuReply(player), nil
		}
		return a.useConsumableFromInventory(ctx, player, telegramID)
	}

	if combat == nil || combat.IsFinished {
		return a.useInventoryItemDirect(ctx, player, itemID)
	}

	combat.Player = player
	return a.resolvePlayerTurn(ctx, player, telegramID, combat, domain.ActionItem, itemID)
}

func (a *App) resolvePlayerTurn(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
	combat *domain.CombatState,
	action domain.CombatActionType,
	actionID string,
) (*Reply, error) {
	message := a.Game.ExecutePlayerTurn(combat, action, actionID)

	if a.Game.IsCombatOver(combat) {
		return a.finishCombat(ctx, player, telegramID, combat, message)
	}

	monsterMsg := a.Game.ExecuteMonsterTurn(combat)
	message += "\n" + monsterMsg

	if a.Game.IsCombatOver(combat) {
		return a.finishCombat(ctx, player, telegramID, combat, message)
	}

	player.StoryState.InteractionState = domain.StateChoosingCombatAction
	if err := a.Game.SaveGame(ctx, player); err != nil {
		return nil, err
	}
	a.setCombat(telegramID, combat)

	return a.combatReply(combat, message), nil
}

func (a *App) finishCombat(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
	combat *domain.CombatState,
	message string,
) (*Reply, error) {
	a.clearCombat(telegramID)

	if !a.Game.PlayerWon(combat) {
		player.StoryState.PreviousOutcome = "The player was defeated."
		player.StoryState.InteractionState = domain.StateAdventureComplete
		player.StoryState.Completed = true
		_ = a.Game.SaveGame(ctx, player)
		return &Reply{
			Text: message + "\n\nGame over. Send /start to try again.",
			Buttons: []Button{
				{ID: CallbackNewGame, Label: "New Game"},
			},
		}, nil
	}

	player.StoryState.PreviousOutcome = fmt.Sprintf(
		"The player defeated %s. The monster is dead and the encounter is finished.",
		combat.Monster.Name,
	)

	sceneNumber := player.StoryState.SceneNumber
	if sceneNumber == 15 {
		player.StoryState.Completed = true
		player.StoryState.InteractionState = domain.StateAdventureComplete
		player.StoryState.CurrentScene = domain.GeneratedScene{}
		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}
		return &Reply{
			Text: message + "\n\nYou defeated the 🐉 Dragon and completed the adventure!",
			Buttons: []Button{
				{ID: CallbackNewGame, Label: "New Game"},
			},
		}, nil
	}

	// After combat that started with the scene, show post-combat choices if any,
	// otherwise advance to the next scene.
	scene := player.StoryState.CurrentScene
	if len(scene.Choices) > 0 {
		player.StoryState.InteractionState = domain.StateChoosingScene
		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}
		reply := a.sceneReply(player)
		reply.Text = message + "\n\n" + reply.Text
		return reply, nil
	}

	player.StoryState.PreviousScene = scene.Description
	player.StoryState.CurrentScene = domain.GeneratedScene{}
	player.StoryState.SceneNumber++
	if err := a.Game.SaveGame(ctx, player); err != nil {
		return nil, err
	}

	reply, err := a.generateAndShowScene(ctx, player, telegramID)
	if err != nil {
		return nil, err
	}
	reply.Text = message + "\n\n" + reply.Text
	return reply, nil
}
