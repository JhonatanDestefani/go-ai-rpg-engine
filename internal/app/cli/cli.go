package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go-inventory-management/internal/app"
)

// Run drives the shared session handlers through the terminal.
func Run(ctx context.Context, application *app.App) error {
	reader := bufio.NewReader(os.Stdin)
	telegramID := app.CLITelegramID

	fmt.Println("=== RPG CLI mode (APP_MODE=cli) ===")
	fmt.Println("Type /start to begin, or quit to exit.")
	fmt.Println()

	reply, err := application.HandleStart(ctx, telegramID)
	if err != nil {
		return err
	}
	printReply(reply)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.EqualFold(line, "quit") || strings.EqualFold(line, "exit") {
			fmt.Println("Bye.")
			return nil
		}
		if line == "/start" {
			reply, err = application.HandleStart(ctx, telegramID)
			if err != nil {
				return err
			}
			printReply(reply)
			continue
		}

		input := app.Input{}
		if reply != nil && len(reply.Buttons) > 0 {
			if n, err := strconv.Atoi(line); err == nil {
				if n >= 1 && n <= len(reply.Buttons) {
					input.CallbackID = reply.Buttons[n-1].ID
				} else {
					fmt.Println("Invalid option.")
					printReply(reply)
					continue
				}
			} else {
				// Free text (e.g. name) even when buttons are shown is rare;
				// treat non-numeric as text.
				input.Text = line
			}
		} else {
			input.Text = line
		}

		reply, err = application.HandleInput(ctx, telegramID, input)
		if err != nil {
			return err
		}
		printReply(reply)
	}
}

func printReply(reply *app.Reply) {
	if reply == nil {
		return
	}
	if reply.Status != "" {
		fmt.Println(reply.Status)
	}
	fmt.Println()
	fmt.Println(reply.Text)
	fmt.Println()
	for i, button := range reply.Buttons {
		fmt.Printf("%d - %s\n", i+1, button.Label)
	}
}
