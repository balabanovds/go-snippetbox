package mysql

import (
	"database/sql"
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
	return 0, nil
}

// Get one record by id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Latest will return last 10 records
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
