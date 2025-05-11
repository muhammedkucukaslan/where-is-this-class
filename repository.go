package main

import (
	"database/sql"
	"os"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgreStore struct {
	db *sql.DB
}

func NewPostgreStore() (*PostgreStore, error) {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreStore{db: db}, nil
}

func (p *PostgreStore) GetClassRoom(code string, language string) (GetClassRoomResponse, error) {
	query := ` SELECT code, building, floor, floor_name, directions FROM classrooms WHERE code = $1 AND language = $2`

	var classroom GetClassRoomResponse
	if err := p.db.QueryRow(query, code, language).Scan(&classroom.Code, &classroom.Building, &classroom.Floor); err != nil {
		return GetClassRoomResponse{}, err
	}
	return classroom, nil
}

func (p *PostgreStore) CreateClassRoom(req *AddClassRoomRequest) error {

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	classRoomId := uuid.New()
	addClassRoomQuery := `INSERT INTO class_rooms (id, code, floor, image_url) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(addClassRoomQuery, classRoomId, req.Code, req.Floor, req.ImageUrl)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrClassRoomAlreadyExists
		}
		return err
	}

	addClassRoomTranslationQuery := `INSERT INTO class_room_translations (id, class_room_id, language, building,  description, detail) VALUES ($1, $2, $3, $4, $5, $6)`
	classRoomTranslationId := uuid.New()
	_, err = tx.Exec(addClassRoomTranslationQuery, classRoomTranslationId, classRoomId, req.Language, req.Building, req.Description, req.Detail)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
