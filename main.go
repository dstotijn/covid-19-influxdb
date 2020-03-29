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

	reports, err := reports()
	if err != nil {
		log.Fatalf("[ERROR] Getting country reports failed: %v", err)
	}

	for _, report := range reports {
		writeCount, err := influxDB.writeMetrics(ctx, report)
		if err != nil {
			log.Fatalf("[ERROR] Cannot write metrics to InfluxDB: %v", err)
		}

		log.Printf("[INFO] Wrote %v metric(s) (%+v, %+v).", writeCount, report.Country, report.Province)
	}
}

func mustGetenv(env string) string {
	token := os.Getenv(env)
	if token == "" {
		log.Fatalf("[FATAL] Environment variable `%v` is required.", env)
	}

	return token
}
