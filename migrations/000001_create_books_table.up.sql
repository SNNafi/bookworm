CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    author text NOT NULL,
    publisher text NOT NULL,
    image text NOT NULL,
    cover_image text NOT NULL DEFAULT '',
    types text[] NOT NULL
);