package cmd

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"doc0x1/text2babe/internal/crypto"
	"doc0x1/text2babe/internal/style"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [data]",
	Short: "Encrypt data using current settings",
	Long:  "Encrypt text or binary data using AES-GCM encryption with the current configuration.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data := strings.Join(args, " ")
		result, err := crypto.EncryptData(data, cfg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
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
	},
}