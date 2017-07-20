package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gphat/lowed"
	"github.com/gphat/lowed/ssf"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/Sirupsen/logrus"
	"github.com/go-yaml/yaml"
)

var (
	configFile = flag.String("f", "", "The config file to read for settings.")
	randSource = rand.NewSource(time.Now().UnixNano())
	rando      = rand.New(randSource)
)

func main() {
	flag.Parse()

	if configFile == nil || *configFile == "" {
		logrus.Fatal("You must specify a config file")
	}

	config, err := ReadConfig(*configFile)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to read config file")
	}

	delay, err := time.ParseDuration(config.Delay)
	if err != nil {
		logrus.WithError(err).Fatal("Cannot parse delay from config file")
	}
	logrus.WithField("delay", delay).Info("Starting generation, vroom vroom!")

	var emitter func(c lowed.Config)

	if config.Protocol == "dogstatsd" {
		statsClient, err := statsd.New(config.StatsAddress)
		if err != nil {
			logrus.WithError(err).Fatal("Unable to create stats client")
		}
		emitter = func(c lowed.Config) {
			emitDDMetric(config, statsClient)
		}
	} else if config.Protocol == "ssf" {
		conn, err := net.Dial("udp", config.StatsAddress)
		if err != nil {
			logrus.WithError(err).Fatal("Could not connect to SSF address")
		}
		emitter = func(c lowed.Config) {
			emitSSFMetric(c, conn)
		}
	} else {
		logrus.WithField("protocol", config.Protocol).Fatal("Unknown protocol")
	}

	ticker := time.NewTicker(delay)
	go func() {
		for _ = range ticker.C {
			emitter(config)
		}
	}()

	select {}
}

// ReadConfig attempts to read a config file and unmarshal it's YAML
// in to a config struct.
func ReadConfig(path string) (c lowed.Config, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	bts, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(bts, &c)
	if err != nil {
		return
	}

	return c, nil
}

func emitDDMetric(c lowed.Config, client *statsd.Client) {
	for _, service := range c.Services {
		for _, counter := range c.Metrics.Counters {
			client.Count(
				fmt.Sprintf("%s.%s", service, counter.Name), 1, nil, 1.0,
			)
		}
		for _, timer := range c.Metrics.Timers {
			client.TimeInMilliseconds(
				fmt.Sprintf("%s.%s", service, timer.Name), float64(rando.Intn(timer.Range.Max-timer.Range.Min)+timer.Range.Min), nil, 1.0,
			)
		}
		for _, histo := range c.Metrics.Histograms {
			client.Histogram(
				fmt.Sprintf("%s.%s", service, histo.Name), float64(rando.Intn(histo.Range.Max-histo.Range.Min)+histo.Range.Min), nil, 1.0,
			)
		}
		for _, gauge := range c.Metrics.Gauges {
			client.Gauge(
				fmt.Sprintf("%s.%s", service, gauge.Name), float64(rando.Intn(gauge.Range.Max-gauge.Range.Min)+gauge.Range.Min), nil, 1.0,
			)
		}
		for _, set := range c.Metrics.Sets {
			client.Set(
				fmt.Sprintf("%s.%s", service, set.Name), strconv.Itoa(rando.Intn(set.UniqueValues)), nil, 1.0,
			)
		}
	}
}

func emitSSFMetric(c lowed.Config, conn net.Conn) {
	for _, service := range c.Services {
		for _, counter := range c.Metrics.Counters {
			s := ssf.SSFSample{
				Name:   fmt.Sprintf("%s.%s", service, counter.Name),
				Value:  1,
				Metric: ssf.SSFSample_COUNTER,
			}
			sp := ssf.SSFSpan{
				Metrics: []*ssf.SSFSample{&s},
			}
			d, _ := proto.Marshal(&sp)
			conn.Write(d)
		}

		for _, timer := range c.Metrics.Timers {
			s := ssf.SSFSample{
				Name:   fmt.Sprintf("%s.%s", service, timer.Name),
				Value:  float32(rando.Intn(timer.Range.Max-timer.Range.Min) + timer.Range.Min),
				Metric: ssf.SSFSample_HISTOGRAM,
			}
			sp := ssf.SSFSpan{
				Metrics: []*ssf.SSFSample{&s},
			}
			d, _ := proto.Marshal(&sp)
			conn.Write(d)
		}

		for _, histo := range c.Metrics.Timers {
			s := ssf.SSFSample{
				Name:   fmt.Sprintf("%s.%s", service, histo.Name),
				Value:  float32(rando.Intn(histo.Range.Max-histo.Range.Min) + histo.Range.Min),
				Metric: ssf.SSFSample_HISTOGRAM,
			}
			sp := ssf.SSFSpan{
				Metrics: []*ssf.SSFSample{&s},
			}
			d, _ := proto.Marshal(&sp)
			conn.Write(d)
		}

		for _, gauge := range c.Metrics.Gauges {
			s := ssf.SSFSample{
				Name:   fmt.Sprintf("%s.%s", service, gauge.Name),
				Value:  float32(rando.Intn(gauge.Range.Max-gauge.Range.Min) + gauge.Range.Min),
				Metric: ssf.SSFSample_GAUGE,
			}
			sp := ssf.SSFSpan{
				Metrics: []*ssf.SSFSample{&s},
			}
			d, _ := proto.Marshal(&sp)
			conn.Write(d)
		}

		for _, set := range c.Metrics.Sets {
			s := ssf.SSFSample{
				Name:    fmt.Sprintf("%s.%s", service, set.Name),
				Message: strconv.Itoa(rando.Intn(set.UniqueValues)),
				Metric:  ssf.SSFSample_SET,
			}
			sp := ssf.SSFSpan{
				Metrics: []*ssf.SSFSample{&s},
			}
			d, _ := proto.Marshal(&sp)
			conn.Write(d)
		}
	}
}
