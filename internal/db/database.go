package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"summer/practice2022/internal/config"
	st "summer/practice2022/internal/structures"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(cfg config.Config) (*Database, error) {
	sqlBd, err := sql.Open("sqlite3", cfg.DBUri)
	if err != nil {
		return nil, nil
	}

	db := &Database{
		db: sqlBd,
	}

	return db, nil
}

func (db *Database) GetHashedPassword(
	login string,
) (
	st.HashedPassword,
	error,
) {
	row := db.db.QueryRow(
		`SELECT
			password_hash,
			password_salt
		FROM users
		WHERE login = ?`,
		login,
	)

	if err := row.Err(); err != nil {
		return st.HashedPassword{}, err
	}

	hashedPassword := st.HashedPassword{}

	err := row.Scan(
		&hashedPassword.Hash,
		&hashedPassword.Salt,
	)
	if err != nil {
		err = errors.Wrap(err, "scan fields")
		return st.HashedPassword{}, err
	}

	return hashedPassword, nil
}

func (db *Database) SaveRefreshTokens(refresh, login string) error {
	_, err := db.db.Exec(
		`INSERT INTO refresh_tokens (
			token,
			login,
			is_used
		) VALUES (
			?,
			?,
			?
		)`,
		refresh,
		login,
		false,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetLoginByRefreshToken(
	refresh string,
) (
	string,
	bool,
	error,
) {
	row := db.db.QueryRow(
		`SELECT
			login,
			is_used
		FROM refresh_tokens
		WHERE token = ?`,
		refresh,
	)

	if err := row.Err(); err != nil {
		return "", false, err
	}

	var (
		login  string
		isUsed bool
	)
	if err := row.Scan(&login, &isUsed); err != nil {
		return "", false, err
	}

	return login, isUsed, nil
}

func (db *Database) MarkTokenAsUsed(refresh string) error {
	_, err := db.db.Exec(
		`UPDATE refresh_tokens
		SET is_used = true
		WHERE token = ?`,
		refresh,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) LogAlreadyUsedRefreshToken(refresh string) error {
	_, err := db.db.Exec(
		`INSERT INTO already_used_refresh_tokens (
			token,
			detected_at
		) VALUES (
			?,
			?
		)`,
		refresh,
		time.Now().Unix(),
	)
	if err != nil {
		return err
	}

	return nil
}
