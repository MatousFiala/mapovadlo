-- +goose up 
create table kraje (
    kraj_id text primary key,
    kraj_nazev text not null,
    geom GEOMETRY(Geometry, 4326) not null
);

create index kraje_location_idx on kraje using GIST (geom);

-- +goose down
drop table if exists kraje;
