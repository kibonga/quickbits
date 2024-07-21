package models

import (
	"database/sql"
	"errors"
	"time"
)

type Bit struct {
	Id        int
	Title     string
	Content   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type BitModel struct {
	DB *sql.DB
}

func (m *BitModel) Insert(title string, content string, daysValid int) (int, error) {
	stmt := `insert into bits (title, content, created, expires)
	values(?, ?, utc_timestamp(), date_add(utc_timestamp(), interval ? day))`

	sqlRes, err := m.DB.Exec(stmt, title, content, daysValid)
	if err != nil {
		return 0, err
	}

	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *BitModel) Get(id int) (*Bit, error) {
	stmt := `select id, title, content, created, expires
	from bits where id = ? and expires > utc_timestamp()`

	row := m.DB.QueryRow(stmt, id)

	b := &Bit{}

	err := row.Scan(&b.Id, &b.Title, &b.Content, &b.CreatedAt, &b.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return b, nil
}

func (m *BitModel) Latest() ([]*BitModel, error) {
	return nil, nil
}
