package prompt

import (
	"io"
	"strings"

	"github.com/chzyer/readline"

	"doc0x1/text2babe/internal/style"
)

type Prompt struct {
	rl *readline.Instance
}

func New() (*Prompt, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          style.PromptWithMode("encrypt"), // Start with encrypt mode
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return nil, err
	}

	return &Prompt{rl: rl}, nil
}

// UpdatePrompt changes the prompt to reflect current mode
func (p *Prompt) UpdatePrompt(mode string) {
	p.rl.SetPrompt(style.PromptWithMode(mode))
}

func (p *Prompt) ReadLine() (string, error) {
	line, err := p.rl.Readline()
	switch err {
	case readline.ErrInterrupt:
		return "", io.EOF
	case io.EOF:
		return "exit", nil
	}
	return strings.TrimSpace(line), err
}

func (p *Prompt) Close() error {
	return p.rl.Close()
}

func (p *Prompt) SetPrompt(prompt string) {
	p.rl.SetPrompt(prompt)
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("help"),
	readline.PcItem("settings"),
	readline.PcItem("config"),
	readline.PcItem("mode",
		readline.PcItem("encrypt"),
		readline.PcItem("decrypt"),
		readline.PcItem("toggle"),
		readline.PcItem("e"),
		readline.PcItem("d"),
		readline.PcItem("t"),
	),
	readline.PcItem("set",
		readline.PcItem("mode",
			readline.PcItem("encrypt"),
			readline.PcItem("decrypt"),
		),
		readline.PcItem("output",
			readline.PcItem("hex"),
			readline.PcItem("base64"),
			readline.PcItem("binary"),
		),
		readline.PcItem("discord",
			readline.PcItem("on"),
			readline.PcItem("off"),
			readline.PcItem("enable"),
			readline.PcItem("disable"),
			readline.PcItem("true"),
			readline.PcItem("false"),
		),
		readline.PcItem("encryption",
			readline.PcItem("on"),
			readline.PcItem("off"),
			readline.PcItem("enable"),
			readline.PcItem("disable"),
			readline.PcItem("true"),
			readline.PcItem("false"),
		),
		readline.PcItem("discord-id"),
		readline.PcItem("dmid"),
	),
	readline.PcItem("toggle",
		readline.PcItem("mode"),
		readline.PcItem("output"),
		readline.PcItem("discord"),
		readline.PcItem("encryption"),
	),
	readline.PcItem("encrypt"),
	readline.PcItem("decrypt"),
	readline.PcItem("key"),
	readline.PcItem("discord",
		readline.PcItem("test"),
		readline.PcItem("fetch"),
		readline.PcItem("decrypt"),
	),
	readline.PcItem("exit"),
	readline.PcItem("quit"),
)
