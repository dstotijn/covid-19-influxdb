package main

import (
	"context"
	"log"
	"os"
)

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
}

func mustGetenv(env string) string {
	token := os.Getenv(env)
	if token == "" {
		log.Fatalf("[FATAL] Environment variable `%v` is required.", env)
	}

	return token
}
