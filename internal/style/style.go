package style

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	// Colors
	Blue    = color.New(color.FgBlue, color.Bold)
	Green   = color.New(color.FgGreen, color.Bold)
	Red     = color.New(color.FgRed, color.Bold)
	Yellow  = color.New(color.FgYellow, color.Bold)
	Cyan    = color.New(color.FgCyan, color.Bold)
	Magenta = color.New(color.FgMagenta, color.Bold)
	White   = color.New(color.FgWhite, color.Bold)
	Gray    = color.New(color.FgHiBlack)

	// Special formatting
	Success = color.New(color.FgGreen, color.Bold)
	Error   = color.New(color.FgRed, color.Bold)
	Warning = color.New(color.FgYellow, color.Bold)
	Info    = color.New(color.FgCyan, color.Bold)
	Accent  = color.New(color.FgMagenta, color.Bold)
)

// Banner creates a styled banner with title
func Banner() string {
	border := strings.Repeat("═", 32)

	return fmt.Sprintf("%s\n%s %s %s\n%s",
		Cyan.Sprint(border),
		Accent.Sprint("🔐"),
		Blue.Sprint("Text2Babe Encryption Tool"),
		Accent.Sprint("🔐"),
		Cyan.Sprint(border))
}

// Prompt creates a styled prompt
func Prompt() string {
	return fmt.Sprintf("%s%s%s %s%s%s %s ",
		Gray.Sprint("["),
		Cyan.Sprint("🔐"),
		Gray.Sprint("]"),
		Green.Sprint("me"),
		Gray.Sprint("@"),
		Magenta.Sprint("babe"),
		Blue.Sprint("❯"))
}

// PromptWithMode creates a prompt that shows the current mode
func PromptWithMode(mode string) string {
	var lockEmoji string

	if mode == "decrypt" {
		lockEmoji = "🔓" // Unlock for decrypt
	} else {
		lockEmoji = "🔒" // Lock for encrypt
	}

	return fmt.Sprintf("%s%s%s %s%s%s %s ",
		Gray.Sprint("["),
		Cyan.Sprint(lockEmoji),
		Gray.Sprint("]"),
		Green.Sprint("me"),
		Gray.Sprint("@"),
		Magenta.Sprint("babe"),
		Blue.Sprint("❯"))
}

// Result formats a result message
func Result(label, value string) string {
	return fmt.Sprintf("%s %s",
		Info.Sprintf("%s:", label),
		White.Sprint(value))
}

// Success message with clipboard indicator
func SuccessWithClipboard(message string) string {
	return fmt.Sprintf("%s %s",
		Success.Sprint("✓"),
		Gray.Sprintf("%s 📋 Copied to clipboard!", message))
}

// Error message
func ErrorMsg(err error) string {
	return Error.Sprintf("✗ Error: %v", err)
}

// Warning message
func WarningMsg(message string) string {
	return Warning.Sprintf("⚠ Warning: %s", message)
}

// Setting display
func Setting(name, value string) string {
	return fmt.Sprintf("  %s %s",
		Cyan.Sprintf("%-12s", name+":"),
		White.Sprint(value))
}

// Command help
func Command(name, description string) string {
	return fmt.Sprintf("  %s %s",
		Green.Sprintf("%-20s", name),
		Gray.Sprint("- "+description))
}

// Section header
func Section(title string) string {
	return fmt.Sprintf("\n%s %s",
		Accent.Sprint("⚙️"),
		Blue.Sprint(title))
}

// Example
func Example(cmd, description string) string {
	return fmt.Sprintf("  %s %s",
		Yellow.Sprintf("%-25s", cmd),
		Gray.Sprint("# "+description))
}
