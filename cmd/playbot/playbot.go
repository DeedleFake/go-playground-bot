package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime/pprof"
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

func profile() func() {
	path, ok := os.LookupEnv("PPROF")
	if !ok {
		return func() {}
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	err = pprof.StartCPUProfile(file)
	if err != nil {
		panic(err)
	}

	return func() {
		pprof.StopCPUProfile()
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	defer profile()()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bot := dgutil.Bot{
		Commands: slices.Values(commands),
	}

	dg, err := bot.Session()
	if err != nil {
		slog.Error("failed to initialize Discord session", "err", err)
		os.Exit(1)
	}
	dgutil.AddHandler(ctx, dg, handleCommand)

	err = bot.Run(ctx)
	if err != nil {
		slog.Error("failed", "err", err)
		os.Exit(1)
	}
}
