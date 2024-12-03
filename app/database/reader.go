package database

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const Limit = 10

//go:generate mockery --name Reader
type Reader interface {
	FindLikesByRecipientIdPaginated(ctx context.Context, recipientId string, page int) ([]DecisionModel, error)
	FindNewLikesByRecipientIdPaginated(ctx context.Context, recipientId string, page int) ([]DecisionModel, error)
	GetIsMatch(ctx context.Context, ActorId string, RecipientId string) (bool, error)
	GetUserById(ctx context.Context, userId string) (UserModel, error)
	GetLimit() int
}

type DatabaseReader struct {
	db *sql.DB
}

func NewDatabaseReader(db *sql.DB) DatabaseReader {
	return DatabaseReader{db}
}

const readDecisionsWithLikeByRecipientIdPaginated = `
SELECT
	id,
	actor_id, 
	recipient_id, 
	liked,
	UNIX_TIMESTAMP(created_at) as created_at,
	UNIX_TIMESTAMP(updated_at) as updated_at
FROM decisions 
WHERE recipient_id = ?
AND liked = 1
LIMIT ?
OFFSET ?;
`

const readNewDecisionsWithLikeByRecipientIdPaginated = `
SELECT
	id,
	actor_id, 
	recipient_id, 
	liked,
	UNIX_TIMESTAMP(created_at) as created_at,
	UNIX_TIMESTAMP(updated_at) as updated_at
FROM decisions 
WHERE recipient_id = ?
AND liked = 1
AND is_new = 1
LIMIT ?
OFFSET ?;
`

const readActiveUsersById = `
SELECT 
	id,
    name,
    likes,
    gender,
    UNIX_TIMESTAMP(created_at) as created_at,
    UNIX_TIMESTAMP(updated_at) as updated_at,
	is_active
FROM users
WHERE id = ? 
AND is_active = 1;
`

const readGetMatchByActorIdAndRecipientId = `
SELECT
	id,
	actor_id, 
	recipient_id, 
	liked,
	UNIX_TIMESTAMP(created_at) as created_at,
	UNIX_TIMESTAMP(updated_at) as updated_at
FROM decisions
WHERE (
	actor_id = ?
	AND recipient_id = ?
	AND liked = 1)
OR (
	actor_id = ?
	AND recipient_id = ?
	AND liked = 1
);
`

// FindLikesByRecipientIdPaginated finds all likes on  decisions table for a given recipient user ID with pagination
func (r DatabaseReader) FindLikesByRecipientIdPaginated(ctx context.Context, recipientId string, page int) ([]DecisionModel, error) {
	offset := (page - 1) * Limit
	rows, err := r.db.QueryContext(ctx, readDecisionsWithLikeByRecipientIdPaginated, recipientId, Limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decisions []DecisionModel
	for rows.Next() {

		var decision DecisionModel
		err = rows.Scan(
			&decision.Id,
			&decision.ActorId,
			&decision.RecipientId,
			&decision.Liked,
			&decision.Created_at,
			&decision.Updated_at,
		)

		decisions = append(decisions, decision)

		if err != nil {
			return nil, err
		}
	}

	return decisions, nil
}

// FindNewLikesByRecipientIdPaginated finds all new/unchecked likes on  decisions table for a given recipient user ID with pagination
func (r DatabaseReader) FindNewLikesByRecipientIdPaginated(ctx context.Context, recipientId string, page int) ([]DecisionModel, error) {
	offset := (page - 1) * Limit
	rows, err := r.db.QueryContext(ctx, readNewDecisionsWithLikeByRecipientIdPaginated, recipientId, Limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decisions []DecisionModel
	for rows.Next() {

		var decision DecisionModel
		err = rows.Scan(
			&decision.Id,
			&decision.ActorId,
			&decision.RecipientId,
			&decision.Liked,
			&decision.Created_at,
			&decision.Updated_at,
		)

		decisions = append(decisions, decision)

		if err != nil {
			return nil, err
		}
	}

	return decisions, nil
}

func (r DatabaseReader) GetLimit() int {
	return Limit
}

// GetUserById get user information for a given user ID
func (r DatabaseReader) GetUserById(ctx context.Context, userId string) (UserModel, error) {
	rows, err := r.db.QueryContext(ctx, readActiveUsersById, userId)

	if err != nil {
		return UserModel{}, err
	}
	defer rows.Close()

	var user UserModel
	for rows.Next() {
		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Likes,
			&user.Gender,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IsAactive,
		)

		if err != nil {
			return UserModel{}, err
		}
	}

	return user, nil
}

// GetIsMatch check for match on decisions between actor user id and recipient user id and return if match is true or false
func (r DatabaseReader) GetIsMatch(ctx context.Context, ActorId string, RecipientId string) (bool, error) {
	rows, err := r.db.QueryContext(
		ctx,
		readGetMatchByActorIdAndRecipientId,
		ActorId,
		RecipientId,
		RecipientId,
		ActorId,
	)

	if err != nil {
		return false, err
	}
	defer rows.Close()

	var decisions []DecisionModel
	for rows.Next() {
		var decision DecisionModel
		err = rows.Scan(
			&decision.Id,
			&decision.ActorId,
			&decision.RecipientId,
			&decision.Liked,
			&decision.Created_at,
			&decision.Updated_at,
		)

		decisions = append(decisions, decision)

		if err != nil {
			return false, err
		}
	}

	if len(decisions) == 2 {
		return true, nil
	}
	return false, nil
}

// Interface guards
var (
	_ Reader = (*DatabaseReader)(nil)
)
