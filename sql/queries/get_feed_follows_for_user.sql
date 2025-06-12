-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name AS feedname
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;
