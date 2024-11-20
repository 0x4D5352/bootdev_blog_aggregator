-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT * FROM posts
WHERE feed_id = ANY($1::uuid[])
ORDER BY published_at LIMIT $2;

-- name: GetPostByURL :one
SELECT * FROM posts
WHERE url = $1 LIMIT 1;

-- name: GetPostByID :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: ResetPosts :exec
DELETE FROM posts;

-- name: GetPosts :many
SELECT * FROM posts;
