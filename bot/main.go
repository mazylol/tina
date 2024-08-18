package main

import (
	"log"
	"os"
	"os/signal"
	"tina/handlers"
	"tina/structs"

	"github.com/bwmarrin/discordgo"
)

var config = structs.Config{}

var state = structs.State{}

var s *discordgo.Session

func init() {
	config.Load()
	var err error
	s, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("Invalid bot token: %v", err)
	}

	s.Identify.Intents |= discordgo.IntentsAll

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handlers.IntentAppend(s, m)
        handlers.NewMessage(s, m, &state)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		err := s.UpdateGameStatus(0, config.Status)
		if err != nil {
			log.Fatalf("Cannot set status: %v", err)
		}
	})
}

func main() {
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			log.Fatalf("Cannot close the session: %v", err)
		}
	}(s)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}
