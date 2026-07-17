package engine

import (
	"context"

	"go-inventory-management/internal/domain"
)

type SceneGenerator interface {
	GenerateScene(
		ctx context.Context,
		request domain.SceneGenerationRequest,
	) (domain.GeneratedScene, error)
}
