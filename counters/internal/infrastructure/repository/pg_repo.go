package repository

import (
	"context"
	"errors"

	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGCounterRepository struct {
	db *pgxpool.Pool
}

func NewPGCounterRepository(db *pgxpool.Pool) *PGCounterRepository {
	return &PGCounterRepository{db: db}
}

func (r *PGCounterRepository) GetCounter(ctx context.Context, dialogID string) (*domain.DialogCounter, error) {
	var counter domain.DialogCounter
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("Started", "dialogID", dialogID)

	err := r.db.QueryRow(ctx, "SELECT dialog_id, unread_count FROM dialog_counters WHERE dialog_id = $1", dialogID).
		Scan(&counter.DialogID, &counter.UnreadCount)
	if errors.Is(err, pgx.ErrNoRows) {
		counter = domain.DialogCounter{DialogID: dialogID, UnreadCount: 0}
	} else if err != nil {
		return nil, err
	}
	log.Debugw("Success", "counter", counter)
	return &counter, nil
}

func (r *PGCounterRepository) SaveCounter(ctx context.Context, counter *domain.DialogCounter) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("Started", "counter", counter)

	_, err := r.db.Exec(ctx, `
		INSERT INTO dialog_counters (dialog_id, unread_count) 
		VALUES ($1, $2) 
		ON CONFLICT (dialog_id) 
		DO UPDATE SET unread_count = EXCLUDED.unread_count
	`, counter.DialogID, counter.UnreadCount)

	if err != nil {
		log.Debugw("r.db.Exec() returned error", "err", err)
		return err
	}
	log.Debug("Success")
	return nil

}
