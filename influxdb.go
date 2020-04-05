package main

import (
	"context"
	"strconv"
	"time"

	"github.com/influxdata/influxdb-client-go"
)

const (
	countryReportsMeasurement = "country_reports"
	nlReportsMeasurement      = "nl_reports"
)

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

func (influxDB *influxDB) writeNLMetrics(ctx context.Context, field string, history map[time.Time]int) (int, error) {
	var metrics []influxdb.Metric

	for date, amount := range history {
		fields := map[string]interface{}{
			field: amount,
		}
		metrics = append(metrics, influxdb.NewRowMetric(
			fields,
			nlReportsMeasurement,
			nil,
			date,
		))
	}

	return influxDB.client.Write(ctx, influxDB.bucket, influxDB.org, metrics...)
}

func (influxDB *influxDB) writeNLHospitalizedMunicipalityMetrics(ctx context.Context, history nlHospitalizedMunicipalityHistory) (int, error) {
	var metrics []influxdb.Metric

	for date, reports := range history {
		for _, report := range reports {
			tags := map[string]string{
				"municipality":      report.Municipality,
				"municipality_code": strconv.Itoa(report.MunicipalityCode),
				"province":          report.Province,
			}
			fields := map[string]interface{}{
				"hospitalized": report.Amount,
			}
			metrics = append(metrics, influxdb.NewRowMetric(
				fields,
				nlReportsMeasurement,
				tags,
				date,
			))
		}
	}

	return influxDB.client.Write(ctx, influxDB.bucket, influxDB.org, metrics...)
}

func (influxDB *influxDB) writeNLCasesProvinceMetrics(ctx context.Context, history nlCasesProvinceHistory) (int, error) {
	var metrics []influxdb.Metric

	for date, reports := range history {
		for _, report := range reports {
			tags := map[string]string{
				"province": report.Province,
			}
			fields := map[string]interface{}{
				"confirmed": report.Amount,
			}
			metrics = append(metrics, influxdb.NewRowMetric(
				fields,
				nlReportsMeasurement,
				tags,
				date,
			))
		}
	}

	return influxDB.client.Write(ctx, influxDB.bucket, influxDB.org, metrics...)
}
