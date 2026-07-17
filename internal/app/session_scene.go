package app

import (
	"context"
	"fmt"
	"strings"
	"math/rand"
	"go-inventory-management/internal/domain"
	"go-inventory-management/internal/engine"
)

func (a *App) renderCurrent(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
) (*Reply, error) {
	if player.StoryState.Completed {
		player.StoryState.InteractionState = domain.StateAdventureComplete
		_ = a.Game.SaveGame(ctx, player)
		return &Reply{
			Text: "Adventure complete! Send /start for a new game.",
			Buttons: []Button{
				{ID: CallbackNewGame, Label: "New Game"},
			},
		}, nil
	}

	combat := a.getCombat(telegramID)
	if combat != nil && combat.IsActive && !combat.IsFinished {
		player.StoryState.InteractionState = domain.StateChoosingCombatAction
		_ = a.Game.SaveGame(ctx, player)
		return a.combatReply(combat, ""), nil
	}

	if player.StoryState.CurrentScene.Title != "" {
		player.StoryState.InteractionState = domain.StateChoosingScene
		_ = a.Game.SaveGame(ctx, player)
		return a.sceneReply(player), nil
	}

	return a.generateAndShowScene(ctx, player, telegramID)
}

func (a *App) generateAndShowScene(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
) (*Reply, error) {
	sceneNumber := player.StoryState.SceneNumber
	if sceneNumber <= 0 {
		sceneNumber = 1
		player.StoryState.SceneNumber = 1
	}

	if player.StoryState.CurrentScene.Title != "" &&
		player.StoryState.SceneNumber == sceneNumber {
		player.StoryState.InteractionState = domain.StateChoosingScene
		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}
		return a.sceneReply(player), nil
	}

	scene, err := a.generateValidScene(ctx, player, sceneNumber)
	if err != nil {
		return &Reply{Text: "Failed to generate scene: " + err.Error()}, nil
	}

	player.StoryState.CurrentScene = scene
	player.StoryState.CurrentRegion = engine.GetRegionForScene(sceneNumber)
	player.StoryState.InteractionState = domain.StateChoosingScene

	if scene.MonsterID != "" {
		monster, err := domain.CreateMonster(scene.MonsterID)
		if err != nil {
			return nil, err
		}

		combat := a.Game.StartCombat(player, &monster)
		a.setCombat(telegramID, &combat)
		player.StoryState.InteractionState = domain.StateChoosingCombatAction

		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}

		intro := formatSceneHeader(player, scene) +
			"\n\nA fight begins against " + monster.Name + "!"
		return a.combatReply(&combat, intro), nil
	}

	if err := a.Game.SaveGame(ctx, player); err != nil {
		return nil, err
	}

	return a.sceneReply(player), nil
}

func (a *App) generateValidScene(
	ctx context.Context,
	player *domain.Player,
	sceneNumber int,
) (domain.GeneratedScene, error) {
	currentRegion := engine.GetRegionForScene(sceneNumber)
	combatMode := engine.GetCombatModeForScene(sceneNumber)

	fmt.Printf(
	"COMBAT BEFORE ROLL -> scene=%d mode=%q\n",
	sceneNumber,
	combatMode,
)
	
	if combatMode == domain.CombatOptional {
	roll := rand.Intn(2)

	fmt.Printf(
		"OPTIONAL COMBAT ROLL -> scene=%d roll=%d\n",
		sceneNumber,
		roll,
	)

	if roll == 0 {
		combatMode = domain.CombatForbidden
	} else {
		combatMode = domain.CombatRequired
	}

	fmt.Printf(
	"COMBAT AFTER ROLL -> scene=%d mode=%q\n",
	sceneNumber,
	combatMode,
)

}
	availableMonsters := engine.GetAvailableMonstersForRegion(
		currentRegion,
		sceneNumber,
	)

	fmt.Printf(
	"AVAILABLE MONSTERS DEBUG -> scene=%d region=%q count=%d monsters=%+v\n",
	sceneNumber,
	currentRegion,
	len(availableMonsters),
	availableMonsters,
)

	const maxAttempts = 3
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		scene, err := a.Game.GenerateScene(
			ctx,
			player,
			player.StoryState.PreviousScene,
			player.StoryState.SelectedChoice,
			player.StoryState.PreviousOutcome,
			storyGoal,
			currentRegion,
			sceneNumber,
			combatMode,
			availableMonsters,
		)
		if err == nil {
			return scene, nil
		}
		lastErr = err
	}

	return domain.GeneratedScene{}, fmt.Errorf(
		"failed after %d attempts: %w",
		maxAttempts,
		lastErr,
	)
}

func (a *App) sceneReply(player *domain.Player) *Reply {
	scene := player.StoryState.CurrentScene
	text := formatSceneHeader(player, scene)

	buttons := make([]Button, 0, len(scene.Choices)+1)
	for _, choice := range scene.Choices {
		buttons = append(buttons, Button{
			ID:    fmt.Sprintf("%s%d", PrefixChoice, choice.ID),
			Label: choice.Text,
		})
	}
	buttons = append(buttons, Button{
		ID:    CallbackInventory,
		Label: "Open inventory",
	})

	return &Reply{Text: text, Buttons: buttons}
}

func formatSceneHeader(player *domain.Player, scene domain.GeneratedScene) string {
	return fmt.Sprintf(
		"Scene %d/15 | %s\n\n%s\n\n%s\n\n❤️ HP %d/%d | 🔷 MP %d/%d | ⬆️ Lv %d",
		player.StoryState.SceneNumber,
		player.StoryState.CurrentRegion,
		scene.Title,
		scene.Description,
		player.Stats.HP,
		player.Stats.MaxHP,
		player.Stats.MP,
		player.Stats.MaxMP,
		player.Level,
	)
}

func (a *App) handleSceneInput(
	ctx context.Context,
	player *domain.Player,
	telegramID int64,
	callback string,
) (*Reply, error) {
	if callback == CallbackInventory {
		player.StoryState.InteractionState = domain.StateChoosingInventory
		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}
		return a.inventoryReply(player), nil
	}

	choiceID, ok := parsePrefixedInt(callback, PrefixChoice)
	if !ok {
		return a.sceneReply(player), nil
	}

	choice, err := a.Game.GetGeneratedChoice(player.StoryState.CurrentScene, choiceID)
	if err != nil {
		reply := a.sceneReply(player)
		reply.Text = "Invalid choice.\n\n" + reply.Text
		return reply, nil
	}

	player.StoryState.PreviousScene = player.StoryState.CurrentScene.Description
	player.StoryState.SelectedChoice = choice.Text
	if strings.TrimSpace(player.StoryState.PreviousOutcome) == "" {
		player.StoryState.PreviousOutcome = "No combat occurred in this scene."
	}
	player.StoryState.CurrentScene = domain.GeneratedScene{}
	player.StoryState.SceneNumber++

	if player.StoryState.SceneNumber > 15 {
		player.StoryState.Completed = true
		player.StoryState.InteractionState = domain.StateAdventureComplete
		if err := a.Game.SaveGame(ctx, player); err != nil {
			return nil, err
		}
		return &Reply{
			Text: "Adventure complete! You finished the journey.",
			Buttons: []Button{
				{ID: CallbackNewGame, Label: "New Game"},
			},
		}, nil
	}

	if err := a.Game.SaveGame(ctx, player); err != nil {
		return nil, err
	}

	reply, err := a.generateAndShowScene(ctx, player, telegramID)
	if err != nil {
		return nil, err
	}
	reply.Status = "Generating next scene..."
	return reply, nil
}
