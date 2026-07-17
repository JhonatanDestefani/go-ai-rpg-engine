package domain

type CombatState struct {
	Player  *Player
	Monster *Monster

	IsPlayerTurn bool
	IsFinished   bool
	PlayerWon    bool

	PlayerEffects  []CombatEffect
	MonsterEffects []CombatEffect

	IsActive bool
}
