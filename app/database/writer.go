package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//go:generate mockery --name Writer
type Writer interface {
	InsertOrUpdateDecision(ctx context.Context, entry PutDecisionEntry) error
	UpdateUserTotalLikes(ctx context.Context, RecipientId string) error
	UpdateLikesAsViewed(ctx context.Context, RecipientId string, likes []DecisionModel) error
}

type DatabaseWriter struct {
	db *sql.DB
}

func NewDatabaseWriter(db *sql.DB) Writer {
	return DatabaseWriter{db}
}

const insertOrUpdateDecisionQuery = `
INSERT INTO decisions (
	actor_id, 
	recipient_id, 
	liked
) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE liked=?;
`

const updateUserTotalLikesQuery = `
UPDATE users 
SET likes = (
	SELECT count(*)
    FROM decisions
    WHERE recipient_id = ?
    AND liked = 1)
WHERE id = ?;
`

const updateDecisionLikesQuery = `
UPDATE decisions 
SET is_new = 0 
WHERE recipient_id = ? 
AND id in (%s) 
AND is_new = 1;
`

func (w DatabaseWriter) InsertOrUpdateDecision(ctx context.Context, entry PutDecisionEntry) error {
	rows, err := w.db.QueryContext(ctx, insertOrUpdateDecisionQuery,
		entry.ActorId,
		entry.RecipientId,
		entry.Like,
		entry.Like,
	)
	if err != nil {
		return fmt.Errorf("unable to insert or update decision: %w", err)
	}
	defer rows.Close()

	return nil
}

func (w DatabaseWriter) UpdateUserTotalLikes(ctx context.Context, RecipientId string) error {
	rows, err := w.db.QueryContext(ctx, updateUserTotalLikesQuery,
		RecipientId,
		RecipientId,
	)
	if err != nil {
		return fmt.Errorf("unable to update user total likes: %w", err)
	}
	defer rows.Close()

	return nil

}

func (w DatabaseWriter) UpdateLikesAsViewed(ctx context.Context, RecipientId string, likes []DecisionModel) error {
	decisionIds := make([]string, 0)

	for _, like := range likes {
		decisionIds = append(decisionIds, like.Id)
	}

	rows, err := w.db.QueryContext(ctx, fmt.Sprintf(updateDecisionLikesQuery, strings.Join(decisionIds, ", ")),
		RecipientId,
	)
	if err != nil {
		return fmt.Errorf("unable to update decisions: %w", err)
	}
	defer rows.Close()

	return nil

}

// Interface guards
var (
	_ Writer = (*DatabaseWriter)(nil)
)
