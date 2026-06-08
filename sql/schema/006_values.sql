-- +goose Up 
create table values (
    series text references series(id),

);

-- +goose Down
drop table values;
