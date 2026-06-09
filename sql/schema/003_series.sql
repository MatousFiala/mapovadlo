-- +goose Up 
create table series (
    id text primary key,
    name text not null,
    unit text not null,
    description text,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table series;
