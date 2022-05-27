package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mjehanno/new/internal/pipe"
)

func RunTemplate(url string) tea.Cmd {
	pipe.Chan <- url
	return nil
}
