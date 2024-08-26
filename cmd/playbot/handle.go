package main

import (
	"log/slog"
	"strings"

	"deedles.dev/xiter"
	"github.com/DeedleFake/go-playground-bot/internal/extract"
	"github.com/DeedleFake/go-playground-bot/internal/play"
	"github.com/bwmarrin/discordgo"
)

func handleBlock(dg *discordgo.Session, i *discordgo.Interaction, block extract.CodeBlock) {
	err := setupResponse(dg, i)
	if err != nil {
		slog.Error("setup response", "err", err)
		return
	}

	result, err := play.Run(block.Source)
	if err != nil {
		slog.Error("run code", "err", err)
		return
	}
	if result.Errors != "" {
		err = updateResponse(dg, i, "Error:\n```\n"+result.Errors+"\n```")
		if err != nil {
			slog.Error("update response with error", "err", err)
		}
		return
	}

	var output strings.Builder
	events := xiter.SliceChunksFunc(result.Events, func(ev play.Event) string { return ev.Kind })
	for chunk := range events {
		output.WriteString(chunk[0].Kind)
		output.WriteString(":\n```\n")
		for _, ev := range chunk {
			output.WriteString(ev.Message)
		}
		output.WriteString("\n```")
	}

	err = updateResponse(dg, i, output.String())
	if err != nil {
		slog.Error("update response", "err", err)
		return
	}
}

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
