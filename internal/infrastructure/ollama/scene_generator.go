package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go-inventory-management/internal/domain"
)

type SceneGenerator struct {
	baseURL string
	model   string
	client  *http.Client
}

type generateRequest struct {
	Model   string          `json:"model"`
	System  string          `json:"system"`
	Prompt  string          `json:"prompt"`
	Format  map[string]any  `json:"format"`
	Stream  bool            `json:"stream"`
	Think bool				`json:"think"`
	Options generateOptions `json:"options"`
}

type generateOptions struct {
	Temperature float64 `json:"temperature"`
}

type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewSceneGenerator(baseURL, model string) *SceneGenerator {
	return &SceneGenerator{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
	}
}

func sceneJSONSchema(monsters []domain.AvailableMonster,
	combatMode domain.CombatMode,) map[string]any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required": []string{
			"title",
			"description",
			"monster_id",
			"choices",
		},
		"properties": map[string]any{
			"title": map[string]any{
				"type":      "string",
				"minLength": 1,
			},
			"description": map[string]any{
				"type":      "string",
				"minLength": 1,
			},
			"monster_id": map[string]any{
				"type": "string",
			},
			"choices": map[string]any{
				"type":     "array",
				"minItems": 1,
				"maxItems": 3,
				"items": map[string]any{
					"type":                 "object",
					"additionalProperties": false,
					"required": []string{
						"id",
						"text",
					},
					"properties": map[string]any{
						"id": map[string]any{
							"type":    "integer",
							"minimum": 1,
						},
						"text": map[string]any{
							"type":      "string",
							"minLength": 1,
						},
					},
				},
			},
		},
	}
}

func buildSystemPrompt() string {
	return `You are the scene generator for a turn-based fantasy RPG.

Rules:
- Write all scene content in English.
- Return only data matching the provided JSON schema.
- Never create monsters.
- You may only use a monster_id present in the available monsters list.
- Use an empty monster_id when the scene has no monster.
- Do not describe combat results.
- Do not change player statistics.
- Do not create items, rewards, experience or gold.
- Create between 1 and 3 meaningful choices.
- Choice IDs must start at 1 and be sequential.
- Keep the scene concise.
- Never return the text "null" as monster_id.
- monster_id must always be a JSON string.
- When no monster exists, use exactly "".
- When combat is required, use exactly one ID from the available monsters list.
- Continue naturally from the previous scene and selected choice.`
}

