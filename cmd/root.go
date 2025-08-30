package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"doc0x1/text2babe/internal/config"
	"doc0x1/text2babe/internal/crypto"
	"doc0x1/text2babe/internal/style"
	"doc0x1/text2babe/pkg/prompt"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "text2babe",
	Short: "🔐 Text2Babe - Encryption/Decryption CLI Tool",
	Long: `Text2Babe is a CLI tool for encrypting and decrypting text and binary data.
It features an interactive shell with zsh-style prompts and easy toggling between modes.`,
	Run: runInteractiveShell,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cfg = config.New()

	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(shellCmd)
}

func runInteractiveShell(cmd *cobra.Command, args []string) {
	fmt.Println(style.Banner())
	fmt.Println()
	fmt.Printf("%s\n", style.Gray.Sprint("Type 'help' for available commands or 'exit' to quit."))
	fmt.Println()

	p, err := prompt.New()
	if err != nil {
		fmt.Printf("Error creating prompt: %v\n", err)
		return
	}
	defer p.Close()

	for {
		line, err := p.ReadLine()
		if err == io.EOF {
			fmt.Println(style.Info.Sprint("\nGoodbye! 👋"))
			break
		}
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		if line == "" {
			continue
		}

		handleInteractiveCommand(line, p)
	}
}

func handleInteractiveCommand(input string, p *prompt.Prompt) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])

	switch command {
	case "help", "h":
		showHelp()
	case "settings", "config":
		showSettings()
	case "mode", "m":
		if len(parts) >= 2 {
			handleModeCommand(parts[1], p)
		} else {
			fmt.Printf("%s\n", style.Info.Sprintf("Current mode: %s", cfg.Mode))
		}
	case "set":
		if len(parts) >= 3 {
			handleSet(parts[1], parts[2], p)
		} else {
			fmt.Println("Usage: set <setting> <value>")
		}
	case "encrypt", "e":
		if len(parts) >= 2 {
			data := strings.Join(parts[1:], " ")
			result, err := crypto.EncryptData(data, cfg)
			if err != nil {
				fmt.Println(style.ErrorMsg(err))
			} else {
				fmt.Println(style.Result("Encrypted", result))

				// Copy to clipboard
				if err := clipboard.WriteAll(result); err != nil {
					fmt.Println(style.WarningMsg("Failed to copy to clipboard: " + err.Error()))
				} else {
					fmt.Println(style.SuccessWithClipboard("Encrypted"))
				}

				// Send to Discord if enabled
				discord := cfg.GetDiscord()
				if cfg.SendToDiscord && discord.IsEnabled() {
					if err := discord.SendEncryptedData(result, cfg.Mode); err != nil {
						fmt.Println(style.WarningMsg("Failed to send to Discord: " + err.Error()))
					} else {
						fmt.Println(style.Success.Sprint("📨 Sent to Discord!"))
					}
				}
			}
		} else {
			fmt.Println("Usage: encrypt <data>")
		}
	case "decrypt", "d":
		if len(parts) >= 2 {
			data := strings.Join(parts[1:], " ")
			result, err := crypto.DecryptData(data, cfg)
			if err != nil {
				fmt.Println(style.ErrorMsg(err))
			} else {
				fmt.Println(style.Result("Decrypted", result))
				if err := clipboard.WriteAll(result); err != nil {
					fmt.Println(style.WarningMsg("Failed to copy to clipboard: " + err.Error()))
				} else {
					fmt.Println(style.SuccessWithClipboard("Decrypted"))
				}
			}
		} else {
			fmt.Println("Usage: decrypt <data>")
		}
	case "toggle", "t":
		if len(parts) >= 2 {
			handleToggle(parts[1], p)
		} else {
			fmt.Println("Usage: toggle <setting>")
		}
	case "key":
		if len(parts) >= 2 {
			password := strings.Join(parts[1:], " ")
			cfg.SetKey(password)
			fmt.Printf("%s\n", style.Success.Sprint("✓ Encryption key updated"))
			fmt.Printf("%s %s\n", style.Info.Sprint("Key fingerprint:"), cfg.GetKeyFingerprint())
			if len(password) < 8 {
				fmt.Printf("%s\n", style.Warning.Sprint("⚠ Consider using a longer password for better security"))
			}
		} else {
			fmt.Println("Usage: key <password>")
		}
	case "discord":
		if len(parts) >= 2 {
			subCommand := strings.ToLower(parts[1])
			switch subCommand {
			case "test":
				discord := cfg.GetDiscord()
				if !discord.IsEnabled() {
					fmt.Println(style.ErrorMsg(fmt.Errorf("discord not configured (missing DISCORD_USER_TOKEN or DISCORD_DM_ID)")))
				} else {
					fmt.Println(style.Info.Sprint("Testing Discord connection..."))
					if err := discord.TestConnection(); err != nil {
						fmt.Println(style.ErrorMsg(err))
					} else {
						fmt.Println(style.Success.Sprint("✓ Discord connection successful!"))
					}
				}
			case "fetch", "decrypt":
				discord := cfg.GetDiscord()
				if !discord.IsEnabled() {
					fmt.Println(style.ErrorMsg(fmt.Errorf("discord not configured (missing DISCORD_USER_TOKEN or DISCORD_DM_ID)")))
					return
				}

				fmt.Println(style.Info.Sprint("Fetching last Discord message..."))
				data, mode, err := discord.FetchAndValidateLastMessage()
				if err != nil {
					fmt.Println(style.ErrorMsg(err))
					return
				}

				fmt.Printf("%s Found text2babe message (%s mode)\n", style.Success.Sprint("✓"), mode)

				// Decode/decrypt the message using auto-detection
				result, decryptErr := crypto.DecryptData(data, cfg)
				
				if decryptErr != nil {
					// If that fails, try without encryption (plain encoding)
					oldEncryption := cfg.UseEncryption
					cfg.UseEncryption = false
					result, decryptErr = crypto.DecryptData(data, cfg)
					cfg.UseEncryption = oldEncryption
				}

				if decryptErr != nil {
					fmt.Println(style.ErrorMsg(decryptErr))
				} else {
					fmt.Println(style.Result("Decoded from Discord", result))
				}
			default:
				discord := cfg.GetDiscord()
				status, message := discord.GetStatus()
				fmt.Printf("%s: %s\n", style.Info.Sprint("Discord Status"), status)
				fmt.Printf("%s: %s\n", style.Info.Sprint("Details"), message)
			}
		} else {
			discord := cfg.GetDiscord()
			status, message := discord.GetStatus()
			fmt.Printf("%s: %s\n", style.Info.Sprint("Discord Status"), status)
			fmt.Printf("%s: %s\n", style.Info.Sprint("Details"), message)
		}
	case "exit", "quit", "q":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
	}
}

