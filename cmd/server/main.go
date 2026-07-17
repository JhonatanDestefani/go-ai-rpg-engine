package main

import (
	"context"
	"log"
	"strings"

	"go-inventory-management/internal/app"
	"go-inventory-management/internal/app/cli"
	"go-inventory-management/internal/config"
	"go-inventory-management/internal/database"
	"go-inventory-management/internal/engine"
	"go-inventory-management/internal/infrastructure/ollama"
	"go-inventory-management/internal/repository"
	"go-inventory-management/internal/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresConnection(cfg.DBConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	playerRepo := repository.NewPlayerRepository(db)
	sceneGenerator := ollama.NewSceneGenerator(cfg.OllamaURL, cfg.OllamaModel)
	game := engine.NewGameEngine(playerRepo, sceneGenerator)
	application := app.NewApp(game, playerRepo)

	ctx := context.Background()
	mode := strings.ToLower(strings.TrimSpace(cfg.AppMode))

	switch mode {
	case "cli":
		log.Println("Starting in CLI mode")
		if err := cli.Run(ctx, application); err != nil {
			log.Fatal(err)
		}
	default:
		if cfg.TelegramBotToken == "" {
			log.Fatal("TELEGRAM_BOT_TOKEN is not configured (or set APP_MODE=cli)")
		}
		log.Println("Starting Telegram bot")
		if err := telegram.Start(ctx, cfg.TelegramBotToken, application); err != nil {
			log.Fatal(err)
		}
	}
}
