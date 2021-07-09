package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	URL      = "https://api.covid19api.com/country/turkey"
	fileName = "turkey_covid.csv"
)

type CovidDataUnit struct {
	Country   string
	Confirmed int
	Deaths    int
	Recovered int
	Active    int
	Date      string
}

func (c *CovidDataUnit) getHeaders() string {
	return strings.Join([]string{"Country", "Confirmed", "Deaths", "Recovered", "Active", "Date"}, ",")
}

func (c *CovidDataUnit) getDataArray() string {

	return strings.Join(
		[]string{c.Country,
			strconv.Itoa(c.Confirmed),
			strconv.Itoa(c.Deaths),
			strconv.Itoa(c.Recovered),
			strconv.Itoa(c.Active),
			c.Date,
		}, ",")
}

type CovidData []CovidDataUnit

func FetchCovidData() error {
	resp, err := http.Get(URL)
	if err != nil {
		return fmt.Errorf("ERROR while fetching %s\n", URL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status code in %s : %d", URL, resp.StatusCode)
	}

	var covidData CovidData

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&covidData)
	if err != nil {
		return fmt.Errorf("error while decoding json : %v\n", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error while creating %s : %v\n", fileName, err)
	}

	defer file.Close()

	for i, line := range covidData {
		if i == 0 {
			_, err := fmt.Fprintf(file, "%s\n", line.getHeaders())
			if err != nil {
				return fmt.Errorf("error while writing data headers :%v", err)
			}
		}
		if i != len(covidData)-1 {
			_, err := fmt.Fprintf(file, "%s\n", line.getDataArray())
			if err != nil {
				return fmt.Errorf("error while writing data array :%v", err)
			}
		}
	}

	fmt.Printf("Write to %s successfully\n", fileName)
	return nil
}
