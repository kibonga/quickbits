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

func (m *BitModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

func (m *BitModel) Get(id int) (*Bit, error) {
	return nil, nil
}

func (m *BitModel) Latest() ([]*BitModel, error) {
	return nil, nil
}
