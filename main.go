package main

import "log"

func main() {
	err := FetchCovidData()
	if err != nil {
		log.Fatal(err)
	}
}
