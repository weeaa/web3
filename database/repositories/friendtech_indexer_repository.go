package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/weeaa/nft/database/models"
)

type IndexerRepository struct {
	db *pgxpool.Pool
}

func NewFriendTechIndexerRepository(db *pgxpool.Pool) *IndexerRepository {
	return &IndexerRepository{db: db}
}

func (r *IndexerRepository) InsertUser(u *models.FriendTechIndexer, ctx context.Context) error {
	query := `INSERT INTO indexer (base_address, twitter_username, user_id) VALUES (@base_address, @twitter_username, @user_id)`

	args := pgx.NamedArgs{
		"base_address":     u.BaseAddress,
		"twitter_username": u.TwitterUsername,
		"user_id":          u.UserID,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (r *IndexerRepository) GetUserByAddress(baseAddress string, ctx context.Context) (*models.FriendTechIndexer, error) {
	var user models.FriendTechIndexer

	query := `SELECT user_id, twitter_username, base_address FROM indexer WHERE base_address = $1`

	err := r.db.QueryRow(ctx, query, baseAddress).
		Scan(&user.UserID, &user.TwitterUsername, &user.BaseAddress)

	if err != nil {
		return nil, fmt.Errorf("unable to fetch user: %w", err)
	}

	return &user, nil
}

func (r *IndexerRepository) GetAllUsers(ctx context.Context) ([]*models.FriendTechIndexer, error) {
	var u []*models.FriendTechIndexer

	query := `SELECT * FROM indexer`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to insert row: %w", err)
	}

	for rows.Next() {
		var indexer *models.FriendTechIndexer

		if err = rows.Scan(
			&indexer.UserID,
			&indexer.TwitterUsername,
			&indexer.BaseAddress,
		); err != nil {
			//log.Fatalf("Scan error: %v", err)
		}

		u = append(u, indexer)
	}

	return u, nil
}
