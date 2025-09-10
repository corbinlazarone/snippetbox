package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// snippet struct to represent a individual snippet type
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// snippet model that wraps a postgres db connection
type SnippetModel struct {
	DB *pgxpool.Pool
}

// insert a new snippet into the db, and returns the created snippet id
func (s *SnippetModel) Insert(Title string, Content string, Expiers int) (int, error) {
	statement := `INSERT INTO snippets (title, content, created, expires)
								VALUES ($1, $2, NOW() AT TIME ZONE 'UTC', (NOW() AT TIME ZONE 'UTC') + $3 * INTERVAL '1 day') RETURNING id;`
	var id int64
	err := s.DB.QueryRow(context.Background(), statement, Title, Content, Expiers).Scan(&id)
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (s *SnippetModel) Get(id int) (*Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets
								WHERE expires > NOW() AT TIME ZONE 'UTC' AND id = $1;`
	newSnip := &Snippet{}
	err := s.DB.QueryRow(context.Background(), statement, id).Scan(
		&newSnip.ID, &newSnip.Title, &newSnip.Content, &newSnip.Created, &newSnip.Expires,
	)
	if err != nil {
		return nil, ErrNoRecord
	}
	return newSnip, nil
}

// get the 10 latest snippets created
func (s *SnippetModel) Latest() ([]Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets
								WHERE expires > NOW() at TIME zone 'UTC' order by id desc limit 10;`

	rows, _ := s.DB.Query(context.Background(), statement)
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[Snippet])

	if err != nil {
		return nil, err
	}

	return res, nil
}
