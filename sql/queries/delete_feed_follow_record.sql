-- name: Deletefeedfollowrecord :exec
 
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1  
AND feed_id IN (SELECT feed_id FROM feeds WHERE feeds.url = $2);