func showHelp() {
	fmt.Println(style.Section("📖 Available Commands:"))
	fmt.Println(style.Command("help, h", "Show this help message"))
	fmt.Println(style.Command("settings, config", "Show current settings"))
	fmt.Println(style.Command("mode, m [encrypt/e/decrypt/d]", "Set or show current mode"))
	fmt.Println(style.Command("encrypt, e <data>", "Encrypt data"))
	fmt.Println(style.Command("decrypt, d <data>", "Decrypt data"))
	fmt.Println(style.Command("key <password>", "Set encryption key from password"))
	fmt.Println(style.Command("discord [test/fetch/decrypt]", "Show Discord status, test connection, or fetch+decrypt last message"))
	fmt.Println(style.Command("set <setting> <val>", "Set a configuration value"))
	fmt.Println(style.Command("toggle, t <setting>", "Toggle a setting"))
	fmt.Println(style.Command("exit, quit, q", "Exit the program"))

	fmt.Println(style.Section("⚙️  Settings:"))
	fmt.Println(style.Setting("mode", "encrypt/decrypt (shown by lock emoji in prompt)"))
	fmt.Println(style.Setting("output", "hex/base64/binary (encrypted data format, default: hex)"))
	fmt.Println(style.Setting("discord", "true/false (auto-send encrypted data to Discord DM)"))
	fmt.Println(style.Setting("discord-id", "channel_id (set Discord DM channel ID)"))

	fmt.Println(style.Section("🔄 How It Works:"))
	fmt.Printf("  %s\n", style.Info.Sprint("ENCRYPT: text input → AES-GCM → hex/base64/binary output → clipboard + Discord"))
	fmt.Printf("  %s\n", style.Info.Sprint("DECRYPT: hex/base64/binary input → AES-GCM → text output → clipboard"))

	fmt.Println(style.Section("💡 Examples:"))
	fmt.Println(style.Example("mode encrypt", "switch to encrypt mode (🔒)"))
	fmt.Println(style.Example("mode d", "switch to decrypt mode (🔓)"))
	fmt.Println(style.Section("📨 Discord Integration:"))
	fmt.Printf("  %s\n", style.Info.Sprint("Requires DISCORD_BOT_TOKEN and DISCORD_DM_ID in .env file"))
	fmt.Printf("  %s\n", style.Info.Sprint("Automatically sends encrypted data to specified Discord DM"))

	fmt.Println(style.Section("💡 Examples:"))
	fmt.Println(style.Example("encrypt hello world", "encrypt text + send to Discord"))
	fmt.Println(style.Example("decrypt a1b2c3d4...", "decrypt any format to text"))
	fmt.Println(style.Example("set output base64", "use base64 encoding"))
	fmt.Println(style.Example("set discord on", "enable Discord sending"))
	fmt.Println(style.Example("set discord-id 123456789", "set Discord DM channel ID"))
	fmt.Println(style.Example("discord fetch", "fetch and decrypt last Discord message"))
	fmt.Println(style.Example("toggle discord", "toggle Discord on/off"))
	fmt.Println(style.Example("key secretpassword", "set encryption key"))
	fmt.Println()
}

