package main

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreStore struct {
	db *gorm.DB
}

func NewPostgreStore() (*PostgreStore, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PostgreStore{db: db}, nil
}

func (p *PostgreStore) GetClassRoom(code string) (GetClassRoomResponse, error) {
	var classroom GetClassRoomResponse
	if err := p.db.Where("code = ?", code).First(&classroom).Error; err != nil {
		return GetClassRoomResponse{}, err
	}
	return classroom, nil
}

func (p *PostgreStore) CreateClassRoom(req *AddClassRoomRequest) error {
	classroom := GetClassRoomResponse{
		Code:       req.Code,
		Building:   req.Building,
		Floor:      req.Floor,
		FloorName:  req.FloorName,
		Directions: req.Directions,
	}
	if err := p.db.Create(&classroom).Error; err != nil {
		return err
	}
	return nil
}
