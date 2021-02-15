package main

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mkalus/fritzflux/fritzbox"
	"log"
	"os"
	"sync"
)

func main() {
	// read data from environment
	fritzUrl := os.Getenv("FRITZURL")
	fritzUser := os.Getenv("FRITZUSER")
	fritzPw := os.Getenv("FRITZPW")

	influxUrl := os.Getenv("INFLUXURL")
	influxAuth := os.Getenv("INFLUXAUTH")
	influxOrg := os.Getenv("INFLUXORG")
	influxBucket := os.Getenv("INFLUXBUCKET")

	// fallbacks
	if fritzUrl == "" {
		fritzUrl = "https://fritz.box:443"
	}
	if influxUrl == "" {
		influxUrl = "http://localhost:8086"
	}
	if influxBucket == "" {
		influxBucket = "fritz"
	}

	// create influx endpoint
	client := influxdb2.NewClient(influxUrl, influxAuth)
	writeAPI := client.WriteAPI(influxOrg, influxBucket)

	// log errors
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			log.Printf("influxdb api write error: %s", err)
		}
	}()

	var wg sync.WaitGroup

	// login into fritzbox
	c, err := fritzbox.LoginFritzbox(fritzUrl, fritzUser, fritzPw)
	if err != nil {
		log.Fatal(err)
	}

	// start logging traffic stats
	wg.Add(1)
	go func() {
		fritzbox.LogStats(c, writeAPI)
	}()

	// for some reason we have a login for the fritzbox and one for home auto...
	h, err := fritzbox.LoginHomeAuto(fritzUrl, fritzUser, fritzPw)
	if err != nil {
		log.Fatal(err)
	}

	// start logging thermostats
	wg.Add(1)
	go func() {
		fritzbox.LogThermostats(h, writeAPI)
	}()

	// wait for go routines to end
	wg.Wait()
}
