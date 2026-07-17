package domain

type SceneGenerationRequest struct {
	PlayerName        string
	PlayerClass       ClassType
	PlayerLevel       int
	PreviousScene     string
	SelectedChoice    string
	PreviousOutcome   string
	StoryGoal         string
	CurrentRegion     string
	SceneNumber       int
	CombatMode CombatMode
	AvailableMonsters []AvailableMonster
}
