package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

type appModel struct {
	table table.Model
}

func InitialModel() *appModel {
	return &appModel{}
}

func (m appModel) Init() tea.Cmd {
	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return tea.Model(m), nil
}

func (m appModel) View() string {
	return "hello from TUI"
}
