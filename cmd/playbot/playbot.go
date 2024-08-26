package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

//func playground(cfg *config, s *discordgo.Session, m *discordgo.Message, res *parsingResult) {
//	needHelp := findBoolOption(res.options, "help", "h")
//	if needHelp {
//		help(cfg, s, m, res)
//
//		return
//	}
//
//	debug := findBoolOption(res.options, "debug", "d", "explain", "e")
//	b, err := CompileAndRun(res.content, debug)
//	if err != nil {
//		sendDeletable(s, m, fmt.Sprintf("```\n%v```", err), 5*time.Minute)
//	}
//
//	var response playgroundResponse
//
//	err = json.Unmarshal(b, &response)
//	if err != nil {
//		log.Println(err)
//
//		return
//	}
//
//	if len(response.Errors) > 0 && len(response.Events) == 0 {
//		sendDeletable(s, m, fmt.Sprintf("```go\n%v```", response.Errors), 5*time.Minute)
//
//		return
//	}
//
//	plain := findBoolOption(res.options, "plain", "p")
//	if plain {
//		result := ""
//		for _, e := range response.Events {
//			if len(result) >= 2000 {
//				break
//			}
//
//			if !isPrintable(e.Message) {
//				continue
//			}
//
//			result = fmt.Sprintf("%s\n_%s_\n%s\n", result, e.Kind, e.Message)
//		}
//
//		if len(response.Events) == 0 {
//			result = "There's nothing to print out.\nReact with 😐 to delete this message."
//		}
//
//		if len(response.Errors) > 0 {
//			result = response.Errors + result
//		}
//
//		const plainOutputTempalte = "*Result*:\n```\n%s\n```"
//
//		if len(result) > 2000-len(plainOutputTempalte) {
//			result = result[:2000-len(plainOutputTempalte)]
//		}
//
//		result = fmt.Sprintf(plainOutputTempalte, result)
//
//		sendDeletable(s, m, result, 5*time.Minute)
//
//		return
//	}
//
//	length := 0
//	emb := &discordgo.MessageEmbed{
//		Title: "Result:",
//	}
//
//	length += len("Result:")
//
//	for _, e := range response.Events {
//		switch {
//		case len(e.Message) > 1024:
//			e.Message = e.Message[:1024]
//		case len(e.Message) == 0:
//			continue
//		}
//
//		if !isPrintable(e.Message) {
//			continue
//		}
//
//		length += len(e.Kind)
//		length += len(e.Message)
//
//		if length > 6000-len("Message is too long...") {
//			emb.Description = "Message is too long..."
//			length += len(e.Message)
//
//			break
//		}
//
//		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
//			Name:  e.Kind,
//			Value: e.Message,
//		})
//
//		if len(emb.Fields) == 25 {
//			emb.Description = "The maximum field amount is 25.\nThe result will be cut off..."
//
//			break
//		}
//	}
//
//	if len(response.Errors) > 0 && len(response.Errors)+length < 6000 {
//		emb.Description = fmt.Sprintf("```go\n%s\n```", response.Errors)
//	}
//
//	if len(response.Events) == 0 {
//		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
//			Name:  "success",
//			Value: "There's nothing to print out.\nReact with 😐 to delete this message.",
//		})
//	}
//
//	sendDeletable(s, m, emb, 5*time.Minute)
//}
//
//func commandHandler(cfg *config, s *discordgo.Session, msg any) func() {
//	var (
//		content string
//		pmsg    *discordgo.Message
//	)
//
//	switch m := msg.(type) {
//	case *discordgo.MessageCreate:
//		if m.Author.Bot {
//			return nil
//		}
//
//		pmsg = m.Message
//		content = m.Content
//
//	case *discordgo.MessageUpdate:
//		if m.Author != nil && m.Author.Bot {
//			return nil
//		}
//
//		pmsg = m.Message
//		content = m.Content
//
//	default:
//		return nil
//	}
//
//	content = catchPrefix(content, cfg.prefix, cfg.botID)
//	if content == "" {
//		return nil
//	}
//
//	res := parseCommand(content, " \t\n", []string{"-", "--"}, []string{"="})
//	if len(res.command) < 1 {
//		return nil
//	}
//
//	command, ok := cfg.commands[res.command]
//	if !ok {
//		return nil
//	}
//
//	return func() {
//		command(cfg, s, pmsg, res)
//	}
//
//}
//
//func help(cfg *config, s *discordgo.Session, m *discordgo.Message, res *parsingResult) {
//	sendDeletable(s, m, "```\nNo help, no hope, human. But if you like, just write it down yourself and tag @English Learner, they're in charge on me.\n"+
//		"Well, basically, I evaluate a code, then give the result of it and stuff. Use go command and get them!\n"+
//		"Btw, react with 😐 within 5 mins to rid of anything I reply to you.\n```", 5*time.Minute)
//}
//
//func findBoolOption(m map[string]any, variants ...string) bool {
//	for _, v := range variants {
//		r, ok := m[v].(bool)
//		if ok {
//			return r
//		}
//	}
//
//	return false
//}
//
//func sendDeletable(s *discordgo.Session, ctx *discordgo.Message, content any, delay time.Duration) {
//	var (
//		msg *discordgo.Message
//		err error
//	)
//
//	ref := &discordgo.MessageReference{
//		MessageID: ctx.ID,
//		ChannelID: ctx.ChannelID,
//		GuildID:   ctx.GuildID,
//	}
//
//	switch c := content.(type) {
//	case string:
//		msg, err = s.ChannelMessageSendComplex(ctx.ChannelID, &discordgo.MessageSend{
//			Content:   c,
//			Reference: ref,
//		})
//	case *discordgo.MessageEmbed:
//		msg, err = s.ChannelMessageSendComplex(ctx.ChannelID, &discordgo.MessageSend{
//			Embed:     c,
//			Reference: ref,
//		})
//	default:
//		return
//	}
//
//	if err != nil {
//		log.Println("sendDeletable:", err)
//
//		return
//	}
//
//	var (
//		cancel1, cancel2, cancel3 func()
//	)
//
//	mtx := &sync.Mutex{}
//	canceled := false
//	cancelAll := func() bool {
//		mtx.Lock()
//		defer mtx.Unlock()
//
//		if !canceled {
//			cancel1()
//			cancel2()
//			cancel3()
//
//			canceled = true
//
//			return !canceled
//		}
//
//		return canceled
//	}
//
//	mtx.Lock()
//	defer mtx.Unlock()
//
//	votes := 0
//
//	cancel1 = s.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
//		if r.MessageID == msg.ID && r.MessageReaction.Emoji.Name == "😐" {
//			votes++
//
//			if ctx.Author.ID != r.UserID && votes != 3 && !hasRoleName(s, ctx.GuildID, r.UserID, "Gopher Herder") {
//				return
//			}
//
//			if cancelAll() {
//				return
//			}
//
//			time.Sleep(3 * time.Second)
//			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
//		}
//	})
//
//	cancel2 = s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
//		if m.Message.ID == ctx.ID {
//			if cancelAll() {
//				return
//			}
//
//			s.ChannelMessageDelete(msg.ChannelID, msg.ID)
//		}
//	})
//
//	cancel3 = s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
//		if m.Message.ID == ctx.ID {
//			if cancelAll() {
//				return
//			}
//
//			s.ChannelMessageDelete(msg.ChannelID, msg.ID)
//		}
//	})
//
//	time.AfterFunc(delay, func() {
//		cancelAll()
//	})
//}

func hasRoleName(s *discordgo.Session, guildID, userID, roleName string) bool {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		log.Println("hasRoleName(session.member):", err)

		return false
	}

	for _, rid := range member.Roles {
		role, err := s.State.Role(guildID, rid)
		if err != nil {
			log.Println("hasRoleName(state.role):", err)

			return false
		}

		if role.Name == roleName {
			return true
		}
	}

	return false
}

func isPrintable(content string) bool {
	for _, r := range content {
		if unicode.IsPrint(r) {
			return true
		}
	}

	return false
}

func commands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{Name: "go", Description: "Run Go code", Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionBoolean, Name: "public", Description: "Display output publicly"},
		}},
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
	dg.AddHandler(func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("command", "name", i.ApplicationCommandData().Name, "user", i.Member.User)

		err := dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Not implemented.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error("respond to interaction", "err", err)
			return
		}
	})

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
