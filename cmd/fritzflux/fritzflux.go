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
	influxBucketThermo := os.Getenv("INFLUXBTHERMO")
	influxBucketTraffic := os.Getenv("INFLUXBTRAFFIC")

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
	if influxBucketThermo == "" {
		influxBucketThermo = influxBucket
	}
	if influxBucketTraffic == "" {
		influxBucketTraffic = influxBucket
	}

	// create influx endpoint
	client := influxdb2.NewClient(influxUrl, influxAuth)
	writeAPIThermo := client.WriteAPI(influxOrg, influxBucketThermo)
	writeAPITraffic := client.WriteAPI(influxOrg, influxBucketTraffic)

	// log errors
	errorsCh1 := writeAPIThermo.Errors()
	errorsCh2 := writeAPITraffic.Errors()
	go func() {
		for err := range errorsCh1 {
			log.Printf("influxdb thermo api write error: %s", err)
		}
		for err := range errorsCh2 {
			log.Printf("influxdb traffic api write error: %s", err)
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
		fritzbox.LogStats(c, writeAPITraffic)
	}()

	// for some reason we have a login for the fritzbox and one for home auto...
	h, err := fritzbox.LoginHomeAuto(fritzUrl, fritzUser, fritzPw)
	if err != nil {
		log.Fatal(err)
	}

	// start logging thermostats
	wg.Add(1)
	go func() {
		fritzbox.LogThermostats(h, writeAPIThermo)
	}()

	// wait for go routines to end
	wg.Wait()
}
