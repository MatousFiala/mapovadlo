-- +goose up 
create table obce (
    id text primary key,
    name text not null,
    orp_id text references orp(id) not null,
    geom GEOMETRY(Geometry, 4326) not null
);

create index obce_location_idx on obce using GIST (geom);

-- +goose down
drop table if exists obce;