func showSettings() {
	fmt.Println(style.Section("Current Configuration:"))

	// Mode with visual indicator and encryption type
	var modeDisplay string
	var emoji string

	if cfg.Mode == "encrypt" {
		emoji = "🔒"
	} else {
		emoji = "🔓"
	}

	encType := "AES-GCM"
	if !cfg.UseEncryption {
		encType = "plain"
	}

	modeDisplay = fmt.Sprintf("%s %s (%s)", cfg.Mode, emoji, encType)
	fmt.Println(style.Setting("Mode", modeDisplay))
	fmt.Println(style.Setting("Output Format", cfg.OutputMode))

	// Key information
	keyInfo := cfg.GetKeyFingerprint()
	if cfg.IsDefaultKey() {
		keyInfo = keyInfo + " " + style.Warning.Sprint("(default - consider changing!)")
	} else {
		keyInfo = keyInfo + " " + style.Success.Sprint("(custom)")
	}
	fmt.Println(style.Setting("Key Fingerprint", keyInfo))

	if cfg.UseEncryption {
		fmt.Println(style.Setting("Encryption", "Enabled - AES-256-GCM"))
		fmt.Println(style.Setting("Key Derivation", "SHA-256"))
	} else {
		fmt.Println(style.Setting("Encryption", "Off"))
	}

	// Discord integration
	discord := cfg.GetDiscord()
	discordStatus, discordMessage := discord.GetStatus()
	discordDisplay := discordStatus
	if discordStatus == "enabled" {
		if cfg.SendToDiscord {
			discordDisplay = discordDisplay + " ✅ (active)"
		} else {
			discordDisplay = discordDisplay + " ⏸️ (paused)"
		}
		discordDisplay = discordDisplay + " - " + discordMessage
	} else {
		discordDisplay = discordDisplay + " - " + discordMessage
	}
	fmt.Println(style.Setting("Discord", discordDisplay))

	// Show Discord channel ID if available
	if dmID := discord.GetDMID(); dmID != "" {
		fmt.Println(style.Setting("Discord DM ID", dmID))
	} else {
		fmt.Println(style.Setting("Discord DM ID", style.Warning.Sprint("not set")))
	}

	fmt.Println()
}

