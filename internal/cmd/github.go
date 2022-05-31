package cmd

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v44/github"
	"github.com/mjehanno/new/internal/message"
)

const (
	DisplaySize = 25
	RequestSize = 100
	search      = "cookiecutter topic:cookiecutter-template"
)

var (
	TotalPageNumber int
	TotalRow        int
)

func Search(ctx context.Context, pageNumber int) tea.Cmd {
	opt := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: RequestSize,
			Page:    pageNumber,
		},
	}
	gclient := github.NewClient(nil)
	searchresult, _, err := gclient.Search.Repositories(ctx, search, &opt)
	if err != nil {
		return func() tea.Msg {
			return message.ErrorMessage{
				Error: err,
			}
		}
	}
	TotalRow = searchresult.GetTotal()
	TotalPageNumber = searchresult.GetTotal() / DisplaySize
	result := make([]message.SearchResult, 0)

	for _, r := range searchresult.Repositories {
		sr := message.SearchResult{}

		sr.Fullname = *r.FullName
		if r.Description != nil {
			sr.Description = *r.Description
		}
		if r.HTMLURL != nil {
			sr.HTMLURL = *r.HTMLURL
		}
		if r.Language != nil {
			sr.Language = *r.Language
		}
		sr.Stars = *r.StargazersCount
		result = append(result, sr)
	}

	if pageNumber == 1 {
		return func() tea.Msg {
			return message.FirstSearchResultMessage{
				SearchResults: result,
			}
		}
	} else if pageNumber == searchresult.GetTotal()/100 {
		return func() tea.Msg {
			return message.LastSearchResultMessage{
				SearchResults: result,
			}
		}
	}
	return func() tea.Msg {
		return message.SearchResultMessage{
			SearchResults: result,
		}
	}
}
