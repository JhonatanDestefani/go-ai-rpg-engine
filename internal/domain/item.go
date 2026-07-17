package domain

type ItemType string

const (
	ItemConsumable ItemType = "consumable"
	ItemWeapon     ItemType = "weapon"
	ItemArmor      ItemType = "armor"
	KeyItem        ItemType = "keyitem"
)

type ItemRarity string

const (
	RarityCommon   ItemRarity = "common"
	RarityUncommon ItemRarity = "uncommon"
	RarityRare     ItemRarity = "rare"
	RarityUnique   ItemRarity = "unique"
)

type Item struct {
	ID          string
	Name        string
	Description string
	Type        ItemType
	Rarity      ItemRarity
	Value       int
	Quantity    int

	Attack  int
	Magic int
	Defense int
	Heal    int
	MPHeal  int
}
