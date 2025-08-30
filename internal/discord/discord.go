package discord

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	session *discordgo.Session
	token   string
	dmID    string
	enabled bool
}

func New() *Client {
	token := os.Getenv("DISCORD_USER_TOKEN")
	dmID := os.Getenv("DISCORD_DM_ID")
	
	client := &Client{
		token:   token,
		dmID:    dmID,
		enabled: token != "" && dmID != "",
	}
	
	return client
}

func (c *Client) IsEnabled() bool {
	return c.enabled
}

func (c *Client) GetStatus() (string, string) {
	if !c.enabled {
		if c.token == "" && c.dmID == "" {
			return "disabled", "DISCORD_USER_TOKEN and DISCORD_DM_ID not set"
		} else if c.token == "" {
			return "disabled", "DISCORD_USER_TOKEN not set"
		} else if c.dmID == "" {
			return "disabled", "DISCORD_DM_ID not set"
		}
	}
	return "enabled", fmt.Sprintf("Ready to send to %s", c.dmID[:8]+"...")
}

// SetDMID allows setting the Discord DM ID from within the program
func (c *Client) SetDMID(dmID string) bool {
	if dmID == "" {
		return false
	}
	c.dmID = dmID
	c.enabled = c.token != "" && c.dmID != ""
	return true
}

// GetDMID returns the current Discord DM ID
func (c *Client) GetDMID() string {
	return c.dmID
}

func (c *Client) Connect() error {
	if !c.enabled {
		return fmt.Errorf("discord not configured")
	}
	
	var err error
	c.session, err = discordgo.New(c.token) // User token, no "Bot " prefix
	if err != nil {
		return fmt.Errorf("failed to create Discord session: %w", err)
	}
	
	// Open websocket connection to Discord
	err = c.session.Open()
	if err != nil {
		return fmt.Errorf("failed to connect to Discord: %w", err)
	}
	
	return nil
}

// FetchLastMessage retrieves the last message from the DM channel
func (c *Client) FetchLastMessage() (*discordgo.Message, error) {
	if !c.enabled {
		return nil, fmt.Errorf("discord not configured")
	}
	
	if c.session == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}
	
	// Fetch the last message (limit 1)
	messages, err := c.session.ChannelMessages(c.dmID, 1, "", "", "")
	if err != nil {
		if discordErr, ok := err.(*discordgo.RESTError); ok {
			switch discordErr.Message.Code {
			case 10003: // Unknown Channel
				return nil, fmt.Errorf("invalid DISCORD_DM_ID - channel not found")
			case 50001: // Missing Access
				return nil, fmt.Errorf("cannot access DM channel - check permissions")
			default:
				return nil, fmt.Errorf("Discord API error %d: %s", discordErr.Message.Code, discordErr.Message.Message)
			}
		}
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}
	
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages found in DM")
	}
	
	return messages[0], nil
}

// ValidateText2BabeMessage checks if a message is from text2babe and extracts the encrypted data
func (c *Client) ValidateText2BabeMessage(message *discordgo.Message) (string, string, error) {
	content := message.Content
	
	// Check if the message matches text2babe format: 🔒/**🔓 **Text2Babe encrypt/decrypt**
	text2babePattern := regexp.MustCompile(`^(🔒|🔓)\s\*\*Text2Babe\s(encrypt|decrypt)\*\*\s*\n\x60\x60\x60\s*(.*?)\s*\n\x60\x60\x60$`)
	
	matches := text2babePattern.FindStringSubmatch(strings.TrimSpace(content))
	if len(matches) != 4 {
		return "", "", fmt.Errorf("message is not from text2babe (invalid format)")
	}
	
	_ = matches[1] // emoji (not used in current implementation)
	mode := matches[2]
	data := strings.TrimSpace(matches[3])
	
	if data == "" {
		return "", "", fmt.Errorf("no encrypted data found in message")
	}
	
	return data, mode, nil
}

// FetchAndValidateLastMessage fetches the last message and validates it's from text2babe
func (c *Client) FetchAndValidateLastMessage() (string, string, error) {
	message, err := c.FetchLastMessage()
	if err != nil {
		return "", "", err
	}
	
	// Check if it's from the current user (our messages)
	if c.session != nil && c.session.State.User != nil {
		if message.Author.ID != c.session.State.User.ID {
			return "", "", fmt.Errorf("last message is not from your account")
		}
	}
	
	return c.ValidateText2BabeMessage(message)
}

func (c *Client) SendMessage(content string) error {
	if !c.enabled {
		return fmt.Errorf("discord not configured")
	}
	
	if c.session == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}
	
	_, err := c.session.ChannelMessageSend(c.dmID, content)
	if err != nil {
		// Provide more helpful error messages for common issues
		if discordErr, ok := err.(*discordgo.RESTError); ok {
			switch discordErr.Message.Code {
			case 50001: // Missing Access
				return fmt.Errorf("bot cannot send DMs to this user - ensure you share a server and user allows DMs from server members")
			case 50007: // Cannot send messages to this user
				return fmt.Errorf("user has DMs disabled or bot is blocked")
			case 10003: // Unknown Channel
				return fmt.Errorf("invalid channel/user ID - check DISCORD_DM_ID")
			default:
				return fmt.Errorf("Discord API error %d: %s", discordErr.Message.Code, discordErr.Message.Message)
			}
		}
		return fmt.Errorf("failed to send Discord message: %w", err)
	}
	
	return nil
}

func (c *Client) SendEncryptedData(data, mode string) error {
	if !c.enabled {
		return fmt.Errorf("discord not configured")
	}
	
	emoji := "🔒"
	if mode == "decrypt" {
		emoji = "🔓"
	}
	
	message := fmt.Sprintf("%s **Text2Babe %s**\n```\n%s\n```", emoji, mode, data)
	
	return c.SendMessage(message)
}

func (c *Client) Disconnect() {
	if c.session != nil {
		c.session.Close()
		c.session = nil
	}
}

// TestConnection attempts to validate the bot setup
func (c *Client) TestConnection() error {
	if !c.enabled {
		return fmt.Errorf("discord not configured")
	}
	
	if c.session == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}
	
	// Try to get channel info to validate the DM channel
	_, err := c.session.Channel(c.dmID)
	if err != nil {
		if discordErr, ok := err.(*discordgo.RESTError); ok {
			switch discordErr.Message.Code {
			case 10003: // Unknown Channel
				return fmt.Errorf("invalid DISCORD_DM_ID - channel/user not found")
			case 50001: // Missing Access
				return fmt.Errorf("bot cannot access this channel - check permissions")
			default:
				return fmt.Errorf("Discord API error %d: %s", discordErr.Message.Code, discordErr.Message.Message)
			}
		}
		return fmt.Errorf("failed to validate Discord channel: %w", err)
	}
	
	return nil
}