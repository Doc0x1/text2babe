package cmd

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"doc0x1/text2babe/internal/crypto"
	"doc0x1/text2babe/internal/style"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [data]",
	Short: "Decrypt data using current settings",
	Long:  "Decrypt AES-GCM encrypted data using the current configuration.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data := strings.Join(args, " ")
		result, err := crypto.DecryptData(data, cfg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println(style.Result("Decrypted", result))
		if err := clipboard.WriteAll(result); err != nil {
			fmt.Println(style.WarningMsg("Failed to copy to clipboard: " + err.Error()))
		} else {
			fmt.Println(style.SuccessWithClipboard("Decrypted"))
		}
	},
}