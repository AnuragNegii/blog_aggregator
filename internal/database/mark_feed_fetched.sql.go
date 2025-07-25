// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: mark_feed_fetched.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds 
SET last_fetched_at =CURRENT_TIMESTAMP , updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}