func buildPrompt(request domain.SceneGenerationRequest) (string, error) {
	monstersJSON, err := json.MarshalIndent(
		request.AvailableMonsters,
		"",
		"  ",
	)
	if err != nil {
		return "", fmt.Errorf("failed to marshal available monsters: %w", err)
	}

	previousScene := request.PreviousScene
	if previousScene == "" {
		previousScene = "This is the beginning of the adventure."
	}

	selectedChoice := request.SelectedChoice
	if selectedChoice == "" {
		selectedChoice = "The player has not selected a previous choice."
	}
	previousOutcome := request.PreviousOutcome
	if strings.TrimSpace(previousOutcome) == "" {
	previousOutcome = "No previous outcome."
}
	combatMode := request.CombatMode

	combatInstruction := ""

	switch combatMode {
	case domain.CombatForbidden:
	combatInstruction = `
	Combat is forbidden in this scene.
	You must return an empty monster_id.
	Create an exploration, introduction, discovery or decision scene.
	`

	case domain.CombatRequired:
	combatInstruction = `
	Combat is mandatory in this scene.

	You must choose exactly one monster_id from the available monsters list.

	IMPORTANT:
	- Do not return an empty monster_id.
	- Do not return "null".
	- Do not return null.
	- Do not invent or modify an ID.
	- Copy the chosen monster ID exactly, character by character.

	The monster encounter begins automatically.
	The choices must describe actions available after the monster is defeated.
	`

default:
	combatInstruction = `
	Combat is optional in this scene.
	Use a monster only when it meaningfully advances the story.
	Otherwise return an empty monster_id.
	`
}
	regionInstruction := buildRegionInstruction(
	request.SceneNumber,
)
	

	return fmt.Sprintf(
	`Create scene number %d of the adventure.

Story goal:
%s

Current region:
%s

Region rules:
%s

Player:
- Name: %s
- Class: %s
- Level: %d

Previous scene:
%s

Choice selected by the player:
%s

Confirmed outcome:
%s

Combat rules for this scene:
%s

Available monsters:
%s

Important rules, in priority order:

1. The confirmed outcome is factual and must never be contradicted, changed, reinterpreted, or ignored.
2. The next scene must be a direct consequence of:
   - the confirmed outcome;
   - the player's selected choice;
   - the current region;
   - the previous scene.
3. Do not repeat, paraphrase, summarize, or recreate the previous scene.
   The new scene must introduce at least one new event, discovery, obstacle, consequence, character interaction, or environmental change.
4. The player's selected choice has already been performed.
   Do not present the same action, destination, object, or investigation as still available.
5. Begin the new scene after the selected action has occurred.
   Do not describe the player deciding whether to perform the selected choice.
6. Reuse previous-scene elements only when showing a new consequence.
   Do not reuse the same description, purpose, or set of available actions.
7. The new scene must contain at least one concrete consequence that would not exist if the player had selected a different choice.
8. Every choice must be physically possible in the location described by the current scene.
9. Every choice must be consistent with the current region's biome, atmosphere, enemies, landmarks, and rules.
10. Never mention, reveal, foreshadow, or move the player toward locations belonging to future regions unless an explicit transition is allowed by the region rules.
11. Do not create a biome transition unless:
   - the region rules explicitly allow it;
   - the current scene number permits it;
   - the engine-provided transition condition has been met.
12. Choices must be meaningfully different.
   Do not provide two choices that lead to the same action, destination, or consequence.
13. The player's selected choice must meaningfully affect the next scene.
   The next scene must clearly reflect what the player chose.
14. Only use a monster_id exactly as provided in the available monster list.
    Never invent, alter, translate, shorten, or approximate a monster_id.
15. If there is no combat in the scene, monster_id must be an empty string.
16. If combat is required, monster_id must not be empty and must match one of the available monster IDs.
17. Do not create enemies, bosses, creatures, NPCs, landmarks, items, or events that violate the current region rules.
18. Follow the output format exactly.
    Do not add explanations, markdown, comments, notes, or text outside the requested JSON.`,
	request.SceneNumber,
	request.StoryGoal,
	request.CurrentRegion,
	regionInstruction,
	request.PlayerName,
	request.PlayerClass,
	request.PlayerLevel,
	previousScene,
	selectedChoice,
	previousOutcome,
	combatInstruction,
	string(monstersJSON),
), nil
	}

