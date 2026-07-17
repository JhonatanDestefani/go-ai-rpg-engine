package engine

import (
	"errors"
	"fmt"
	"strings"

	"go-inventory-management/internal/domain"
)


func ValidateGeneratedScene(
	scene domain.GeneratedScene,
	availableMonsters []domain.AvailableMonster,
) error {
	if strings.TrimSpace(scene.Title) == "" {
		return errors.New("generated scene has no title")
	}

	if strings.TrimSpace(scene.Description) == "" {
		return errors.New("generated scene has no description")
	}

	if len(scene.Choices) == 0 {
		return errors.New("generated scene has no choices")
	}

	if len(scene.Choices) > 3 {
		return fmt.Errorf(
			"generated scene has too many choices: %d",
			len(scene.Choices),
		)
	}

	if err := validateGeneratedChoices(scene.Choices); err != nil {
		return err
	}

	if err := validateGeneratedMonsterID(
		scene.MonsterID,
		availableMonsters,
	); err != nil {
		return err
	}

	return nil
}
func validateGeneratedChoices(choices []domain.GeneratedChoice) error {
	usedIDs := make(map[int]struct{}, len(choices))

	for index, choice := range choices {
		expectedID := index + 1

		if choice.ID != expectedID {
			return fmt.Errorf(
				"invalid choice ID: expected %d, received %d",
				expectedID,
				choice.ID,
			)
		}

		if strings.TrimSpace(choice.Text) == "" {
			return fmt.Errorf(
				"choice %d has no text",
				choice.ID,
			)
		}

		if _, exists := usedIDs[choice.ID]; exists {
			return fmt.Errorf(
				"duplicated choice ID: %d",
				choice.ID,
			)
		}

		usedIDs[choice.ID] = struct{}{}
	}

	return nil
}
func validateGeneratedMonsterID(
	monsterID string,
	availableMonsters []domain.AvailableMonster,
) error {
	monsterID = strings.TrimSpace(monsterID)

	if monsterID == "" {
		return nil
	}

	for _, monster := range availableMonsters {
		if monster.ID == monsterID {
			return nil
		}
	}

	return fmt.Errorf(
		"Ollama returned an unavailable monster ID: %q",
		monsterID,
	)
}
