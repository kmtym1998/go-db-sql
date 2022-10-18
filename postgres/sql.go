package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/guregu/null"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBCfg struct {
	MaxIdleConns    null.Int
	MaxOpenConns    null.Int
	ConnMaxLifetime null.Int
}

func New(url string, cfg DBCfg) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)

	if err != nil {
		return db, err
	}

	if err := db.Ping(); err != nil {
		return db, err
	}

	if cfg.MaxIdleConns.Valid {
		db.SetMaxIdleConns(int(cfg.MaxIdleConns.Int64))
	}

	if cfg.MaxOpenConns.Valid {
		db.SetMaxOpenConns(int(cfg.MaxOpenConns.Int64))
	}

	if cfg.ConnMaxLifetime.Valid {
		db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime.Int64))
	}

	return db, nil
}

func Tx(ctx context.Context, db *sql.DB, do func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	if err := do(tx); err != nil {
		if innerErr := tx.Rollback(); innerErr != nil {
			return fmt.Errorf("tx: rollback error: %w (outer error: %v)", innerErr, err)
		}
		return err
	}
	return tx.Commit()
}
