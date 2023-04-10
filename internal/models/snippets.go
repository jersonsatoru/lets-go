package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func NewSnippetModel(db *sql.DB) *SnippetModel {
	return &SnippetModel{
		DB: db,
	}
}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `
		INSERT INTO snippets 
		(title, content, expires, created) 
		VALUES ($1, $2, now() + INTERVAL '1 day' * $3, now())
		RETURNING id
	`
	result := sm.DB.QueryRow(stmt, title, content, expires)
	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (sm *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `
		SELECT id, title, content, expires, created
		FROM snippets 
		WHERE expires > now() AND id = $1
	`

	var snippet Snippet
	result := sm.DB.QueryRow(stmt, id)
	err := result.Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Expires,
		&snippet.Created,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &snippet, nil
}

func (sm *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `
		SELECT id, title, content, expires, created
		FROM snippets 
		WHERE expires > now()
	`

	result, err := sm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var snippets []*Snippet

	for result.Next() {
		snippet := &Snippet{}
		err := result.Scan(
			&snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.Expires,
			&snippet.Created,
		)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
