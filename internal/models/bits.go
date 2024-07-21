package models

import (
	"database/sql"
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
	return nil, nil
}

func (m *BitModel) Latest() ([]*BitModel, error) {
	return nil, nil
}
