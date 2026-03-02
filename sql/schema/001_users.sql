-- +goose Up

create table users(
    id UUID primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    email Text not null unique
);

-- +goose Down
DROP TABLE users;