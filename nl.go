package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	nlMunicipalityUnknown     = "Unknown"
	nlMunicipalityCodeUnknown = -1
)

type nlHospitalizedMunicipalityReport struct {
	Municipality     string
	MunicipalityCode int
	Province         string
	Amount           int
}

type nlHospitalizedMunicipalityHistory map[time.Time][]nlHospitalizedMunicipalityReport

type nlCasesProvinceReport struct {
	Province string
	Amount   int
}

type nlCasesProvinceHistory map[time.Time][]nlCasesProvinceReport

func nlHospitalizedReports() (map[time.Time]int, error) {
	return nlHistory("https://raw.githubusercontent.com/J535D165/CoronaWatchNL/master/data/rivm_corona_in_nl_hosp.csv")
}

func nlDeathsReports() (map[time.Time]int, error) {
	return nlHistory("https://raw.githubusercontent.com/J535D165/CoronaWatchNL/master/data/rivm_corona_in_nl_fatalities.csv")
}

func nlHistory(url string) (map[time.Time]int, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)

	_, err = csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("cannot read header CSV record: %v", err)
	}

	history := make(map[time.Time]int)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read CSV record: %v", err)
		}

		if len(record) != 2 {
			return nil, fmt.Errorf("invalid column count: %v instead of 5", len(record))
		}

		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, fmt.Errorf("cannot parse date: %v", err)
		}

		amount, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			return nil, fmt.Errorf("cannot parse amount: %v", err)
		}

		history[date] = amount
	}

	return history, nil
}

func nlHospitalizedMunicipalityReports() (nlHospitalizedMunicipalityHistory, error) {
	resp, err := httpClient.Get("https://raw.githubusercontent.com/J535D165/CoronaWatchNL/master/data/rivm_NL_covid19_hosp_municipality.csv")
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)

	_, err = csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("cannot read header CSV record: %v", err)
	}

	history := make(nlHospitalizedMunicipalityHistory)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read CSV record: %v", err)
		}

		if len(record) != 5 {
			return nil, fmt.Errorf("invalid column count: %v instead of 5", len(record))
		}

		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, fmt.Errorf("cannot parse date: %v", err)
		}

		if _, ok := history[date]; !ok {
			history[date] = []nlHospitalizedMunicipalityReport{}
		}

		munCode, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("cannot parse municipality code: %v", err)
		}

		amount, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, fmt.Errorf("cannot parse amount: %v", err)
		}

		history[date] = append(history[date], nlHospitalizedMunicipalityReport{
			Municipality:     record[1],
			MunicipalityCode: munCode,
			Province:         record[3],
			Amount:           amount,
		})
	}

	return history, nil
}

func nlCasesProvinceReports() (nlCasesProvinceHistory, error) {
	resp, err := httpClient.Get("https://raw.githubusercontent.com/J535D165/CoronaWatchNL/master/data/rivm_NL_covid19_province.csv")
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)

	_, err = csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("cannot read header CSV record: %v", err)
	}

	history := make(nlCasesProvinceHistory)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read CSV record: %v", err)
		}

		if len(record) != 3 {
			return nil, fmt.Errorf("invalid column count: %v instead of 3", len(record))
		}

		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, fmt.Errorf("cannot parse date: %v", err)
		}

		if _, ok := history[date]; !ok {
			history[date] = []nlCasesProvinceReport{}
		}

		province := record[1]
		if province == "" {
			province = "Onbekend"
		}

		amount, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("cannot parse amount: %v", err)
		}

		history[date] = append(history[date], nlCasesProvinceReport{
			Province: province,
			Amount:   amount,
		})
	}

	return history, nil
}
