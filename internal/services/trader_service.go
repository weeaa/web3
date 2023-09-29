package services

import (
	"github.com/weeaa/nft/db"
)

type UserService struct {
	DB *db.DB
}

func NewUserService(db *db.DB) *UserService {
	return &UserService{DB: db}
}
