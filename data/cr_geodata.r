library("tidyverse")
library("RCzechia")
library("czso")
library("RPostgres")
source("./dbConn.r")

cat <- czso::czso_get_catalogue() |> filter(str_detect(dataset_id, "cis"))

kraje <- RCzechia::kraje() # 100
orp <- RCzechia::orp_polygony() # 65
obce <- RCzechia::obce_polygony() # 43


kraje
orp
obce

locations <- bind_rows(
kraje |>
    mutate(id = paste0("100-", KOD_KRAJ), type = "kraj", parent = NA_character_) |>
    select(id, name = NAZ_CZNUTS3, type, parent, geom = geometry),
orp |>
    mutate(id = paste0("65-", KOD_ORP), type = "orp", parent = paste0("100-", KOD_KRAJ)) |>
    select(id, name = NAZ_ORP, type, parent, geom = geometry),
obce |>
    mutate(id = paste0("43-", KOD_OBEC), type = "obec", parent = paste0("65-", KOD_ORP)) |>
    select(id, name = NAZ_ORP, type, parent, geom = geometry))

con <- get_con()

dbExecute(con, "delete from locations;")

st_write(
    obj = locations,
    dsn = con,
    layer = "locations",
    append = TRUE,
    delete_layer = FALSE
)

dbDisconnect(con)
