package domain

type DropTable struct {
	MonsterID string
	Drops     []Drop
}

func CreateDropTable(
	monsterID string,
	drops ...Drop,
) DropTable {
	return DropTable{
		MonsterID: monsterID,
		Drops:     drops,
	}
}

func NewDrop(item Item, min, max, qtyMin, qtyMax int) Drop {
	return Drop{
		Item:        item,
		DiceMin:     min,
		DiceMax:     max,
		QuantityMin: qtyMin,
		QuantityMax: qtyMax,
	}
}
