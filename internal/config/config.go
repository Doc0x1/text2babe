package config

import (
	"crypto/sha256"
	"doc0x1/text2babe/internal/discord"
	"encoding/hex"
)

type Config struct {
	Mode          string
	DataType      string
	OutputMode    string
	Key           []byte
	KeySource     string // Track what password/source was used
	Discord       *discord.Client
	SendToDiscord bool // Toggle for Discord sending
	UseEncryption bool // Toggle for AES encryption vs plain encoding
}

func New() *Config {
	return &Config{
		Mode:          "encrypt",
		DataType:      "text",
		OutputMode:    "hex",
		Key:           generateKey("default-password"),
		KeySource:     "default-password",
		Discord:       nil, // Initialize lazily
		SendToDiscord: true,
		UseEncryption: false, // Default to encryption disabled
	}
}

func (c *Config) SetMode(mode string) bool {
	if mode == "encrypt" || mode == "decrypt" {
		c.Mode = mode
		return true
	}
	return false
}

func (c *Config) SetDataType(dataType string) bool {
	if dataType == "text" || dataType == "binary" {
		c.DataType = dataType
		return true
	}
	return false
}

func (c *Config) SetOutputMode(outputMode string) bool {
	if outputMode == "hex" || outputMode == "base64" || outputMode == "binary" {
		c.OutputMode = outputMode
		return true
	}
	return false
}

func (c *Config) ToggleMode() {
	if c.Mode == "encrypt" {
		c.Mode = "decrypt"
	} else {
		c.Mode = "encrypt"
	}
}

func (c *Config) ToggleDataType() {
	if c.DataType == "text" {
		c.DataType = "binary"
	} else {
		c.DataType = "text"
	}
}

func (c *Config) ToggleOutputMode() {
	switch c.OutputMode {
	case "hex":
		c.OutputMode = "base64"
	case "base64":
		c.OutputMode = "binary"
	case "binary":
		c.OutputMode = "hex"
	default:
		c.OutputMode = "hex"
	}
}

func (c *Config) SetKey(password string) {
	c.Key = generateKey(password)
	c.KeySource = password
}

// GetKeyFingerprint returns a short hex representation of the key for display
func (c *Config) GetKeyFingerprint() string {
	if len(c.Key) >= 4 {
		return hex.EncodeToString(c.Key[:4]) + "..."
	}
	return "unknown"
}

// IsDefaultKey returns true if using the default password
func (c *Config) IsDefaultKey() bool {
	return c.KeySource == "default-password"
}

// GetDiscord lazily initializes and returns the Discord client
func (c *Config) GetDiscord() *discord.Client {
	if c.Discord == nil {
		c.Discord = discord.New()
	}
	return c.Discord
}

// ToggleDiscord toggles Discord sending on/off
func (c *Config) ToggleDiscord() bool {
	discord := c.GetDiscord()
	if !discord.IsEnabled() {
		return false // Can't enable if not configured
	}
	c.SendToDiscord = !c.SendToDiscord
	return true
}

// SetDiscordSending sets Discord sending state
func (c *Config) SetDiscordSending(enabled bool) bool {
	discord := c.GetDiscord()
	if enabled && !discord.IsEnabled() {
		return false // Can't enable if not configured
	}
	c.SendToDiscord = enabled
	return true
}

// ToggleEncryption toggles between encrypted and plain encoding modes
func (c *Config) ToggleEncryption() {
	c.UseEncryption = !c.UseEncryption
	// Update mode to match encryption state
	if c.UseEncryption {
		switch c.Mode {
		case "encode":
			c.Mode = "encrypt"
		case "decode":
			c.Mode = "decrypt"
		}
	} else {
		switch c.Mode {
		case "encrypt":
			c.Mode = "encode"
		case "decrypt":
			c.Mode = "decode"
		}
	}
}

// SetEncryption sets encryption enabled/disabled
func (c *Config) SetEncryption(enabled bool) {
	c.UseEncryption = enabled
	// Update mode to match encryption state
	if c.UseEncryption {
		switch c.Mode {
		case "encode":
			c.Mode = "encrypt"
		case "decode":
			c.Mode = "decrypt"
		}
	} else {
		switch c.Mode {
		case "encrypt":
			c.Mode = "encode"
		case "decrypt":
			c.Mode = "decode"
		}
	}
}

func generateKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}
