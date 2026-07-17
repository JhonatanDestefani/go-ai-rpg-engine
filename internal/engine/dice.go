package engine

import (
	"fmt"
	"math/rand"
	"strings"
)

type DiceRoll struct {
	Quantity int
	Sides    int
	Rolls    []int
	Total    int
}

var d6Faces = map[int]string{
	1: "⚀",
	2: "⚁",
	3: "⚂",
	4: "⚃",
	5: "⚄",
	6: "⚅",
}

func RollDice(quantity int, sides int) DiceRoll {
	result := DiceRoll{
		Quantity: quantity,
		Sides:    sides,
		Rolls:    []int{},
		Total:    0,
	}
	for i := 0; i < quantity; i++ {
		roll := rand.Intn(sides) + 1
		result.Rolls = append(result.Rolls, roll)
		result.Total += roll
	}
	return result
}

func Roll2D6() DiceRoll {
	return RollDice(2, 6)
}

func RollD20() DiceRoll {
	return RollDice(1, 20)
}

// FormatDiceFaces renders each die as a face emoji when possible (d6),
// otherwise falls back to the numeric value.
func FormatDiceFaces(roll DiceRoll) string {
	parts := make([]string, 0, len(roll.Rolls))
	for _, value := range roll.Rolls {
		if face, ok := d6Faces[value]; ok && roll.Sides == 6 {
			parts = append(parts, face)
			continue
		}
		parts = append(parts, fmt.Sprintf("%d", value))
	}
	return strings.Join(parts, " ")
}

// FormatRollLine formats one contested roll line, e.g. "Attack: ⚂ ⚄ = 8".
func FormatRollLine(label string, roll DiceRoll) string {
	return fmt.Sprintf("%s: %s = %d", label, FormatDiceFaces(roll), roll.Total)
}

// FormatContestedRolls formats attack and defense roll lines together.
func FormatContestedRolls(attackRoll, defenseRoll DiceRoll) string {
	return FormatRollLine("Attack", attackRoll) + "\n" +
		FormatRollLine("Defense", defenseRoll)
}
