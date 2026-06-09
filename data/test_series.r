library("tidyverse")
library("obcevkostce")
library("czso")
library("RCzechia")


obyvatele <- obcevkostce::get_data("POCET-OBYVATEL") |> mutate(municipality_id = as.character(municipality_id))
obyvatele

obce <- RCzechia::obce_body() |> st_drop_geometry() |> as_tibble()

obyvplus <- obyvatele |>
    inner_join(obce, by = c("municipality_id" = "KOD_OBEC"))

pocet_obyv <- bind_rows(
obyvplus |>
    mutate(location = paste0("43-", municipality_id), time = as_date(paste0(year,"-01-01")), series = "pocet-obyvatel") |>
    select(location, series, time, value),
obyvplus |>
    summarise(.by = c(KOD_ORP, year), value = sum(value)) |>
    mutate(location = paste0("65-", KOD_ORP), time =con as_date(paste0(year,"-01-01")), series = "pocet-obyvatel") |>
    select(location, series, time, value),
obyvplus |>
    summarise(.by = c(KOD_KRAJ, year), value = sum(value)) |>
    mutate(location = paste0("100-", KOD_KRAJ), time = as_date(paste0(year,"-01-01")), series = "pocet-obyvatel") |>
    select(location, series, time, value)
)

source("./dbConn.r")
con <- get_con()

dbExecute(con, "delete from values where series = 'pocet-obyvatel';")
dbExecute(con, "insert into series (id, name, unit, created_at, updated_at) values ('pocet-obyvatel', 'Počet obyvatel', 'osob', NOW(), NOW());")
dbWriteTable(con, "values", pocet_obyv, append = T, row.names = F)

dbDisconnect(con)



