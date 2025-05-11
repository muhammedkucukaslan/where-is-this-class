package main

import (
	"database/sql"
	"errors"
	"fmt"
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

	tx, err := p.db.Begin()
	if err != nil {
		return GetClassRoomResponse{}, err
	}
	defer tx.Rollback()

	query := ` select  building,floor, image_url, description, detail 
	from class_rooms cr 
	join class_room_translations crt on cr.id = crt.class_room_id
	where code = $1 AND  language = $2`

	var classroom GetClassRoomResponse
	if err := tx.QueryRow(query, code, language).Scan(&classroom.Building, &classroom.Floor, &classroom.ImageUrl, &classroom.Description, &classroom.Detail); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GetClassRoomResponse{}, ErrClassRoomNotFound
		}
		return GetClassRoomResponse{}, err
	}

	visitCountQuery := `UPDATE class_rooms SET visited = visited + 1 WHERE code = $1`
	_, err = tx.Exec(visitCountQuery, code)
	if err != nil {
		return GetClassRoomResponse{}, err
	}

	if err := tx.Commit(); err != nil {
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
	addClassRoomTranslationQuery := `INSERT INTO class_room_translations (id, class_room_id, language, building,  description, detail) VALUES `
	for _, translation := range req.Translations {
		values := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s'),", uuid.New(), classRoomId, translation.Language, translation.Building, translation.Description, translation.Detail)
		addClassRoomTranslationQuery += values
	}
	_, err = tx.Exec(addClassRoomTranslationQuery[:len(addClassRoomTranslationQuery)-1])
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
