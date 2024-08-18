package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"tina/structs"

	"github.com/bwmarrin/discordgo"
)

func NewMessage(s *discordgo.Session, m *discordgo.MessageCreate, state *structs.State) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	botMentioned := false
	mentionString := fmt.Sprintf("<@%s>", s.State.User.ID)

	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			botMentioned = true
			break
		}
	}

	if !botMentioned {
		return
	}

	messageContent := strings.ReplaceAll(m.Content, mentionString, "")
	messageContent = strings.TrimSpace(messageContent)

    resp, err := http.Post(fmt.Sprintf("http://api:6969"), "application/json", strings.NewReader(fmt.Sprintf("{\"query\": \"%s\"}", messageContent)))
	if err != nil {
		log.Printf("Failed to get response from api: %v", err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
	}

	text := strings.ReplaceAll(string(bytes), "\"", "")

	s.ChannelMessageSendReply(m.ChannelID, text, m.Reference())
	if err != nil {
		log.Printf("Failed to send reply: %v", err)
	}
}

func IntentAppend(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if m.MessageReference != nil {
        resp, err := http.Get("http://api:6969/get/intents")
		if err != nil {
			log.Printf("Cannot get intents: %v", err)
			return
		}

		reader, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Cannot read intents: %v", err)
			return
		}

		intents := structs.Intents{}

		err = json.Unmarshal(reader, &intents)
		if err != nil {
			log.Printf("Cannot unmarshal intents: %v", err)
			return
		}

		referenced, err := s.ChannelMessage(m.ChannelID, m.MessageReference.MessageID)
		if err != nil {
			return
		}

		if referenced.Author.ID == s.State.User.ID || referenced.Author.Bot {
			return
		}

		for _, intent := range intents.Intents {
			if intent.Tag == referenced.ID {
				// Check if the response already exists
				for _, response := range intent.Responses {
					if strings.Contains(m.Content, response) {
						return
					}
				}

				// Append the new response
                http.Post("http://api:6969/new/response", "application/json", strings.NewReader(fmt.Sprintf("{\"tag\": \"%s\", \"response\": \"%s\"}", intent.Tag, m.Content)))

				return
			}
		}

		// If the referenced message does not match any existing intent, create a new one
		intent := structs.Intent{
			Tag:       referenced.ID,
			Responses: []string{m.Content},
			Patterns:  []string{referenced.Content},
		}

		str, err := json.Marshal(intent)

        http.Post("http://api:6969/new/intent", "application/json", strings.NewReader(string(str)))

        time.Sleep(100 * time.Millisecond)

        http.Get("http://api:6969/retrain")
	}
}
