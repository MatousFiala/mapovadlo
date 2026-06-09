-- name: GetFeaturesByLayerSeriesAndBbox :many
select loc.id, loc.name, ST_AsGeoJSON(ST_SimplifyPreserveTopology(loc.geom, $8))::json as geom, val.value
from locations loc
    join values val on loc.id = val.location
where loc.type = $1
    and val.series = $2
    and val.time = $3
    and ST_Intersects(geom, ST_MakeEnvelope($4, $5, $6, $7, 4326));
