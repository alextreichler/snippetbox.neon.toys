package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct{ DB *sql.DB }

func (s *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, CURRENT_TIMESTAMP, datetime(CURRENT_TIMESTAMP, '+' || ? || ' day'))`

	result, err := s.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func (s *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
