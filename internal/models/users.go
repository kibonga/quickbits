package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const duplicateEntry int = 1062

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword string
	Created        time.Time
}

type UserModelDB struct {
	DB         *sql.DB
	InsertStmt *sql.Stmt
	AuthStmt   *sql.Stmt
	ExistsStmt *sql.Stmt
}

type UserSignupModel struct {
	Name     string
	Email    string
	Password string
}

type UserLoginModel struct {
	Email   string
	Pasword string
}

func (m *UserModelDB) Insert(u *UserSignupModel) (int, error) {
	tx, err := m.DB.Begin()
	defer tx.Rollback()

	if err != nil {
		return 0, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	sqlRes, err := m.InsertStmt.Exec(u.Name, u.Email, hash)
	if err != nil {
		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) {
			if int(mySqlErr.Number) == duplicateEntry {
				return 0, ErrDuplicateEmail
			}
		}
		return 0, err
	}

	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, nil
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModelDB) Auth(u *UserLoginModel) (int, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	var hash string

	if err := m.AuthStmt.QueryRow(u.Email).Scan(&id, &hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCreds
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Pasword)); err != nil {
		return 0, ErrInvalidCreds
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *UserModelDB) Exists(id int) (bool, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return false, nil
	}
	defer tx.Rollback()

	var exists bool

	if err := m.ExistsStmt.QueryRow(id).Scan(&exists); err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, nil
	}

	return exists, nil
}

func userInsertStmt(db *sql.DB) (*sql.Stmt, error) {
	query := `insert into users (name, email, hashed_password, created)
	values(?, ?, ?, utc_timestamp())`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func userExistsStmt(db *sql.DB) (*sql.Stmt, error) {
	query := "select exists(select 1 from users where id = ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func userAuthStmt(db *sql.DB) (*sql.Stmt, error) {
	query := "select id, hashed_password from users where email = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func UserModelDb(db *sql.DB) (*UserModelDB, error) {

	userInsertStmt, err := userInsertStmt(db)
	if err != nil {
		return nil, err
	}

	userExistsStmt, err := userExistsStmt(db)
	if err != nil {
		return nil, err
	}

	userAuthStmt, err := userAuthStmt(db)
	if err != nil {
		return nil, err
	}

	return &UserModelDB{
		DB:         db,
		InsertStmt: userInsertStmt,
		ExistsStmt: userExistsStmt,
		AuthStmt:   userAuthStmt,
	}, nil
}
