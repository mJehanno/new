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

func InitSearch(ctx context.Context) tea.Cmd {
	opt := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: RequestSize,
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
	return func() tea.Msg {
		return message.SearchResultMessage{
			SearchResults: result,
		}
	}
}

func SearchPage(ctx context.Context, pageNumber int) tea.Cmd {
	gclient := github.NewClient(nil)
	opt := github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    pageNumber,
			PerPage: RequestSize,
		},
	}
	pagedResult, _, err := gclient.Search.Repositories(ctx, search, &opt)
	if err != nil {
		return func() tea.Msg {
			return message.ErrorMessage{
				Error: err,
			}
		}
	}
	result := make([]message.SearchResult, 0)
	for _, r := range pagedResult.Repositories {
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
	return func() tea.Msg {
		return message.SearchResultMessage{
			SearchResults: result,
		}
	}
}
