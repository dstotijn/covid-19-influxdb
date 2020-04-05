package main

import (
	"context"
	"time"

	"github.com/influxdata/influxdb-client-go"
)

const countryReportsMeasurement = "country_reports"

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

func (influxDB *influxDB) writeCountryMetrics(ctx context.Context, metric string, date time.Time, reports []countryReport) (int, error) {
	var metrics []influxdb.Metric

	for _, report := range reports {
		tags := map[string]string{
			"country":  report.Country,
			"province": report.Province,
		}
		fields := map[string]interface{}{
			metric: report.Amount,
		}
		metrics = append(metrics, influxdb.NewRowMetric(
			fields,
			countryReportsMeasurement,
			tags,
			date,
		))
	}

	return influxDB.client.Write(ctx, influxDB.bucket, influxDB.org, metrics...)
}
