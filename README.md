# FritzFlux

Scarpe Fritzbox data into InfluxDB (traffic, thermostats, &hellip;) - works with InfluxDB 1.8+ (especially with 2.x)

This service uses https://github.com/bpicode/fritzctl

## Compiling and Running

Compile and run with:

```bash
go build github.com/mkalus/fritzflux/cmd/fritzflux && ./fritzflux
```

## Environmental variables

The following variables can/must be set:

* `FRITZURL` URL to find fritz box at (default: `https://fritz.box:443`)
* `FRITZUSER` Fritzbox user (empty for default user)
* `FRITZPW` Fritzbox password
* `INFLUXURL` InfluxDB URL (default: `http://localhost:8086`)
* `INFLUXAUTH` InfluxDB auth token or user:password for InfluxDB 1.8
* `INFLUXORG` InfluxDB org name (can be empty)
* `INFLUXBUCKET` InfluxDB bucket/database (default: `fritz`)
* `INFLUXBTHERMO` InfluxDB bucket/database for thermostats (default empty, will overwrite `INFLUXBUCKET`, if set)
* `INFLUXBTRAFFIC` InfluxDB bucket/database for traffic data (default empty, will overwrite `INFLUXBUCKET`, if set)

## Scraped data

Right now:

* Thermostats data (every 15 minutes)
* Traffic data (every minute)

TODO: phone calls, power, etc. - feel free to contribute

## Docker

There is a Docker container of goggler including a headless version of Chromium. Try it using:

```bash
docker run --rm -e "FRITZPW=mysecretpw" ronix/fritzflux
```
