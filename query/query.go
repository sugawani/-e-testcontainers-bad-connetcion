package query

import (
	"database/sql"

	"main/models"
)

type Query struct {
	db *sql.DB
}

func NewQuery(db *sql.DB) *Query {
	return &Query{db: db}
}

func (q *Query) Execute(userID models.ID) (*models.User, error) {
	var u models.User
	row := q.db.QueryRow("SELECT * FROM users WHERE id = ?", userID)
	err := row.Scan(&u.ID, &u.Name)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
