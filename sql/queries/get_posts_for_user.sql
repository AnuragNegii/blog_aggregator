-- name: GetPostsForUsers :many
SELECT * FROM posts
JOIN feeds ON feeds.id = posts.feed_id
JOIN feed_follows ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY published_at DESC LIMIT $2;
