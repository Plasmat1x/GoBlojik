package ms_sql

import (
	"context"
	"database/sql"
	"pl1x/pkg/models"
	"strconv"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) getLast() (int64, error) {
	lastInsertId := 0
	err := m.DB.QueryRow("SELECT TOP(1) id FROM snippet").Scan(&lastInsertId)
	return int64(lastInsertId), err
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
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
