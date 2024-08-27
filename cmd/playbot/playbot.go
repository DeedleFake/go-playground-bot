package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"slices"

	"deedles.dev/dgutil"
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Type: discordgo.MessageApplicationCommand,
		Name: "Run Go Code",
	},
}

func setup(s *dgutil.Setup) error {
	dg := s.Session()
	dg.AddHandler(handleCommand)
	s.RegisterCommands(slices.Values(commands))

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := dgutil.Run(ctx, setup)
	if err != nil {
		slog.Error("failed", "err", err)
		os.Exit(1)
	}
}