func handleSet(setting, value string, p *prompt.Prompt) {
	switch strings.ToLower(setting) {
	case "mode":
		if cfg.SetMode(value) {
			fmt.Printf("%s\n", style.Success.Sprintf("Mode set to: %s", value))
			p.UpdatePrompt(cfg.Mode) // Update prompt with new mode
		} else {
			fmt.Println(style.ErrorMsg(fmt.Errorf("mode must be 'encrypt' or 'decrypt'")))
		}
	case "output", "format":
		if cfg.SetOutputMode(value) {
			fmt.Printf("%s\n", style.Success.Sprintf("Output mode set to: %s", value))
		} else {
			fmt.Println(style.ErrorMsg(fmt.Errorf("output mode must be 'hex', 'base64', or 'binary'")))
		}
	case "discord":
		switch value {
		case "true", "on", "enable":
			if cfg.SetDiscordSending(true) {
				fmt.Printf("%s\n", style.Success.Sprintf("Discord sending enabled"))
			} else {
				fmt.Println(style.ErrorMsg(fmt.Errorf("discord not configured (missing DISCORD_USER_TOKEN or DISCORD_DM_ID)")))
			}
		case "false", "off", "disable":
			cfg.SetDiscordSending(false)
			fmt.Printf("%s\n", style.Success.Sprintf("Discord sending disabled"))
		default:
			fmt.Println(style.ErrorMsg(fmt.Errorf("discord must be 'true/on/enable' or 'false/off/disable'")))
		}
	case "encryption":
		switch value {
		case "true", "on", "enable":
			cfg.SetEncryption(true)
			fmt.Printf("%s\n", style.Success.Sprintf("Encryption enabled - using AES-GCM"))
			p.UpdatePrompt(cfg.Mode)
		case "false", "off", "disable":
			cfg.SetEncryption(false)
			fmt.Printf("%s\n", style.Success.Sprintf("Encryption disabled - using plain encoding"))
			p.UpdatePrompt(cfg.Mode)
		default:
			fmt.Println(style.ErrorMsg(fmt.Errorf("encryption must be 'true/on/enable' or 'false/off/disable'")))
		}
	case "discord-id", "dmid":
		discord := cfg.GetDiscord()
		if discord.SetDMID(value) {
			fmt.Printf("%s\n", style.Success.Sprintf("Discord DM ID set to: %s", value))
		} else {
			fmt.Println(style.ErrorMsg(fmt.Errorf("invalid Discord DM ID (cannot be empty)")))
		}
	default:
		fmt.Printf("%s\n", style.ErrorMsg(fmt.Errorf("unknown setting: %s", setting)))
	}
}

func handleToggle(setting string, p *prompt.Prompt) {
	switch strings.ToLower(setting) {
	case "mode":
		cfg.ToggleMode()
		fmt.Printf("%s\n", style.Success.Sprintf("Mode toggled to: %s", cfg.Mode))
		p.UpdatePrompt(cfg.Mode) // Update prompt with new mode
	case "output", "format":
		cfg.ToggleOutputMode()
		fmt.Printf("%s\n", style.Success.Sprintf("Output mode toggled to: %s", cfg.OutputMode))
	case "discord":
		if cfg.ToggleDiscord() {
			status := "disabled"
			if cfg.SendToDiscord {
				status = "enabled"
			}
			fmt.Printf("%s\n", style.Success.Sprintf("Discord sending toggled to: %s", status))
		} else {
			fmt.Println(style.ErrorMsg(fmt.Errorf("discord not configured (missing DISCORD_USER_TOKEN or DISCORD_DM_ID)")))
		}
	case "encryption":
		cfg.ToggleEncryption()
		status := "disabled (plain encoding)"
		if cfg.UseEncryption {
			status = "enabled (AES-GCM)"
		}
		fmt.Printf("%s\n", style.Success.Sprintf("Encryption toggled to: %s", status))
		p.UpdatePrompt(cfg.Mode)
	default:
		fmt.Printf("%s\n", style.ErrorMsg(fmt.Errorf("cannot toggle setting: %s", setting)))
	}
}

func handleModeCommand(mode string, p *prompt.Prompt) {
	switch strings.ToLower(mode) {
	case "encrypt", "e":
		cfg.SetMode("encrypt")
		encType := "AES-GCM"
		if !cfg.UseEncryption {
			encType = "plain"
		}
		fmt.Printf("%s\n", style.Success.Sprintf("Mode set to: encrypt (%s)", encType))
		p.UpdatePrompt("encrypt")
	case "decrypt", "d":
		cfg.SetMode("decrypt")
		encType := "AES-GCM"
		if !cfg.UseEncryption {
			encType = "plain"
		}
		fmt.Printf("%s\n", style.Success.Sprintf("Mode set to: decrypt (%s)", encType))
		p.UpdatePrompt("decrypt")
	case "toggle", "t":
		cfg.ToggleMode()
		encType := "AES-GCM"
		if !cfg.UseEncryption {
			encType = "plain"
		}
		fmt.Printf("%s\n", style.Success.Sprintf("Mode toggled to: %s (%s)", cfg.Mode, encType))
		p.UpdatePrompt(cfg.Mode)
	default:
		fmt.Printf("%s\n", style.ErrorMsg(fmt.Errorf("invalid mode. Use: encrypt/e, decrypt/d, or toggle/t")))
	}
}
