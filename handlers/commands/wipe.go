package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Wipe deletes a specified number of messages within the channel
func Wipe(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// Ensure the command is not executed by a bot
	if m.Author.Bot {
		return nil
	}

	// Get the number of messages to delete
	parts := strings.Split(m.Content, " ")
	if len(parts) < 2 {
		log.Println("usage: !wipe <number_of_messages>")
		return fmt.Errorf("usage: !wipe <number_of_messages>")
	}

	numMessages := parts[1]
	num, err := strconv.Atoi(numMessages)
	if err != nil {
		log.Printf("invalid number of messages: %v", err)
		return fmt.Errorf("invalid number of messages: %v", err)
	}

	// Ensure the number of messages is within a reasonable limit
	if num <= 0 || num > 100 {
		log.Println("number of messages must be between 1 and 100")
		return fmt.Errorf("number of messages must be between 1 and 100")
	}

	// Fetch the messages from the channel
	messages, err := s.ChannelMessages(m.ChannelID, num, "", "", "")
	if err != nil {
		log.Printf("failed to fetch messages: %v", err)
		return fmt.Errorf("failed to fetch messages: %v", err)
	}

	// Collect message IDs to delete
	var messageIDs []string
	for _, msg := range messages {
		messageIDs = append(messageIDs, msg.ID)
	}

	// Bulk delete the messages
	if len(messageIDs) > 0 {
		err = s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
		if err != nil {
			log.Printf("failed to delete messages: %v", err)
			return fmt.Errorf("failed to delete messages: %v", err)
		}
	}

	// Send a confirmation message
	confirmation, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Deleted %d messages.", len(messageIDs)))
	if err != nil {
		log.Printf("failed to send confirmation message: %v", err)
		return fmt.Errorf("failed to send confirmation message: %v", err)
	}

	// Delete the confirmation message after a few seconds
	time.AfterFunc(5*time.Second, func() {
		s.ChannelMessageDelete(m.ChannelID, confirmation.ID)
	})

	return nil
}
