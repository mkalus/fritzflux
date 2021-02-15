package fritzbox

import (
	"github.com/bpicode/fritzctl/fritz"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"strconv"
	"time"
)

// run thermostat logger
func LogThermostats(h fritz.HomeAuto, writeAPI api.WriteAPI) {
	// wait until full minute divisible by 15
	now := time.Now()
	<-time.After(time.Duration(15-now.Minute()%15)*time.Minute - time.Duration(now.Second())*time.Second - time.Duration(now.Nanosecond()))

	err := saveThermostats(h, writeAPI)
	if err != nil {
		log.Printf("[error] Thermostats failed: %s", err)
	}

	// one stat per 15 minutes
	for range time.Tick(time.Minute * 15) {
		err = saveThermostats(h, writeAPI)
		if err != nil {
			log.Printf("[error] Thermostats failed: %s", err)
		}
	}
}

// save thermostat data to db
func saveThermostats(h fritz.HomeAuto, writeAPI api.WriteAPI) error {
	list, err := h.List()
	if err != nil {
		return err
	}

	// devices
	for _, item := range list.Devices {
		// only those that are relevant for us
		if item.IsThermostat() || item.CanMeasureTemp() {
			p := influxdb2.NewPointWithMeasurement("thermostat").
				AddTag("ain", item.Identifier).
				AddTag("name", item.Name).
				AddTag("productname", item.Productname).
				AddField("version", item.Fwversion).
				AddField("present", item.Present == 1) // convert to bool

			// can measure temperature?
			if item.CanMeasureTemp() {
				measured, _ := strconv.ParseFloat(item.Temperature.FmtCelsius(), 32)

				p.
					AddTag("has_temperature", "1").
					AddField("measured", measured)
			}

			// is thermostat?
			if item.IsThermostat() {
				// we have to convert a lot of stuff to numeric values...
				lock := 0
				if item.Thermostat.Lock == "1" {
					lock = 1
				}
				deviceLock := 0
				if item.Thermostat.DeviceLock == "1" {
					deviceLock = 1
				}

				offset, _ := strconv.ParseFloat(item.Temperature.FmtOffset(), 32)
				want, _ := strconv.ParseFloat(item.Thermostat.FmtGoalTemperature(), 32)
				saving, _ := strconv.ParseFloat(item.Thermostat.FmtSavingTemperature(), 32)
				comfort, _ := strconv.ParseFloat(item.Thermostat.FmtComfortTemperature(), 32)

				windowOpen := 0
				if item.Thermostat.WindowOpen == "1" {
					windowOpen = 1
				}

				batteryChargeLevel, _ := strconv.ParseInt(item.Thermostat.BatteryChargeLevel, 10, 8)

				batteryLow := 0
				if item.Thermostat.BatteryLow == "1" {
					batteryLow = 1
				}

				p.AddField("lock", lock).
					AddTag("is_thermostat", "1").
					AddField("devicelock", deviceLock).
					AddField("offset", offset).
					AddField("want", want).
					AddField("saving", saving).
					AddField("comfort", comfort).
					AddField("windowopen", windowOpen).
					AddField("state", item.Thermostat.ErrorCode). // string ok
					AddField("battery", batteryChargeLevel).
					AddField("batterylow", batteryLow)
			}

			p.SetTime(time.Now())

			writeAPI.WritePoint(p)
		}
	}

	log.Print("thermostats saved")
	return nil
}
