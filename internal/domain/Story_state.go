package domain

type StoryState struct {
	SceneNumber      int              `json:"scene_number"`
	CurrentRegion    string           `json:"current_region"`
	CurrentScene     GeneratedScene   `json:"current_scene"`
	PreviousScene    string           `json:"previous_scene"`
	SelectedChoice   string           `json:"selected_choice"`
	PreviousOutcome  string           `json:"previous_outcome"`
	Completed        bool             `json:"completed"`
	InteractionState InteractionState `json:"interaction_state"`
}

type InteractionState string

const (
	StateChoosingStart        InteractionState = "choosing_start"
	StateChoosingName         InteractionState = "choosing_name"
	StateChoosingClass        InteractionState = "choosing_class"
	StateChoosingScene        InteractionState = "choosing_scene"
	StateChoosingCombatAction InteractionState = "choosing_combat_action"
	StateChoosingSkill        InteractionState = "choosing_skill"
	StateChoosingItem         InteractionState = "choosing_item"
	StateChoosingInventory    InteractionState = "choosing_inventory"
	StateChoosingEquip        InteractionState = "choosing_equip"
	StateAdventureComplete    InteractionState = "adventure_complete"
)
