package ms_sql

import (
	"context"
	"database/sql"
	"fmt"
	"pl1x/pkg/models"
	"strconv"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	ctx := context.Background()

	err := m.DB.PingContext(ctx)
	if err != nil {
		return 0, err
	}

	tsql := `INSERT INTO [GoBlojikDB].[dbo].[snippets] ([title], [content], [created], [expires]) 
	VALUES
	(@title, @content, CURRENT_TIMESTAMP, DATEADD(day, @range, CURRENT_TIMESTAMP));`

	stmt, err := m.DB.Prepare(tsql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rng, _ := strconv.Atoi(expires)

	row := stmt.QueryRowContext(ctx,
		sql.Named("title", title),
		sql.Named("content", content),
		sql.Named("range", rng))

	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return 0, err
	}

	return int(newID), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := m.DB.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	tsql := fmt.Sprintf(`USE GoBlojikDB;
	SELECT id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP AND id = @id 
	`)

	// Execute query
	rows, err := m.DB.QueryContext(ctx, tsql, sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var count int
	s := &models.Snippet{}

	// Iterate through the result set.
	for rows.Next() {
		// Get values from row.
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		count++
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
