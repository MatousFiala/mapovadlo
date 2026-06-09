-- +goose Up
create extension if not exists postgis;

-- +goose Down
drop extension if exists postgis;
