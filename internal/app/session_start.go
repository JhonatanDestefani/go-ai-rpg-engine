package app

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-inventory-management/internal/domain"
)

// HandleStart begins or resumes a session for the given telegram/session id.
func (a *App) HandleStart(ctx context.Context, telegramID int64) (*Reply, error) {
	player, err := a.PlayerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if player != nil && !player.StoryState.Completed {
		player.StoryState.InteractionState = domain.StateChoosingStart
		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}

		return &Reply{
			Text: fmt.Sprintf(
				"Welcome back, %s!\nWhat would you like to do?",
				player.Name,
			),
			Buttons: []Button{
				{ID: CallbackContinue, Label: "Continue"},
				{ID: CallbackNewGame, Label: "New Game"},
			},
		}, nil
	}

	a.setPending(telegramID, &pendingPlayer{})
	return &Reply{
		Text: "Welcome, adventurer!\nEnter your hero's name:",
	}, nil
}

// HandleInput processes one message or button press for a session.
func (a *App) HandleInput(
	ctx context.Context,
	telegramID int64,
	input Input,
) (*Reply, error) {
	callback := strings.TrimSpace(input.CallbackID)
	text := strings.TrimSpace(input.Text)

	if callback == "" && text != "" {
		if pending := a.getPending(telegramID); pending != nil && pending.Name == "" {
			return a.handleNameInput(telegramID, text)
		}
	}

	player, err := a.PlayerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if pending := a.getPending(telegramID); pending != nil && pending.Name != "" &&
		(callback == CallbackClassWarrior || callback == CallbackClassMage) {
		return a.handleClassPick(ctx, telegramID, pending, callback)
	}

	if player == nil {
		if pending := a.getPending(telegramID); pending != nil {
			if pending.Name == "" {
				return a.handleNameInput(telegramID, text)
			}
			return a.handleClassPick(ctx, telegramID, pending, callback)
		}
		return a.HandleStart(ctx, telegramID)
	}

	switch player.StoryState.InteractionState {
	case domain.StateChoosingStart:
		return a.handleStartChoice(ctx, player, telegramID, callback)
	case domain.StateChoosingName:
		if text != "" {
			return a.handleNameInput(telegramID, text)
		}
		return &Reply{Text: "Enter your hero's name:"}, nil
	case domain.StateChoosingClass:
		pending := a.getPending(telegramID)
		if pending == nil {
			pending = &pendingPlayer{Name: player.Name}
			a.setPending(telegramID, pending)
		}
		return a.handleClassPick(ctx, telegramID, pending, callback)
	case domain.StateChoosingScene:
		return a.handleSceneInput(ctx, player, telegramID, callback)
	case domain.StateChoosingCombatAction:
		return a.handleCombatAction(ctx, player, telegramID, callback)
	case domain.StateChoosingSkill:
		return a.handleSkillPick(ctx, player, telegramID, callback)
	case domain.StateChoosingItem:
		return a.handleItemPick(ctx, player, telegramID, callback)
	case domain.StateChoosingInventory:
		return a.handleInventoryMenu(ctx, player, telegramID, callback)
	case domain.StateChoosingEquip:
		return a.handleEquipPick(ctx, player, telegramID, callback)
	case domain.StateAdventureComplete:
		if callback == CallbackNewGame {
			a.clearCombat(telegramID)
			a.setPending(telegramID, &pendingPlayer{})
			return &Reply{Text: "New adventure!\nEnter your hero's name:"}, nil
		}
		return &Reply{
			Text: "Adventure complete! Send /start to play again.",
			Buttons: []Button{
				{ID: CallbackNewGame, Label: "New Game"},
			},
		}, nil
	default:
		return a.renderCurrent(ctx, player, telegramID)
	}
}

func (a *App) handleNameInput(telegramID int64, name string) (*Reply, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return &Reply{Text: "Please enter a valid name:"}, nil
	}

	a.setPending(telegramID, &pendingPlayer{Name: name})

	return &Reply{
		Text: fmt.Sprintf("Nice to meet you, %s!\nChoose your class:", name),
		Buttons: []Button{
			{ID: CallbackClassWarrior, Label: "Warrior"},
			{ID: CallbackClassMage, Label: "Mage"},
		},
	}, nil
}

func (a *App) handleClassPick(
	ctx context.Context,
	telegramID int64,
	pending *pendingPlayer,
	callback string,
) (*Reply, error) {
	var class domain.ClassType

	switch callback {
	case CallbackClassWarrior:
		class = domain.ClassWarrior
	case CallbackClassMage:
		class = domain.ClassMage
	default:
		return &Reply{
			Text: "Choose your class:",
			Buttons: []Button{
				{ID: CallbackClassWarrior, Label: "Warrior"},
				{ID: CallbackClassMage, Label: "Mage"},
			},
		}, nil
	}

	if pending == nil || pending.Name == "" {
		return &Reply{Text: "Enter your hero's name first:"}, nil
	}

	player := a.Game.CreatePlayer(pending.Name, class)
	player.TelegramID = telegramID
	player.StoryState.InteractionState = domain.StateChoosingScene
	player.StoryState.SceneNumber = 1

	existing, err := a.PlayerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		player.ID = existing.ID
		if err := a.PlayerRepo.Update(ctx, &player); err != nil {
			return nil, err
		}
	} else {
		if err := a.PlayerRepo.Create(ctx, &player); err != nil {
			return nil, err
		}
	}

	a.clearPending(telegramID)
	a.clearCombat(telegramID)

	return a.generateAndShowScene(ctx, &player, telegramID)
}

func (a *App) handleStartChoice(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
	callback string,
) (*Reply, error) {
	switch callback {
	case CallbackContinue:
		return a.renderCurrent(ctx, player, telegramID)
	case CallbackNewGame:
		a.clearCombat(telegramID)
		a.setPending(telegramID, &pendingPlayer{})
		player.StoryState.InteractionState = domain.StateChoosingName
		_ = a.Game.SaveGame(ctx, player)
		return &Reply{Text: "New adventure!\nEnter your hero's name:"}, nil
	default:
		return &Reply{
			Text: "Choose an option:",
			Buttons: []Button{
				{ID: CallbackContinue, Label: "Continue"},
				{ID: CallbackNewGame, Label: "New Game"},
			},
		}, nil
	}
}

func parsePrefixedID(callback, prefix string) (string, bool) {
	if !strings.HasPrefix(callback, prefix) {
		return "", false
	}
	return strings.TrimPrefix(callback, prefix), true
}

func parsePrefixedInt(callback, prefix string) (int, bool) {
	raw, ok := parsePrefixedID(callback, prefix)
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return n, true
}
