-- +goose Up
-- +goose StatementBegin
    CREATE TABLE IF NOT EXISTS urls (
           id INTEGER PRIMARY KEY,
           url TEXT NOT NULL,
           alias VARCHAR(255) UNIQUE,
           create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
