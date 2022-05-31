package message

type SearchResult struct {
	Fullname    string
	Description string
	HTMLURL     string
	Stars       int
	Language    string
}

type FirstSearchResultMessage struct {
	SearchResults []SearchResult
}

type SearchResultMessage struct {
	SearchResults []SearchResult
}

type LastSearchResultMessage struct {
	SearchResults []SearchResult
}
