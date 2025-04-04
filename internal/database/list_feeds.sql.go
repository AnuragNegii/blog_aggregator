// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: list_feeds.sql

package database

import (
	"context"
)

const listFeeds = `-- name: ListFeeds :many
SELECT feeds.name, feeds.url, users.name AS user_name
FROM feeds
JOIN users ON feeds.user_id = users.id
`

type ListFeedsRow struct {
	Name     string
	Url      string
	UserName string
}

func (q *Queries) ListFeeds(ctx context.Context) ([]ListFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, listFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListFeedsRow
	for rows.Next() {
		var i ListFeedsRow
		if err := rows.Scan(&i.Name, &i.Url, &i.UserName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
