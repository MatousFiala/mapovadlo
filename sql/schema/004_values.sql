-- +goose Up 
create table values (
    series text references series(id) not null,
    location text references locations(id) not null,
    time date not null,
    value double precision,
    primary key (series, location, time)
);

-- +goose Down
drop table values;
