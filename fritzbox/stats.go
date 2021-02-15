package fritzbox

import (
	"github.com/bpicode/fritzctl/fritz"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"time"
)

// run stat logger
func LogStats(c *fritz.Client, writeAPI api.WriteAPI) {
	// wait until full minute
	now := time.Now()
	<-time.After(time.Duration(60-now.Second())*time.Second - time.Duration(now.Nanosecond()))

	err := saveStats(c, writeAPI)
	if err != nil {
		log.Printf("[error] Internet stats failed: %s", err)
	}

	// one stat per minute
	for range time.Tick(time.Minute) {
		err = saveStats(c, writeAPI)
		if err != nil {
			log.Printf("[error] Internet stats failed: %s", err)
		}
	}
}

// save stats to fritzbox
func saveStats(c *fritz.Client, writeAPI api.WriteAPI) error {
	stats, err := getLogStats(c)
	if err != nil {
		return err
	}

	// stats represent measurements within a five second interval each, so for a minute, we want to average the first 12 entries
	p := influxdb2.NewPointWithMeasurement("traffic").
		AddTag("direction", "downstream").
		AddField("internet", average(stats.DownstreamInternet, 12)).
		AddField("media", average(stats.DownStreamMedia, 12)).
		AddField("guest", average(stats.DownStreamGuest, 12))

	writeAPI.WritePoint(p)

	// stats represent measurements within a five second interval each, so for a minute, we want to average the first 12 entries
	p = influxdb2.NewPointWithMeasurement("traffic").
		AddTag("direction", "upstream").
		AddField("internet", average(stats.UpstreamRealtime, 12)).
		AddField("guest", average(stats.UpstreamGuest, 12))

	writeAPI.WritePoint(p)

	log.Print("stats saved")
	return nil
}

// averages the first count elements in the array
func average(in []float64, count int) float64 {
	var sum float64
	max := len(in)
	if max > count {
		max = count
	}

	if max == 0 {
		return 0
	}

	for i := 0; i < max; i++ {
		sum += in[i]
	}

	return sum / float64(max)
}

func getLogStats(c *fritz.Client) (*fritz.TrafficMonitoringData, error) {
	f := fritz.NewInternal(c)
	stats, err := f.InternetStats()
	if err != nil {
		return nil, err
	}

	return stats, nil
}
