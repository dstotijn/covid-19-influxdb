package main

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go"
)

const measurementName = "reports"

type influxDB struct {
	client *influxdb.Client
	bucket string
	org    string
}

type influxConfig struct {
	url    string
	token  string
	bucket string
	org    string
}

func newInfluxDB(cfg influxConfig) (*influxDB, error) {
	client, err := influxdb.New(cfg.url, cfg.token)
	if err != nil {
		return nil, err
	}

	return &influxDB{
		client: client,
		bucket: cfg.bucket,
		org:    cfg.org,
	}, nil
}

func (influxDB *influxDB) writeMetrics(ctx context.Context, report report) (int, error) {
	var metrics []influxdb.Metric

	tags := map[string]string{
		"country":  report.Country,
		"province": report.Province,
	}

	for date, caseCount := range report.Timeline.Cases {
		ts, err := time.Parse("1/2/06", date)
		if err != nil {
			return 0, fmt.Errorf("cannot parse report date: %v", err)
		}

		fields := map[string]interface{}{
			"cases": caseCount,
		}

		if deaths, ok := report.Timeline.Deaths[date]; ok {
			fields["deaths"] = deaths
		}

		metrics = append(metrics, influxdb.NewRowMetric(
			fields,
			measurementName,
			tags,
			ts,
		))
	}

	return influxDB.client.Write(ctx, influxDB.bucket, influxDB.org, metrics...)
}
