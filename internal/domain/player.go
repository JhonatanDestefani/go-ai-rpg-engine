package domain

type Player struct {
	ID           int64
	TelegramID   int64
	Name         string
	Class        ClassType
	CurrentScene string
	Level        int
	XP           int
	Gold         int
	Stats        Stats
	Skills       []Skill
	Inventory    Inventory

	EquippedWeaponID string
	EquippedArmorID  string

	StoryState StoryState
}
