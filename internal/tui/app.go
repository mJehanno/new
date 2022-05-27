package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/knipferrc/teacup/statusbar"
	command "github.com/mjehanno/new/internal/cmd"
	"github.com/mjehanno/new/internal/message"
	"github.com/mjehanno/new/internal/pipe"
)

type languageColor string

const (
	javascript = "#e3f542"
	golang     = "#42ddf5"
	java       = "#c41d1d"
	php        = "#055c85"
	csharp     = "#5f0469"
	cpp        = "#043669"
	python     = "#039603"
	swift      = "#bf3c08"
	rust       = "#fa5b1b"
	typescript = "#2a379c"
	ruby       = "#e60505"
)

type appModel struct {
	ctx       context.Context
	table     table.Model
	statusBar statusbar.Bubble
	rows      []table.Row
	page      int
}

const (
	columnKeyFullname = "fullname"
	columnKeyDesc     = "description"
	columnKeyLang     = "language"
	columnKeyStars    = "stars"
	columnHTMLURL     = "HTMLURL"
)

func InitialModel() *appModel {
	ctx := context.Background()
	status := statusbar.New(
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#A550DF", Dark: "#A550DF"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#6124DF", Dark: "#6124DF"},
		},
	)
	return &appModel{
		ctx:       ctx,
		statusBar: status,
		page:      1,
		table: table.New(
			[]table.Column{
				table.NewColumn(columnKeyFullname, "Fullname", 40).WithFiltered(true),
				table.NewColumn(columnKeyDesc, "Description", 100).WithFiltered(true),
				table.NewColumn(columnKeyLang, "Language", 10).WithFiltered(true),
				table.NewColumn(columnKeyStars, "Stars", 5),
				table.NewColumn(columnHTMLURL, "", 0),
			},
		).Filtered(true).Focused(true).WithPageSize(command.DisplaySize),
	}
}

func (m appModel) Init() tea.Cmd {
	return command.InitSearch(m.ctx)
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case message.ErrorMessage:
		m.statusBar.SetContent("", msg.Error.Error(), "", "")
	case message.SearchResultMessage:
		for _, result := range msg.SearchResults {
			color := ""
			switch result.Language {
			case "Python":
				color = python
			case "JavaScript":
				color = javascript
			case "TypeScript":
				color = typescript
			case "C#":
				color = csharp
			case "Swift":
				color = swift
			case "Java":
				color = java
			case "Ruby":
				color = ruby
			case "C++":
				color = cpp
			case "Go":
				color = golang
			case "Rust":
				color = rust
			case "PHP":
				color = php
			}

			m.rows = append(m.rows, table.NewRow(
				table.RowData{
					columnKeyFullname: result.Fullname,
					columnKeyDesc:     result.Description,
					columnKeyLang:     table.NewStyledCell(result.Language, lipgloss.NewStyle().Foreground(lipgloss.Color(color))),
					columnKeyStars:    result.Stars,
					columnHTMLURL:     result.HTMLURL,
				},
			))
		}

		m.table = m.table.WithRows(m.rows).WithStaticFooter(fmt.Sprintf("Page %v / %v", m.table.CurrentPage(), command.TotalPageNumber))
		m.table, cmd = m.table.Update(msg)
		return tea.Model(m), cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return tea.Model(m), tea.Batch(tea.Quit, func() tea.Msg {
				pipe.Chan <- ""
				return nil
			})
		case "enter":
			HTMLURL := m.table.HighlightedRow().Data[columnHTMLURL]
			s := HTMLURL.(string)
			return tea.Model(m), tea.Batch(tea.Quit, command.RunTemplate(s))
		}
	case tea.WindowSizeMsg:
		m.statusBar.SetSize(msg.Width)
	}

	m.table, cmd = m.table.Update(msg)

	return tea.Model(m), cmd
}

func (m appModel) View() string {
	s := ""
	s += m.table.View()
	s += "\n \n"
	s += m.statusBar.View()
	s += "\n"
	return s
}
