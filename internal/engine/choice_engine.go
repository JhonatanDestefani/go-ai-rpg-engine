package engine

import (
	"go-inventory-management/internal/repository"
)

type GameEngine struct {
	playerRepository repository.PlayerRepository
	sceneGenerator   SceneGenerator
}

func NewGameEngine(
	playerRepository repository.PlayerRepository,
	sceneGenerator SceneGenerator,
) *GameEngine {
	return &GameEngine{
		playerRepository: playerRepository,
		sceneGenerator:   sceneGenerator,
	}
}
