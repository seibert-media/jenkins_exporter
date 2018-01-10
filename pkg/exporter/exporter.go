package exporter

import (
	"strconv"
	"sync"

	"github.com/pkg/errors"

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
	queueLength = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "queueLength",
			Help:      "Amount of items in queue",
		},
	)
	stuckBuilds = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "stuckBuilds",
			Help:      "Amount of items stuck in queue",
		},
	)

	jobSuccess       = map[string]prometheus.Gauge{}
	jobFail          = map[string]prometheus.Gauge{}
	jobWeather       = map[string]prometheus.Gauge{}
	jobLegacyWeather = map[string]prometheus.Gauge{}
	jobInQueue       = map[string]prometheus.Gauge{}
	jobColor         = map[string]map[string]prometheus.Gauge{}
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
	ch <- queueLength.Desc()
	ch <- stuckBuilds.Desc()

	for _, metric := range jobSuccess {
		ch <- metric.Desc()
	}
	for _, metric := range jobFail {
		ch <- metric.Desc()
	}
	for _, metric := range jobWeather {
		ch <- metric.Desc()
	}
	for _, metric := range jobLegacyWeather {
		ch <- metric.Desc()
	}
	for _, metric := range jobInQueue {
		ch <- metric.Desc()
	}
	for _, metric := range jobColor {
		for _, metric := range metric {
			ch <- metric.Desc()
		}
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
	ch <- queueLength
	ch <- stuckBuilds

	for _, metric := range jobSuccess {
		ch <- metric
	}
	for _, metric := range jobFail {
		ch <- metric
	}
	for _, metric := range jobWeather {
		ch <- metric
	}
	for _, metric := range jobLegacyWeather {
		ch <- metric
	}
	for _, metric := range jobInQueue {
		ch <- metric
	}
	for _, metric := range jobColor {
		for _, metric := range metric {
			ch <- metric
		}
	}
}

// scrape just starts the scraping loop.
func (e *Exporter) scrape() error {
	log.Debug("start scrape loop")

	var (
		blueRoot   = &BlueRoot{}
		legacyRoot = &LegacyRoot{}
		queueRoot  = &QueueRoot{}
	)

	if err := blueRoot.Fetch(e.address, e.username, e.password); err != nil {
		log.Debugf("%s", err)
		return errors.Wrap(err, "failed to fetch blueRoot data")
	}
	if err := legacyRoot.Fetch(e.address, e.username, e.password); err != nil {
		log.Debugf("%s", err)
		return errors.Wrap(err, "failed to fetch legacyRoot data")
	}
	if err := queueRoot.Fetch(e.address, e.username, e.password); err != nil {
		log.Debugf("%s", err)
		return errors.Wrap(err, "failed to fetch queueRoot data")
	}

	queueLength.Set(float64(len(queueRoot.Items)))
	stuck := 0
	for _, item := range queueRoot.Items {
		if item.Stuck {
			stuck = stuck + 1
		}
		jobID := strconv.Itoa(item.ID)
		if _, ok := jobInQueue[jobID]; ok == false {
			jobInQueue[jobID] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      "job_in_queue",
					Help:      "amount of time a job spends in queue",
					ConstLabels: prometheus.Labels{
						"id": jobID,
					},
				},
			)
		}
		log.Debugf("setting time in queue to %d for %s", item.InQueueSince, item.ID)
		jobInQueue[jobID].Set(float64(item.InQueueSince))
	}
	log.Debugf("setting stuck jobs to %d", stuck)
	stuckBuilds.Set(float64(stuck))

	for _, job := range *blueRoot {
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

	for _, job := range legacyRoot.Jobs {
		log.Debugf("processing %s job", job.Name)
		jobName := job.Name
		//jobName := strings.Replace(strings.ToLower(job.Name), " ", "_", -1)

		if _, ok := jobLegacyWeather[job.Key()]; ok == false {
			jobLegacyWeather[job.Key()] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      "job_legacy_weather",
					Help:      "legacy weather for job",
					ConstLabels: prometheus.Labels{
						"name": jobName,
					},
				},
			)
		}
		if len(job.Health) > 0 {
			health := job.Health[0].Score
			log.Debugf("setting legacy weather to %d for %s", health, job.Name)
			jobLegacyWeather[job.Key()].Set(float64(health))
		} else {
			log.Debugf("setting legacy weather to %d for %s", 100, job.Name)
			jobLegacyWeather[job.Key()].Set(float64(100))
		}

		if _, ok := jobColor[job.Key()]; ok == false {
			jobColor[job.Key()] = make(map[string]prometheus.Gauge)
		}

		for _, branch := range job.Jobs {
			if _, ok := jobColor[job.Key()][branch.Name]; ok == false {
				jobColor[job.Key()][branch.Name] = prometheus.NewGauge(
					prometheus.GaugeOpts{
						Namespace: namespace,
						Name:      "job_color",
						Help:      "color for job",
						ConstLabels: prometheus.Labels{
							"name":   jobName,
							"branch": branch.Name,
						},
					},
				)
			}
			log.Debugf("setting color to %d for %s.%s", colorToGauge(branch.Color), job.Name, branch.Name)
			jobColor[job.Key()][branch.Name].Set(float64(colorToGauge(branch.Color)))
		}
	}

	isUp.Set(1)
	return nil
}
