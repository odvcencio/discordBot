package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	run()
}

func run() {
	// generate token on discord and get the bot permission to post in the server
	token := os.Getenv("TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register the callback funcs
	dg.AddHandler(addPeach)
	dg.AddHandler(heathcliff)
	dg.AddHandler(reactWithUnicorn)
	dg.AddHandler(eightBall)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// adds a peach emoji reaction to any message that ends in "this"
func addPeach(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// normalize message contents to lower case
	contentNoCase := strings.ToLower(m.Content)

	// split the string and get the last word in it
	stringArray := strings.Split(contentNoCase, " ")
	lastPartOfString := stringArray[len(stringArray)-1]

	// if "this" is the last word (or anything including this), react with a peach emoji
	if strings.Contains(lastPartOfString, "this") {
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "🍑")
		if err != nil {
			fmt.Println("failed to react with peach emoji")
		}
	}
}

func reactWithUnicorn(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// normalize message contents to lower case
	contentNoCase := strings.ToLower(m.Content)

	// if "this" is the last word (or anything including this), react with a peach emoji
	if strings.Contains(contentNoCase, "unicorn") {
		// doesnt work in the unicorn servers for some reason! maybe has to do with guilds?

		// // get special custom unicorn emoji, not the ho-hum one
		// uniEmoji := getEmojiID(s, m, "unicorn-1")

		// fmt.Println(uniEmoji)
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "🦄")
		if err != nil {
			fmt.Println("failed to react with unicorn emoji")
		}
	}
}

// posts a link to a heathcliff comic that discord inlines automagically
func heathcliff(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// normalize message contents to lower case
	contentNoCase := strings.ToLower(m.Content)

	// randomly select a date between January 2003 - January 2020, post the url
	if strings.Contains(contentNoCase, "!heathcliff") {
		min := time.Date(2003, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
		max := time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
		delta := max - min

		// make the current unix timestamp in nanoseconds the seed for the rng
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		sec := r1.Int63n(delta) + min

		randoDate := time.Unix(sec, 0)

		heathURL := fmt.Sprintf("https://www.gocomics.com/heathcliff/%d/%d/%d",
			randoDate.Year(),
			randoDate.Month(),
			randoDate.Day())

		_, err := s.ChannelMessageSend(m.ChannelID, heathURL)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func eightBall(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	eightBallSayings := []string{
		"As I see it, yes.",
		"Ask again later.",
		"Better not tell you now.",
		"Cannot predict now.",
		"Concentrate and ask again.",
		"Don’t count on it.",
		"It is certain.",
		"It is decidedly so.",
		"Most likely.",
		"My reply is no.",
		"My sources say no.",
		"Outlook not so good.",
		"Outlook good.",
		"Reply hazy, try again.",
		"Signs point to yes.",
		"Very doubtful.",
		"Without a doubt.",
		"Yes.",
		"Yes – definitely.",
		"You may rely on it.",
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	contentNoCase := strings.ToLower(m.Content)

	if strings.Contains(contentNoCase, "!8ball") {
		answer := eightBallSayings[r1.Intn(len(eightBallSayings))]

		_, err := s.ChannelMessageSend(m.ChannelID, answer)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// UTILITY FUNCTIONS

// getEmojiID gets the safe-for-API-request version of custom emojis
func getEmojiID(s *discordgo.Session, m *discordgo.MessageCreate, name string) string {
	guild, _ := s.Guild(m.GuildID)
	emojis := guild.Emojis
	for _, emoji := range emojis {
		if strings.Contains(emoji.APIName(), name) {
			return emoji.APIName()
		}
	}

	return ""
}
