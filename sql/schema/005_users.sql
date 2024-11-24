-- +goose Up
CREATE SCHEMA IF NOT EXISTS gator;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE gator.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP SCHEMA gator CASCADE;

