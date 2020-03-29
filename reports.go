package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// See: https://docs.corona.lmao-xd.wtf/version-2#get-historical-data.
type report struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	Timeline struct {
		Cases  map[string]int `json:"cases"`
		Deaths map[string]int `json:"deaths"`
	} `json:"timeline"`
}

func reports() ([]report, error) {
	resp, err := httpClient.Get("https://corona.lmao.ninja/v2/historical")
	if err != nil {
		return nil, fmt.Errorf("cannot execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var reports []report
	err = json.NewDecoder(resp.Body).Decode(&reports)
	if err != nil {
		return nil, fmt.Errorf("cannot parse HTTP response body: %v", err)
	}

	return reports, nil
}