func (g *SceneGenerator) GenerateScene(
	ctx context.Context,
	request domain.SceneGenerationRequest,
) (domain.GeneratedScene, error) {
	prompt, err := buildPrompt(request)
	if err != nil {
		return domain.GeneratedScene{}, err
	}

	payload := generateRequest{
		Model:  g.model,
		System: buildSystemPrompt(),
		Prompt: prompt,
		Format: sceneJSONSchema(
		request.AvailableMonsters,
		request.CombatMode,
),
		Stream: false,
		Think:  false,
		Options: generateOptions{
			Temperature: 0.2,
		},
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return domain.GeneratedScene{}, fmt.Errorf(
			"failed to marshal Ollama request: %w",
			err,
		)
	}

	endpoint := g.baseURL + "/api/generate"

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return domain.GeneratedScene{}, fmt.Errorf(
			"failed to create Ollama request: %w",
			err,
		)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := g.client.Do(httpRequest)
	if err != nil {
		return domain.GeneratedScene{}, fmt.Errorf(
			"failed to call Ollama: %w",
			err,
		)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.GeneratedScene{}, fmt.Errorf(
			"failed to read Ollama response: %w",
			err,
		)
	}

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		var ollamaError errorResponse

		if err := json.Unmarshal(body, &ollamaError); err == nil &&
			ollamaError.Error != "" {
			return domain.GeneratedScene{}, fmt.Errorf(
				"Ollama returned status %d: %s",
				response.StatusCode,
				ollamaError.Error,
			)
		}

		return domain.GeneratedScene{}, fmt.Errorf(
			"Ollama returned status %d: %s",
			response.StatusCode,
			string(body),
		)
	}

	var result generateResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return domain.GeneratedScene{}, fmt.Errorf(
			"failed to decode Ollama response: %w",
			err,
		)
	}

	fmt.Println("========== RAW OLLAMA RESPONSE ==========")
	fmt.Println(result.Response)
	fmt.Println("=========================================")

	if strings.TrimSpace(result.Response) == "" {
		return domain.GeneratedScene{}, fmt.Errorf(
			"Ollama returned an empty response",
		)
	}

	var scene domain.GeneratedScene

	if err := json.Unmarshal(
		[]byte(result.Response),
		&scene,
	); err != nil {
		return domain.GeneratedScene{}, fmt.Errorf(
			"failed to decode generated scene: %w; response: %s",
			err,
			result.Response,
		)
	}
	scene.MonsterID = normalizeMonsterID(scene.MonsterID)

	return scene, nil
}
func buildRegionInstruction(sceneNumber int) string {
	switch {
	case sceneNumber >= 1 && sceneNumber <= 5:
		return `
CURRENT REGION: DARK FOREST
SCENES: 1 TO 5

The scene takes place entirely inside the Dark Forest.

Allowed elements:
- dense forest paths
- abandoned camps
- rivers and old wooden bridges
- caves hidden by vegetation
- hunter trails
- fallen trees
- fog, moss and animal tracks

Forbidden:
- temples
- ancient stone sanctuaries
- lava
- dragon chambers
- final boss locations

All choices must remain inside the Dark Forest.
Scene 5 may naturally lead toward the Ancient Ruins.
`

	case sceneNumber >= 6 && sceneNumber <= 10:
		return `
CURRENT REGION: ANCIENT RUINS
SCENES: 6 TO 10

The player has permanently left the Dark Forest.
The player is already inside the Ancient Ruins.

Every description and every choice must involve:
- collapsed stone corridors
- underground chambers
- ancient statues
- broken gates
- forgotten libraries
- traps
- old mechanisms
- inscriptions
- dragon symbols
- sealed doors
- ruined stairways

ABSOLUTELY FORBIDDEN:
- forests
- trees
- branches
- leaves
- animals
- animal tracks
- rivers
- riverbanks
- wooden bridges
- camps
- wilderness exploration

Do not continue clues or environmental elements from the Dark Forest.
Scene 10 may naturally lead toward the Sanctuary.
`

	case sceneNumber >= 11 && sceneNumber <= 15:
		return `
CURRENT REGION: SANCTUARY
SCENES: 11 TO 15

The player has permanently left the Ancient Ruins.
The player is already inside the Sanctuary.

Every description and every choice must involve:
- sacred halls
- white stone chambers
- crystals
- holy statues
- magical energy
- divine light
- ceremonial rooms
- sealed sacred doors
- the path toward the final chamber

ABSOLUTELY FORBIDDEN:
- forests
- trees
- branches
- leaves
- animals
- animal tracks
- rivers
- riverbanks
- bridges
- camps
- caves
- ordinary ruin exploration


Do not continue clues or environmental elements from previous regions.

The Fire Dragon may appear only in scene 15.
Scenes 11 to 14 must advance toward the final chamber.
Scene 15 takes place in the final chamber.
`

	default:
		return `
The scene number is invalid.
Do not generate a scene.
`
	}
}

func normalizeMonsterID(monsterID string) string {
	monsterID = strings.TrimSpace(monsterID)

	for {
		trimmed := strings.Trim(monsterID, `"`)
		if trimmed == monsterID {
			break
		}
		monsterID = strings.TrimSpace(trimmed)
	}

	switch strings.ToLower(monsterID) {
	case "", "null", "nil", "none":
		return ""
	default:
		return monsterID
	}
}