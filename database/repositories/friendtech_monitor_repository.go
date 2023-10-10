package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/weeaa/nft/database/models"
)

type MonitoredUsersRepository struct {
	db *pgxpool.Pool
}

func NewFriendTechMonitorRepository(db *pgxpool.Pool) *MonitoredUsersRepository {
	return &MonitoredUsersRepository{db: db}
}

func (r *MonitoredUsersRepository) InsertUser(u *models.FriendTechMonitor, ctx context.Context) error {
	query := `INSERT INTO users (base_address, status, twitter_username, twitter_name, twitter_url, user_id) VALUES (@base_address, @status, @twitter_username, @twitter_name, @twitter_url, @user_id)`

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

func (r *MonitoredUsersRepository) GetUserByAddress(baseAddress string, ctx context.Context) (*models.FriendTechMonitor, error) {
	var user models.FriendTechMonitor

	query := `SELECT base_address, status, twitter_username, twitter_name, twitter_url, user_id FROM users WHERE base_address = $1`

	err := r.db.QueryRow(ctx, query, baseAddress).Scan(&user.BaseAddress, &user.Status, &user.TwitterUsername, &user.TwitterName, &user.TwitterURL, &user.UserID)

	if err != nil {
		return nil, fmt.Errorf("unable to fetch user: %w", err)
	}

	return &user, nil
}

func (r *MonitoredUsersRepository) GetAllAddresses(ctx context.Context) ([]string, error) {
	var baseAddresses []string

	query := "SELECT base_address FROM users"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var baseAddress string
		if err = rows.Scan(&baseAddress); err != nil {
			return nil, err
		}
		baseAddresses = append(baseAddresses, baseAddress)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return baseAddresses, nil
}

func (r *MonitoredUsersRepository) RemoveUser(baseAddress string, ctx context.Context) error {
	query := `DELETE FROM users WHERE base_address = $1`

	_, err := r.db.Exec(ctx, query, baseAddress)
	if err != nil {
		return fmt.Errorf("unable to remove user: %w", err)
	}

	return nil
}

func (r *MonitoredUsersRepository) RemoveAllUsers(ctx context.Context) error {
	query := `DELETE FROM users`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to remove users: %w", err)
	}

	return nil
}
