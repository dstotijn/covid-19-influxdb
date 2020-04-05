package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type countryReport struct {
	Country  string
	Province string
	Amount   int
}

type reportHistory map[time.Time][]countryReport

const (
	provinceCol = iota
	countryCol
)

func countryConfirmedCasesReports() (reportHistory, error) {
	return countryReports("https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv")
}

func countryDeathsReports() (reportHistory, error) {
	return countryReports("https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv")
}

func countryRecoveredReports() (reportHistory, error) {
	return countryReports("https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_recovered_global.csv")
}

func countryReports(url string) (reportHistory, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)

	headerRecord, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("cannot read header CSV record: %v", err)
	}
	if len(headerRecord) < 5 {
		return nil, fmt.Errorf("unexpected column length (%v)", len(headerRecord))
	}

	dates := make([]time.Time, 0, len(headerRecord)-4)
	for _, date := range headerRecord[4:] {
		ts, err := time.Parse("1/2/06", date)
		if err != nil {
			return nil, fmt.Errorf("cannot parse date: %v", err)
		}
		dates = append(dates, ts)
	}

	reportHist := make(reportHistory, len(dates))

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read CSV record: %v", err)
		}

		if len(record) < 5 || len(record)-4 != len(dates) {
			return nil, fmt.Errorf("unexpected column length (%v)", len(record))
		}

		country := record[countryCol]
		province := record[provinceCol]

		for i, amountString := range record[4:] {
			amount, err := strconv.Atoi(amountString)
			if err != nil {
				return nil, fmt.Errorf("cannot parse amount: %v", err)
			}
			reportHist[dates[i]] = append(reportHist[dates[i]], countryReport{
				Country:  country,
				Province: province,
				Amount:   amount,
			})
		}
	}

	return reportHist, nil
}
