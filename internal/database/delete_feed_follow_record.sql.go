// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: delete_feed_follow_record.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const deletefeedfollowrecord = `-- name: Deletefeedfollowrecord :exec
 
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1  
AND feed_id IN (SELECT feed_id FROM feeds WHERE feeds.url = $2)
`

type DeletefeedfollowrecordParams struct {
	UserID uuid.UUID
	Url    sql.NullString
}

func (q *Queries) Deletefeedfollowrecord(ctx context.Context, arg DeletefeedfollowrecordParams) error {
	_, err := q.db.ExecContext(ctx, deletefeedfollowrecord, arg.UserID, arg.Url)
	return err
}
