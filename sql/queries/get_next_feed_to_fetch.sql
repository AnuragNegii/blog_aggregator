-- name: GetNextFeedToFetch :one

SELECT * FROM feeds
WHERE user_id = $1
ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;
