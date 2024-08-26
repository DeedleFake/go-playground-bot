package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"deedles.dev/xiter"
	"github.com/bwmarrin/discordgo"
)

func handleCommand(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("command", "name", i.ApplicationCommandData().Name, "user", i.Member.User)

	data, ok := i.Data.(discordgo.ApplicationCommandInteractionData)
	if !ok {
		slog.Error("interaction data unexpected type", "data", i.Data)
		return
	}

	for _, msg := range data.Resolved.Messages {
		result, err := CompileAndRun(msg.Content, false)
		if err != nil {
			slog.Error("run code", "err", err)
			continue
		}

		var rsp playgroundResponse
		err = json.Unmarshal(result, &rsp)
		if err != nil {
			slog.Error("decode response", "result", result, "err", err)
			continue
		}
		events := xiter.SliceChunksFunc(rsp.Events, func(ev Event) string { return ev.Kind })

		var output strings.Builder
		for chunk := range events {
			output.WriteString(chunk[0].Kind)
			output.WriteString(":\n```\n")
			for _, ev := range chunk {
				output.WriteString(ev.Message)
			}
			output.WriteString("\n```")
		}

		err = dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: output.String(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error("respond to interaction", "err", err)
			continue
		}
	}
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
