-- name: CreatePost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- CREATE TABLE posts (
--     id UUID PRIMARY KEY,
--     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     title TEXT NOT NULL,
--     url TEXT UNIQUE NOT NULL,
--     description TEXT,
--     published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     feed_id UUID NOT NULL REFERENCES feeds (id) ON DELETE CASCADE
-- );
--


-- name: GetPostsForUser :many
SELECT posts.*, feeds.name AS feed_name FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
JOIN feeds ON posts.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
--
