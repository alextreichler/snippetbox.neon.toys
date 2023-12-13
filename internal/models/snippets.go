package models

import (
	"database/sql"
	"errors"
	"fmt"
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

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	// Construct the query string with actual values for debugging
	debugStmt := fmt.Sprintf("INSERT INTO snippets (title, content, created, expires) VALUES('%s', '%s', CURRENT_TIMESTAMP, datetime(CURRENT_TIMESTAMP, '+%d day'))", title, content, expires)

	// stmt := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, CURRENT_TIMESTAMP, datetime(CURRENT_TIMESTAMP, '+' || ? || 'day'))`
	result, err := sm.DB.Exec(debugStmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > datetime('now') AND id = ?;`

	row := sm.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (sm *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > datetime('now') ORDER BY id DESC LIMIT 10`

	rows, err := sm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
