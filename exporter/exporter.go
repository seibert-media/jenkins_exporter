package exporter

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	// namespace defines the Prometheus namespace for this exporter.
	namespace = "jenkins"
)

var (
	// isUp defines if the API response can get processed.
	isUp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Check if Jenkins response can be processed",
		},
	)

	// jobSuccess defines a map to collect the build color codes.
	jobSuccess = map[string]prometheus.Gauge{}

	// jobFail defines a map to collect the build color codes.
	jobFail = map[string]prometheus.Gauge{}
	// jobWeather defines a map to collect the build color codes.
	jobWeather = map[string]prometheus.Gauge{}
)

// init just defines the initial state of the exports.
func init() {
	isUp.Set(0)
}

// NewExporter gives you a new exporter instance.
func NewExporter(address, username, password string) *Exporter {
	return &Exporter{
		address:  address,
		username: username,
		password: password,
	}
}

// Exporter combines the metric collector and descritions.
type Exporter struct {
	address  string
	username string
	password string
	mutex    sync.RWMutex
}

// Describe defines the metric descriptions for Prometheus.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- isUp.Desc()

	for _, metric := range jobSuccess {
		ch <- metric.Desc()
	}
	for _, metric := range jobFail {
		ch <- metric.Desc()
	}
	for _, metric := range jobWeather {
		ch <- metric.Desc()
	}
}

// Collect delivers the metrics to Prometheus.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if err := e.scrape(); err != nil {
		log.Error(err)

		isUp.Set(0)
		ch <- isUp

		return
	}

	ch <- isUp

	for _, metric := range jobSuccess {
		ch <- metric
	}
	for _, metric := range jobFail {
		ch <- metric
	}
	for _, metric := range jobWeather {
		ch <- metric
	}
}

// scrape just starts the scraping loop.
func (e *Exporter) scrape() error {
	log.Debug("start scrape loop")

	var (
		root = &Root{}
	)

	if err := root.Fetch(e.address, e.username, e.password); err != nil {
		log.Debugf("%s", err)
		return fmt.Errorf("failed to fetch root data")
	}

	for _, job := range *root {
		log.Debugf("processing %s job", job.Name)
		jobName := job.Name
		//jobName := strings.Replace(strings.ToLower(job.Name), " ", "_", -1)
		if _, ok := jobSuccess[job.Key()]; ok == false {
			jobSuccess[job.Key()] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      "job_success",
					Help:      "number of successful branches for job",
					ConstLabels: prometheus.Labels{
						"name": jobName,
					},
				},
			)
		}
		if _, ok := jobFail[job.Key()]; ok == false {
			jobFail[job.Key()] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      "job_fail",
					Help:      "number of failed branches for job",
					ConstLabels: prometheus.Labels{
						"name": jobName,
					},
				},
			)
		}
		if _, ok := jobWeather[job.Key()]; ok == false {
			jobWeather[job.Key()] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      "job_weather",
					Help:      "weather for job",
					ConstLabels: prometheus.Labels{
						"name": jobName,
					},
				},
			)
		}

		log.Debugf("setting success to %d for %s", job.Success, job.Name)
		log.Debugf("setting failed to %d for %s", job.Fail, job.Name)
		log.Debugf("setting weather to %d for %s", job.Weather, job.Name)

		jobSuccess[job.Key()].Set(float64(job.Success))
		jobFail[job.Key()].Set(float64(job.Fail))
		jobWeather[job.Key()].Set(float64(job.Weather))
	}

	isUp.Set(1)
	return nil
}
