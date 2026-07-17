package domain

type GeneratedScene struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	MonsterID   string            `json:"monster_id"`
	Choices     []GeneratedChoice `json:"choices"`
}

type GeneratedChoice struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type AvailableMonster struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
