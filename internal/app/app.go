package app

import (
	"sync"

	"go-inventory-management/internal/domain"
	"go-inventory-management/internal/engine"
	"go-inventory-management/internal/repository"
)

type App struct {
	Game       *engine.GameEngine
	PlayerRepo repository.PlayerRepository

	combatMu sync.Mutex
	combats  map[int64]*domain.CombatState

	pendingMu sync.Mutex
	pending   map[int64]*pendingPlayer
}

type pendingPlayer struct {
	Name string
}

func NewApp(
	game *engine.GameEngine,
	playerRepo repository.PlayerRepository,
) *App {
	return &App{
		Game:       game,
		PlayerRepo: playerRepo,
		combats:    make(map[int64]*domain.CombatState),
		pending:    make(map[int64]*pendingPlayer),
	}
}

func (a *App) setCombat(telegramID int64, combat *domain.CombatState) {
	a.combatMu.Lock()
	defer a.combatMu.Unlock()
	a.combats[telegramID] = combat
}

func (a *App) getCombat(telegramID int64) *domain.CombatState {
	a.combatMu.Lock()
	defer a.combatMu.Unlock()
	return a.combats[telegramID]
}

func (a *App) clearCombat(telegramID int64) {
	a.combatMu.Lock()
	defer a.combatMu.Unlock()
	delete(a.combats, telegramID)
}

func (a *App) setPending(telegramID int64, p *pendingPlayer) {
	a.pendingMu.Lock()
	defer a.pendingMu.Unlock()
	a.pending[telegramID] = p
}

func (a *App) getPending(telegramID int64) *pendingPlayer {
	a.pendingMu.Lock()
	defer a.pendingMu.Unlock()
	return a.pending[telegramID]
}

func (a *App) clearPending(telegramID int64) {
	a.pendingMu.Lock()
	defer a.pendingMu.Unlock()
	delete(a.pending, telegramID)
}
