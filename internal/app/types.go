package app

// Button is a choice presented to the player.
// Telegram maps ID to callback_data; CLI shows numbered labels.
type Button struct {
	ID    string
	Label string
}

// Reply is what an adapter should show after handling one input.
type Reply struct {
	Text     string
	Buttons  []Button
	Status   string // optional short status (e.g. "Generating scene...")
	ClearUI  bool
}

// Input is one user action from Telegram or the terminal.
type Input struct {
	Text       string
	CallbackID string
}

const (
	CallbackContinue   = "start:continue"
	CallbackNewGame    = "start:new"
	CallbackClassWarrior = "class:warrior"
	CallbackClassMage  = "class:mage"
	CallbackAttack     = "combat:attack"
	CallbackSkillMenu  = "combat:skill"
	CallbackItemMenu   = "combat:item"
	CallbackInventory  = "scene:inventory"
	CallbackInvEquip   = "inv:equip"
	CallbackInvUse     = "inv:use"
	CallbackInvBack    = "inv:back"
	CallbackEquipBack  = "equip:back"
	CallbackSkillBack  = "skill:back"
	CallbackItemBack   = "item:back"

	PrefixChoice = "choice:"
	PrefixSkill  = "skill:"
	PrefixItem   = "item:"
	PrefixEquip  = "equip:"
)

// CLITelegramID is the synthetic telegram_id used by APP_MODE=cli.
const CLITelegramID int64 = 1

const storyGoal = `
The player must cross the Dark Forest, investigate the Ancient Ruins,
discover why monsters are gathering near the Sanctuary,
and eventually face the Dragon.

The adventure must contain at least 15 scenes.
The story must build tension gradually.
The Dragon is the final boss and must only appear in scene 15.
Not every scene should contain combat.
Exploration, discoveries and consequences must advance the story.
`
