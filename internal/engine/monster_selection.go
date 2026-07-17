package engine

import (
	"fmt"
	"math/rand"

	"go-inventory-management/internal/domain"
)

func SelectRequiredMonster(
	monsters []domain.AvailableMonster,
) (string, error) {
	if len(monsters) == 0 {
		return "", fmt.Errorf(
			"combat is required but no monsters are available",
		)
	}

	index := rand.Intn(len(monsters))

	return monsters[index].ID, nil
}

