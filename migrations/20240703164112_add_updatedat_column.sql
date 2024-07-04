-- +goose Up
-- +goose StatementBegin
ALTER TABLE burgers
ADD COLUMN updated_at TIMESTAMP DEFAULT current_timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE burgers
DROP COLUMN IF EXISTS created_at;
-- +goose StatementEnd
