-- +goose up 
create table orp (
    id text primary key,
    name text not null,
    kraj_id text references kraje(id) not null,
    geom GEOMETRY(Geometry, 4326) not null
);

create index orp_location_idx on orp using GIST (geom);

-- +goose down
drop table if exists orp;
