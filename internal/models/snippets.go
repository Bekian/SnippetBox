package models

import (
	"database/sql"
	"time"
)

// individual snippet
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// snippet model, sql db connection pool wrapper
type SnippetModel struct {
	DB *sql.DB
}

// row inserter
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// snippet getter
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// return 10 most recent snippets
// i think of this like Pandas.DataFrame.Head()
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
