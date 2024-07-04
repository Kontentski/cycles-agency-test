-- +goose Up
-- +goose StatementBegin
CREATE TABLE burgers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    image_url TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE burgers;
-- +goose StatementEnd
