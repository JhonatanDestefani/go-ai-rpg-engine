package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-inventory-management/internal/domain"
)

func (g *GameEngine) GetGeneratedChoice(
	scene domain.GeneratedScene,
	choiceID int,
) (domain.GeneratedChoice, error) {
	for _, choice := range scene.Choices {
		if choice.ID == choiceID {
			return choice, nil
		}
	}

	return domain.GeneratedChoice{}, errors.New(
		"invalid generated choice",
	)
}

func (g *GameEngine) GenerateScene(
	ctx context.Context,
	player *domain.Player,
	previousScene string,
	selectedChoice string,
	previousOutcome string,
	storyGoal string,
	currentRegion string,
	sceneNumber int,
	combatMode domain.CombatMode,
	availableMonsters []domain.AvailableMonster,
) (domain.GeneratedScene, error) {
	request := domain.SceneGenerationRequest{
		PlayerName:        player.Name,
		PlayerClass:       player.Class,
		PlayerLevel:       player.Level,
		PreviousScene:     previousScene,
		SelectedChoice:    selectedChoice,
		PreviousOutcome:   previousOutcome,
		StoryGoal:         storyGoal,
		CurrentRegion:     currentRegion,
		SceneNumber:       sceneNumber,
		CombatMode:        combatMode,
		AvailableMonsters: availableMonsters,
	}

	fmt.Println("===== AVAILABLE MONSTERS =====")
	for _, m := range request.AvailableMonsters {
	fmt.Println(m.ID)
	}
	fmt.Println("==============================")

	scene, err := g.sceneGenerator.GenerateScene(ctx, request)
	if err != nil {
		return domain.GeneratedScene{}, fmt.Errorf(
			"failed to generate scene: %w",
			err,
		)
	}

	switch request.CombatMode {
	case domain.CombatForbidden:
	scene.MonsterID = ""

	case domain.CombatRequired:
	if strings.TrimSpace(scene.MonsterID) == "" ||
		strings.EqualFold(scene.MonsterID, "null") ||
		scene.MonsterID == `""` {

		monsterID, err := SelectRequiredMonster(request.AvailableMonsters)
		if err != nil {
			return domain.GeneratedScene{}, err
		}

		scene.MonsterID = monsterID
	}
}

	if err := ValidateGeneratedScene(
		scene,
		availableMonsters,
	); err != nil {
		return domain.GeneratedScene{}, fmt.Errorf(
			"invalid generated scene: %w",
			err,
		)
	}
	player.StoryState.SceneNumber = sceneNumber
	player.StoryState.CurrentRegion = currentRegion
	player.StoryState.CurrentScene = scene
	player.StoryState.PreviousScene = previousScene
	player.StoryState.SelectedChoice = selectedChoice
	player.StoryState.PreviousOutcome = previousOutcome
	player.StoryState.Completed = false

	return scene, nil
}