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

type UpdateBit struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type InsertBit struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	DaysValid int    `json:"expires"`
}

type BitModel struct {
	DB             *sql.DB
	InsertStmt     *sql.Stmt
	UpdateStmt     *sql.Stmt
	SelectByIdStmt *sql.Stmt
}

func (m *BitModel) Insert(b *InsertBit) (int, error) {
	tx, err := m.DB.Begin()

	if err != nil {
		return -1, nil
	}

	defer tx.Rollback()

	sqlRes, err := tx.Stmt(m.InsertStmt).Exec(b.Title, b.Content, b.DaysValid)
	if err != nil {
		return -1, nil
	}

	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return -1, nil
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

func (m *BitModel) Latest() ([]*Bit, error) {
	stmt := `select id, title, content, created, expires 
	from bits where expires > utc_timestamp() order by id limit 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var bits []*Bit
	for rows.Next() {
		var b = &Bit{}
		err = rows.Scan(&b.Id, &b.Title, &b.Content, &b.CreatedAt, &b.ExpiresAt)
		if err != nil {
			return nil, err
		}
		bits = append(bits, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bits, nil
}

func (m *BitModel) Update(id int, b *UpdateBit) error {
	tx, err := m.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	var rid int
	err = tx.Stmt(m.SelectByIdStmt).QueryRow(id).Scan(&rid)
	if rid == 0 {
		return ErrNoRecord
	}
	if err != nil {
		return err
	}

	_, err = tx.Stmt(m.UpdateStmt).Exec(b.Title, b.Content, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func insertStmt(db *sql.DB) (*sql.Stmt, error) {
	insertQuery := `insert into bits (title, content, created, expires)
	values(?, ?, utc_timestamp(), date_add(utc_timestamp(), interval ? day))`
	stmt, err := db.Prepare(insertQuery)

	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func updateStmt(db *sql.DB) (*sql.Stmt, error) {
	updateQuery := `update bits set title = ?, content = ? where id = ?`
	stmt, err := db.Prepare(updateQuery)

	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func selectByIdStmt(db *sql.DB) (*sql.Stmt, error) {
	query := `select count(id) from bits where id = ?`
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func CreateBitModel(db *sql.DB) (*BitModel, error) {

	updateStmt, err := updateStmt(db)
	if err != nil {
		return nil, err
	}

	selectByIdStmt, err := selectByIdStmt(db)
	if err != nil {
		return nil, err
	}

	inserStmt, err := insertStmt(db)
	if err != nil {
		return nil, err
	}

	return &BitModel{
		DB:             db,
		InsertStmt:     inserStmt,
		UpdateStmt:     updateStmt,
		SelectByIdStmt: selectByIdStmt,
	}, nil
}
