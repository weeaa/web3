package services

import (
	"github.com/weeaa/nft/database/db"
)

type UserService struct {
	DB *db.DB
}

func NewUserService(db *db.DB) *UserService {
	return &UserService{DB: db}
}
