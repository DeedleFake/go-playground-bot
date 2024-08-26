package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func setupResponse(dg *discordgo.Session, i *discordgo.Interaction) error {
	return dg.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Working on it...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func updateResponse(dg *discordgo.Session, i *discordgo.Interaction, content string) error {
	_, err := dg.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &content,
	})
	return err
}

func commands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Type: discordgo.MessageApplicationCommand,
			Name: "Run Go Code",
		},
	}
}

func run(ctx context.Context) error {
	token, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		return errors.New("$DISCORD_TOKEN not set")
	}

	dg, err := discordgo.New("Bot " + strings.TrimSpace(token))
	if err != nil {
		return fmt.Errorf("create Discord session: %w", err)
	}
	dg.AddHandler(func(dg *discordgo.Session, r *discordgo.Ready) {
		slog.Info("authenticated successfully", "user", r.User)
	})
	dg.AddHandler(handleCommand)

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("open Discord session: %w", err)
	}
	defer dg.Close()

	for _, guild := range dg.State.Guilds {
		for _, cmd := range commands() {
			r, err := dg.ApplicationCommandCreate(dg.State.User.ID, guild.ID, cmd)
			if err != nil {
				return fmt.Errorf("register command %q: %w", cmd.Name, err)
			}
			defer func() {
				err := dg.ApplicationCommandDelete(dg.State.User.ID, guild.ID, r.ID)
				if err != nil {
					slog.Error("unregister command", "command", r.Name, "err", err)
				}
			}()

			slog.Info("command registered", "command", r.Name, "guild_id", guild.ID, "guild_name", guild.Name)
		}
	}

	<-ctx.Done()
	slog.Info("exiting")

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := run(ctx)
	if err != nil {
		slog.Error("failed", "err", err)
		os.Exit(1)
	}
}
