package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
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
	Date      time.Time
}

func (c *CovidDataUnit) getHeaders() []string {
	return []string{"Country", "Confirmed", "Deaths", "Recovered", "Active", "Date"}
}

func (c *CovidDataUnit) getDataArray() []string {
	return []string{c.Country,
		strconv.Itoa(c.Confirmed),
		strconv.Itoa(c.Deaths),
		strconv.Itoa(c.Recovered),
		strconv.Itoa(c.Active),
		c.Date.String(),
	}
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

	csvWriter := csv.NewWriter(file)

	for i, line := range covidData {
		if i == 0 {
			err := csvWriter.Write(line.getHeaders())
			if err != nil {
				return fmt.Errorf("error while writing data headers :%v\n", err)
			}
		}
		err := csvWriter.Write(line.getDataArray())
		if err != nil {
			return fmt.Errorf("error while writing data array :%v\n", err)
		}
	}
	fmt.Printf("Write to %s successfully\n", fileName)
	return nil
}
