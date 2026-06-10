package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/MatousFiala/mapovadlo/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Bbox struct {
	XMin float64
	YMin float64
	XMax float64
	YMax float64
}

type apiConfig struct {
	dbQueries *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
	    panic(fmt.Errorf("Could not open connection to db! %w", err))
	}
	ApiConfig := apiConfig{dbQueries: database.New(db)}
 
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", getIndex)
	mux.HandleFunc("GET /api/features", ApiConfig.getFeatures)

	server := &http.Server{Handler: mux, Addr: ":8080"}
	err = server.ListenAndServe()
	if err != nil {
	    panic(err)
	}
	
}

func getIndex(w http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}

func (cfg *apiConfig) getFeatures(w http.ResponseWriter, req *http.Request) {
	series := req.URL.Query().Get("series")
	if series == "" {
		respondWithError(w, 400, "Series field is needed")
		return
	}
	timeRaw := req.URL.Query().Get("time")
	const layout = "2006-01-02"
	time, err := time.Parse(layout, timeRaw)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Can't parse %s as time, format YYYY-MM-DD expected", timeRaw))
		return
	}
	granularity := database.LocationType(req.URL.Query().Get("granularity"))
	if granularity != database.LocationTypeKraj && granularity != database.LocationTypeOrp && granularity != database.LocationTypeObec {
		respondWithError(w, 400, "Invalid granularity, only kraj, orp and obec supported")
	}
	bboxRaw := req.URL.Query().Get("bbox")
	var bbox Bbox
	if granularity != database.LocationTypeObec {
		bbox = Bbox{XMin: 12, YMin: 48, XMax: 19, YMax: 52}
	} else {
		bbox, err = parseBboxFromString(bboxRaw)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error parsing bbox: %v", err))
			return
		}
	}
	zoom, err := strconv.ParseFloat(req.URL.Query().Get("zoom"), 64)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Can't parse %s as zoom level", req.URL.Query().Get("zoom")))
		return
	}
	tolerance := getSimplifyTolerance(zoom)

	features, err := cfg.dbQueries.GetFeaturesByLayerSeriesAndBbox(req.Context(), database.GetFeaturesByLayerSeriesAndBboxParams{
		Type: granularity,
		Series: series,
		Time: time,
		StMakeenvelope: bbox.XMin,
		StMakeenvelope_2: bbox.YMin,
		StMakeenvelope_3: bbox.XMax,
		StMakeenvelope_4: bbox.YMax,
		StSimplifypreservetopology: tolerance, 
		})


	geojsonFeatures := make([]map[string]any, 0, len(features))
	for _, row := range features {
		geojsonFeatures = append(geojsonFeatures, map[string]any{
			"type": "Feature",
			"geometry": json.RawMessage(row.Geom),
			"properties": map[string]any{
				"id": row.ID,
				"name": row.Name,
				"value": row.Value.Float64,
			},
		})
	}
	responseMap := map[string]any{
		"type": "FeatureCollection",
		"features": geojsonFeatures,
	}

	log.Printf("Got request with params series: %s, time: %v, granularity %s, bbox %s (using %v). Returning %d features\n", series, time, granularity, bboxRaw, bbox, len(geojsonFeatures))

	respondWithJson(w, 200, responseMap)
}

func getSimplifyTolerance(zoomLevel float64) float64 {
	return 1.0 / math.Pow(2, zoomLevel)
}

func parseBboxFromString(bboxString string) (Bbox, error) {
	parts := strings.Split(bboxString, ",")
	if len(parts) != 4 {
		return Bbox{}, fmt.Errorf("Bounding box must be in format 'xmin,ymin,xmax,ymax', could not parse '%s' as 4 coords separated by commas", bboxString)
	}
	var bbox Bbox
	var err error
	values := []*float64{&bbox.XMin, &bbox.YMin, &bbox.XMax, &bbox.YMax}
	for i, p := range parts {
		*values[i], err = strconv.ParseFloat(p, 64)
		if err != nil {
			return Bbox{}, fmt.Errorf("Bbox coord could not be parsed as float: %v", p)
		}
	}
	if bbox.XMin >= bbox.XMax || bbox.YMin >= bbox.YMax {
		return Bbox{}, fmt.Errorf("Degenerate bbox")
	}
	return bbox, nil
}



func respondWithJson(w http.ResponseWriter, code int, payload any) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil { return err }

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payloadJson)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, errorMessage string) error {
	err := respondWithJson(w, code, map[string]string{"error": errorMessage})
	if err != nil { return err }
	return nil
}

