-- name: GetFeed :one
SELECT * FROM feeds
WHERE  feeds.url = $1;
