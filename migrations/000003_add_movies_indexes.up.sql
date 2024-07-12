CREATE INDEX IF NOT EXISTS books_name_idx ON books USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS books_types_idx ON books USING GIN (types);