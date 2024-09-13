package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime/pprof"
	"slices"

	"deedles.dev/dgutil"
	"github.com/bwmarrin/discordgo"
	"github.com/thejerf/suture/v4"
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

type App struct{}

func (a *App) Serve(ctx context.Context) error {
	bot := dgutil.Bot{
		Commands: slices.Values(commands),
	}

	dg, err := bot.Session()
	if err != nil {
		return fmt.Errorf("initialize Discord session: %w", err)
	}
	dgutil.AddHandler(ctx, dg, handleCommand)

	err = bot.Run(ctx)
	if err != nil {
		return fmt.Errorf("run bot: %w", err)
	}

	return nil
}

func main() {
	defer profile()()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	sup := suture.NewSimple("playbot")
	sup.Add(&App{})
	err := sup.Serve(ctx)
	if err != nil {
		slog.Error("supervisor tree failed", "err", err)
		os.Exit(1)
	}
}
