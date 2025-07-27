package services

import (
	"fmt"
	"notification-service/internal/config"
	"notification-service/internal/models"

	"github.com/slack-go/slack"
)

// SlackSender handles Slack notifications
type SlackSender struct {
	config *config.Config
	client *slack.Client
}

// NewSlackSender creates a new Slack sender
func NewSlackSender(config *config.Config) *SlackSender {
	return &SlackSender{
		config: config,
		client: slack.New(config.SlackToken),
	}
}

// Send sends a Slack notification
func (s *SlackSender) Send(notification *models.Notification) error {
	// Determine channel
	channel := s.config.SlackChannel
	if notification.Channel != "" {
		channel = notification.Channel
	}

	// Create message options
	options := []slack.MsgOption{
		slack.MsgOptionText(notification.Message, false),
	}

	// Add blocks if metadata contains them
	if blocks, exists := notification.Metadata["blocks"]; exists {
		if blockArray, ok := blocks.([]slack.Block); ok {
			options = append(options, slack.MsgOptionBlocks(blockArray...))
		}
	}

	// Add attachments if metadata contains them
	if attachments, exists := notification.Metadata["attachments"]; exists {
		if attachmentArray, ok := attachments.([]slack.Attachment); ok {
			options = append(options, slack.MsgOptionAttachments(attachmentArray...))
		}
	}

	// Send message
	_, _, err := s.client.PostMessage(channel, options...)
	if err != nil {
		return fmt.Errorf("failed to send Slack message: %w", err)
	}

	return nil
}

// TestConnection tests the Slack connection
func (s *SlackSender) TestConnection() error {
	// Test authentication
	_, err := s.client.AuthTest()
	if err != nil {
		return fmt.Errorf("failed to authenticate with Slack: %w", err)
	}

	return nil
}

// GetChannels retrieves available Slack channels
func (s *SlackSender) GetChannels() ([]slack.Channel, error) {
	channels, _, err := s.client.GetConversations(&slack.GetConversationsParameters{
		Types: []string{"public_channel", "private_channel"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get Slack channels: %w", err)
	}

	return channels, nil
} 