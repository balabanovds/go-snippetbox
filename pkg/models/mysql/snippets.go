package mysql

import (
	"database/sql"
	"errors"
	"github.com/balabanovds/go-snippetbox/pkg/models"
)

// SnippetModel wraps db pool and have some methods working with db
type SnippetModel struct {
	DB *sql.DB
}

func NewSnippetModel(db *sql.DB) *SnippetModel {
	return &SnippetModel{DB: db}
}

// Insert new record
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt, err := m.DB.Prepare("INSERT INTO snippets (title, content, created, expires) " +
		"VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))")
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get one record by id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	row := m.DB.QueryRow("SELECT id, title, content, created, expires FROM snippets "+
		"WHERE expires > UTC_TIMESTAMP AND id = ?", id)

	s := &models.Snippet{}

	if err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

// Latest will return last 10 records
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	rows, err := m.DB.Query("SELECT id, title, content, created, expires FROM snippets " +
		"WHERE expires > UTC_TIMESTAMP ORDER BY created DESC LIMIT 10")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*models.Snippet

	for rows.Next() {
		s := &models.Snippet{}

		if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
