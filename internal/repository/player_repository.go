package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"go-inventory-management/internal/domain"
)

type playerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) PlayerRepository {
	return &playerRepository{
		db: db,
	}
}

func (r *playerRepository) Create(
	ctx context.Context,
	player *domain.Player,
) error {
	statsJSON, err := json.Marshal(player.Stats)
	if err != nil {
		return fmt.Errorf("failed to serialize player stats: %w", err)
	}

	inventoryJSON, err := json.Marshal(player.Inventory)
	if err != nil {
		return fmt.Errorf("failed to serialize player inventory: %w", err)
	}

	skillsJSON, err := json.Marshal(player.Skills)
	if err != nil {
		return fmt.Errorf("failed to serialize player skills: %w", err)
	}

	storyStateJSON, err := json.Marshal(player.StoryState)
	if err != nil {
		return fmt.Errorf("failed to marshal story state: %w", err)
	}

	query := `
		INSERT INTO players (
			telegram_id,
			name,
			class,
			current_scene,
			level,
			xp,
			gold,
			stats,
			skills,
			inventory,
			equipped_weapon_id,
			equipped_armor_id,
			story_state
		)
		VALUES (
			NULLIF($1::bigint, 0::bigint),
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13
		)
		RETURNING id
	`

	err = r.db.QueryRowContext(
		ctx,
		query,
		player.TelegramID,
		player.Name,
		player.Class,
		player.CurrentScene,
		player.Level,
		player.XP,
		player.Gold,
		statsJSON,
		skillsJSON,
		inventoryJSON,
		player.EquippedWeaponID,
		player.EquippedArmorID,
		storyStateJSON,
	).Scan(&player.ID)
	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	return nil
}

func (r *playerRepository) GetByID(
	ctx context.Context,
	id int64,
) (*domain.Player, error) {
	var (
		player        domain.Player
		class         string
		telegramID    sql.NullInt64
		statsJSON     []byte
		skillsJSON    []byte
		inventoryJSON []byte
		storyJSON     []byte
	)

	query := `
	SELECT
		id,
		telegram_id,
		name,
		class,
		current_scene,
		level,
		xp,
		gold,
		stats,
		skills,
		inventory,
		equipped_weapon_id,
		equipped_armor_id,
		story_state
	FROM players
	WHERE id = $1
`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&player.ID,
		&telegramID,
		&player.Name,
		&class,
		&player.CurrentScene,
		&player.Level,
		&player.XP,
		&player.Gold,
		&statsJSON,
		&skillsJSON,
		&inventoryJSON,
		&player.EquippedWeaponID,
		&player.EquippedArmorID,
		&storyJSON,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("player with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	if telegramID.Valid {
		player.TelegramID = telegramID.Int64
	}

	player.Class = domain.ClassType(class)

	if err := json.Unmarshal(statsJSON, &player.Stats); err != nil {
		return nil, fmt.Errorf("failed to deserialize player stats: %w", err)
	}

	if err := json.Unmarshal(inventoryJSON, &player.Inventory); err != nil {
		return nil, fmt.Errorf("failed to deserialize player inventory: %w", err)
	}

	if err := json.Unmarshal(skillsJSON, &player.Skills); err != nil {
		return nil, fmt.Errorf("failed to deserialize player skills: %w", err)
	}

	if len(storyJSON) > 0 {
		if err := json.Unmarshal(storyJSON, &player.StoryState); err != nil {
			return nil, fmt.Errorf(
				"failed to deserialize story state: %w",
				err,
			)
		}
	}

	return &player, nil
}

func (r *playerRepository) Update(
	ctx context.Context,
	player *domain.Player,
) error {
	statsJSON, err := json.Marshal(player.Stats)
	if err != nil {
		return fmt.Errorf("failed to serialize player stats: %w", err)
	}

	inventoryJSON, err := json.Marshal(player.Inventory)
	if err != nil {
		return fmt.Errorf("failed to serialize player inventory: %w", err)
	}

	skillsJSON, err := json.Marshal(player.Skills)
	if err != nil {
		return fmt.Errorf("failed to serialize player skills: %w", err)
	}
	storyJSON, err := json.Marshal(player.StoryState)
	if err != nil {
		return fmt.Errorf("failed to marshal story state: %w", err)
	}

	query := `
		UPDATE players
		SET
			telegram_id = NULLIF($1::bigint, 0::bigint),
			name = $2,
			class = $3,
			current_scene = $4,
			level = $5,
			xp = $6,
			gold = $7,
			stats = $8,
			skills = $9,
			inventory = $10,
			equipped_weapon_id = $11,
			equipped_armor_id = $12,
			story_state = $13,
			updated_at = NOW()
		WHERE id = $14
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		player.TelegramID,
		player.Name,
		string(player.Class),
		player.CurrentScene,
		player.Level,
		player.XP,
		player.Gold,
		statsJSON,
		skillsJSON,
		inventoryJSON,
		player.EquippedWeaponID,
		player.EquippedArmorID,
		storyJSON,
		player.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check updated rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("player with id %d not found", player.ID)
	}

	return nil
}

func (r *playerRepository) GetLatest(
	ctx context.Context,
) (*domain.Player, error) {

	var latestID int64

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id
		 FROM players
		 ORDER BY id DESC
		 LIMIT 1`,
	).Scan(&latestID)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no saved players found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get latest player id: %w", err)
	}

	return r.GetByID(ctx, latestID)
}

func (r *playerRepository) GetByTelegramID(
	ctx context.Context,
	telegramID int64,
) (*domain.Player, error) {
	var playerID int64

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT id
		FROM players
		WHERE telegram_id = $1
		`,
		telegramID,
	).Scan(&playerID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find player by telegram id: %w",
			err,
		)
	}

	return r.GetByID(ctx, playerID)
}