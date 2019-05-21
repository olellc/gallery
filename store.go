package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

const schemaSQL = `
DROP TABLE IF EXISTS fotos;

CREATE TABLE fotos(
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	raw_foto BLOB NOT NULL,
    preview BLOB
);`

func (s *Store) CreateSchema() error {
	_, err := s.db.Exec(schemaSQL)

	return err
}

func (s *Store) SaveFoto(rawFoto, preview []byte) (int64, error) {
	const query = `
		INSERT INTO fotos(raw_foto, preview) VALUES
			(?, ?)`
	res, err := s.db.Exec(query, rawFoto, preview)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// GetRawFoto returns nil slice if there is no foto with matching id
func (s *Store) GetRawFoto(id int64) ([]byte, error) {
	const query = `
		SELECT raw_foto
		FROM fotos WHERE id = ?`
	var raw_foto []byte
	err := s.db.QueryRow(query, id).Scan(&raw_foto)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return raw_foto, nil
}

// GetPreview returns nil slice if there is no preview with matching id
func (s *Store) GetPreview(id int64) ([]byte, error) {
	const query = `
		SELECT preview
		FROM fotos WHERE id = ?`
	var preview []byte
	err := s.db.QueryRow(query, id).Scan(&preview)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return preview, nil
}

func (s *Store) GetFotos() ([]int64, error) {
	rows, err := s.db.Query("SELECT id FROM fotos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []int64{}
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *Store) RemoveFoto(id int64) error {
	_, err := s.db.Exec("DELETE FROM fotos WHERE id = ?", id)

	return err
}
