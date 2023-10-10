package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/weeaa/nft/database/repositories"
)

type DB struct {
	db         *pgxpool.Pool
	Indexer    *repositories.IndexerRepository
	Monitor    *repositories.MonitoredUsersRepository
	MonitorAll *repositories.MonitoredAllUsersRepository
}

func New(ctx context.Context, connString string) (*DB, error) {
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(ctx, "SELECT 1"); err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &DB{
		db:         db,
		Indexer:    repositories.NewFriendTechIndexerRepository(db),
		Monitor:    repositories.NewFriendTechMonitorRepository(db),
		MonitorAll: repositories.NewFriendTechMonitoredAllUsersRepository(db),
	}, nil
}

func (pg *DB) Disconnect() {
	pg.db.Close()
}
