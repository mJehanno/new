package message

type SearchResult struct {
	Fullname    string
	Description string
	HTMLURL     string
	Stars       int
	Language    string
}

type SearchResultMessage struct {
	SearchResults []SearchResult
}
