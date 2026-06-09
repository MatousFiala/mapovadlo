-- +goose up 
create type location_type as enum ('kraj', 'orp', 'obec');

create table locations (
    id text primary key,
    name text not null,
    type location_type not null,
    parent text references locations(id),
    geom GEOMETRY(Geometry, 4326) not null
);

create index location_idx on locations using GIST (geom);

-- +goose down
drop table if exists locations;
drop type if exists location_type;
