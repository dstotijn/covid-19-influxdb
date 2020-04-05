package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func main() {
	ctx := context.Background()

	influxDB, err := newInfluxDB(influxConfig{
		url:    mustGetenv("INFLUX_URL"),
		token:  mustGetenv("INFLUX_TOKEN"),
		bucket: mustGetenv("INFLUX_BUCKET"),
		org:    mustGetenv("INFLUX_ORG"),
	})
	if err != nil {
		log.Fatalf("[ERROR] Cannot create InfluxDB client: %v", err)
	}

	countryReports := map[string]func() (reportHistory, error){
		"confirmed": countryConfirmedCasesReports,
		"deaths":    countryDeathsReports,
		"recovered": countryRecoveredReports,
	}

	for metric, reportFn := range countryReports {
		reportHistory, err := reportFn()
		if err != nil {
			log.Fatalf("[ERROR] Getting %v country reports failed: %v", metric, err)
		}

		for date, reports := range reportHistory {
			writeCount, err := influxDB.writeCountryMetrics(ctx, metric, date, reports)
			if err != nil {
				log.Fatalf("[ERROR] Cannot write metrics to InfluxDB: %v", err)
			}

			log.Printf("[INFO] Wrote %v `%v` country metric(s) for %v.", writeCount, metric, date)
		}
	}

	nlHospitalizedHistory, err := nlHospitalizedReports()
	if err != nil {
		log.Fatalf("[ERROR] Getting NL hospitalized reports failed: %v", err)
	}
	writeCount, err := influxDB.writeNLMetrics(ctx, "hospitalized", nlHospitalizedHistory)
	if err != nil {
		log.Fatalf("[ERROR] Cannot write metrics to InfluxDB: %v", err)
	}
	log.Printf("[INFO] Wrote %v NL hospitalized metric(s).", writeCount)

	nlHospitalizedMunicipalityHistory, err := nlHospitalizedMunicipalityReports()
	if err != nil {
		log.Fatalf("[ERROR] Getting NL hospitalized municipality reports failed: %v", err)
	}
	writeCount, err = influxDB.writeNLHospitalizedMunicipalityMetrics(ctx, nlHospitalizedMunicipalityHistory)
	if err != nil {
		log.Fatalf("[ERROR] Cannot write metrics to InfluxDB: %v", err)
	}
	log.Printf("[INFO] Wrote %v NL hospitalized municipality metric(s).", writeCount)

	nlDeathsHistory, err := nlDeathsReports()
	if err != nil {
		log.Fatalf("[ERROR] Getting NL deaths reports failed: %v", err)
	}
	writeCount, err = influxDB.writeNLMetrics(ctx, "deaths", nlDeathsHistory)
	if err != nil {
		log.Fatalf("[ERROR] Cannot write metrics to InfluxDB: %v", err)
	}
	log.Printf("[INFO] Wrote %v NL deaths metric(s).", writeCount)

	nlCasesProvinceHistory, err := nlCasesProvinceReports()
	if err != nil {
		log.Fatalf("[ERROR] Getting NL cases province reports failed: %v", err)
	}
	writeCount, err = influxDB.writeNLCasesProvinceMetrics(ctx, nlCasesProvinceHistory)
	if err != nil {
		log.Fatalf("[ERROR] Cannot write metrics to InfluxDB: %v", err)
	}
	log.Printf("[INFO] Wrote %v NL cases province metric(s).", writeCount)
}

func mustGetenv(env string) string {
	token := os.Getenv(env)
	if token == "" {
		log.Fatalf("[FATAL] Environment variable `%v` is required.", env)
	}

	return token
}
