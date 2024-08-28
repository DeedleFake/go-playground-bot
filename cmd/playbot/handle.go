package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"deedles.dev/dgutil"
	"deedles.dev/xiter"
	"github.com/DeedleFake/go-playground-bot/internal/extract"
	"github.com/DeedleFake/go-playground-bot/internal/play"
	"github.com/bwmarrin/discordgo"
)

// buildOutput builds the content of a message from the result of the
// playground.
func buildOutput(result play.Result) string {
	var output strings.Builder

	if result.Errors != "" {
		output.WriteString("Compile errors:\n```\n")
		output.WriteString(result.Errors)
		output.WriteString("\n```")
	}

	if result.VetErrors != "" {
		output.WriteString("Vet errors:\n```\n")
		output.WriteString(result.VetErrors)
		output.WriteString("\n```")
	}

	if result.IsTest {
		fmt.Fprintf(&output, "Test failures: %v\n", result.TestsFailed)
	}

	events := xiter.SliceChunksFunc(result.Events, func(ev play.Event) string { return ev.Kind })
	for chunk := range events {
		output.WriteString(chunk[0].Kind)
		output.WriteString(":\n```\n")
		for _, ev := range chunk {
			output.WriteString(ev.Message)
		}
		output.WriteString("\n```")
	}

	return output.String()
}

// handleBlock runs the logic for a single block of Go code in a
// message.
func handleBlock(dg *discordgo.Session, i *discordgo.Interaction, block extract.CodeBlock) {
	err := dgutil.SetupResponse(dg, i)
	if err != nil {
		slog.Error("setup response", "err", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	result, err := play.Run(ctx, block.Source)
	if err != nil {
		slog.Error("run code", "err", err)
		err = dgutil.UpdateResponse(dg, i, fmt.Sprintf("Playground error:\n```\n%v\n```", err))
		if err != nil {
			slog.Error("update response", "err", err)
		}
		return
	}

	output := buildOutput(result)
	if output == "" {
		output = "No output."
	}
	err = dgutil.UpdateResponse(dg, i, output)
	if err != nil {
		slog.Error("update response", "err", err)
		return
	}
}

// handleMessage responds to a single message.
func handleMessage(dg *discordgo.Session, i *discordgo.Interaction, msg *discordgo.Message) {
	var found bool
	for block := range extract.CodeBlocks(msg.Content) {
		if block.Language != "go" && block.Language != "" {
			continue
		}

		found = true
		handleBlock(dg, i, block)
	}
	if !found {
		err := dg.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No Go code blocks found. Note that all code blocks that are not either `go` blocks or have no language specified are ignored.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error("respond to interaction", "err", err)
			return
		}
	}
}

// handleCommand handles INTERACTION_CREATE events.
func handleCommand(dg *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("command", "name", i.ApplicationCommandData().Name, "user", i.Member.User)

	data, ok := i.Data.(discordgo.ApplicationCommandInteractionData)
	if !ok {
		slog.Error("interaction data unexpected type", "data", i.Data)
		return
	}

	for _, msg := range data.Resolved.Messages {
		handleMessage(dg, i.Interaction, msg)
	}
}
