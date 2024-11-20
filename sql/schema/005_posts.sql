-- +goose Up
CREATE TABLE posts (
	id uuid PRIMARY KEY,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	title TEXT,
	url TEXT UNIQUE NOT NULL,
	description TEXT,
	published_at timestamp NOT NULL,
	feed_id uuid NOT NULL REFERENCES feeds (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
