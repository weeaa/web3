package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/weeaa/nft/database/models"
)

type MonitoredAllUsersRepository struct {
	db *pgxpool.Pool
}

func NewFriendTechMonitoredAllUsersRepository(db *pgxpool.Pool) *MonitoredAllUsersRepository {
	return &MonitoredAllUsersRepository{db: db}
}

func (r *MonitoredAllUsersRepository) InsertUser(u *models.FriendTechMonitorAll, ctx context.Context) error {
	query := `INSERT INTO usersall (base_address, status, twitter_username, twitter_name, twitter_url, user_id) VALUES (@base_address, @status, @twitter_username, @twitter_name, @twitter_url, @user_id)`

	args := pgx.NamedArgs{
		"base_address":     u.BaseAddress,
		"status":           u.Status,
		"twitter_username": u.TwitterUsername,
		"twitter_name":     u.TwitterName,
		"twitter_url":      u.TwitterURL,
		"user_id":          u.UserID,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (r *MonitoredAllUsersRepository) GetUserByAddress(baseAddress string, ctx context.Context) (*models.FriendTechMonitorAll, error) {
	var user models.FriendTechMonitorAll

	query := `SELECT base_address, status, twitter_username, twitter_name, twitter_url, user_id FROM usersall WHERE base_address = $1`

	err := r.db.QueryRow(ctx, query, baseAddress).Scan(&user.BaseAddress, &user.Status, &user.TwitterUsername, &user.TwitterName, &user.TwitterURL, &user.UserID)

	if err != nil {
		return nil, fmt.Errorf("unable to fetch user: %w", err)
	}

	return &user, nil
}
