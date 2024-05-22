package models

import (
	"database/sql"
	"time"
)

// Define a Snippet type to hold the dataa for an individual snippet.
// The fields of the struct corresponds to the fields in our MySQL snippets table.
type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	/*
		Use the Exec() method on the embedded connection pool to execute the statement.
		Parameters are the SQL statement followed by the values for the placeholder parms(?).
		Ths returns a sql.Result type, which contains some basic info when the query got executed.
	*/
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// The LastInsertId method on the result to get the ID of the newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned is int64 so convert it to type int.
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}


