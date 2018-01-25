package main

import (
	"log"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
	"time"
	"bw_exporter/cgminer"
        "strconv"
)

//
func newGaugeMetric(metricName string, docString string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricName,
			Help: docString,
		},
		labels,
	)
}

// Collector interface
type Exporter struct {
	hostname			string
	port				int64
	timeout				time.Duration
	chipStatGauge		*prometheus.GaugeVec
	devsHashRateGauge	*prometheus.GaugeVec
	devsHashCountGauge	*prometheus.GaugeVec
	devsErrorsGauge		*prometheus.GaugeVec
	sync.Mutex
}

//
type CgminerStats struct {
	ChipStat	*cgminer.ChipStat
	Devs		*[]cgminer.Devs
}

//
func (e *Exporter) fetchData() (*CgminerStats, error) {
	log.Printf("fetchData(%s:%d, %s)", e.hostname, e.port, e.timeout)

	miner := cgminer.New(e.hostname, e.port, e.timeout)

	chipStat, err := miner.ChipStat()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	devs, err := miner.Devs()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &CgminerStats{chipStat, devs}, nil
}

//
func NewExporter(cgHost string, cgPort int64, cgTimeout time.Duration) (*Exporter, error) {
	return &Exporter{
		hostname:			cgHost,
		port:				cgPort,
		timeout:			cgTimeout,
		chipStatGauge:		newGaugeMetric("chip_total", "Chip accept/reject total", []string{"id", "state"}),
		devsHashCountGauge:	newGaugeMetric("devs_hashcount_total", "Device hash accept/reject total", []string{"id", "state"}),
		devsHashRateGauge:	newGaugeMetric("devs_hashrate_total", "Device hashrate total", []string{"id", "rate"}),
		devsErrorsGauge:	newGaugeMetric("devs_errors_total", "Device hardware errors total", []string{"id"}),
	}, nil
}

// Parses the map[string]float64 into a gauge with labels.
// {"1_accept": 1, "1_reject": 100}
func collectChipStat(e *Exporter, cgStats *CgminerStats, ch chan<- prometheus.Metric) {
	chipStat := cgStats.ChipStat
	for chip, value := range *chipStat {
		s := strings.Split(chip, "_")
		chipId, chipType := s[0], s[1]
		e.chipStatGauge.WithLabelValues(chipId, chipType).Set(value)
	}

	e.chipStatGauge.Collect(ch)
}

//
func collectDevs(e *Exporter, cgStats *CgminerStats, ch chan<- prometheus.Metric) {
	devs := cgStats.Devs
	for id, value := range *devs {
		idStr := strconv.Itoa(id)
		e.devsHashCountGauge.WithLabelValues(idStr, "accept").Set(float64(value.Accepted))
		e.devsHashCountGauge.WithLabelValues(idStr, "reject").Set(float64(value.Rejected))
		e.devsHashRateGauge.WithLabelValues(idStr, "MHS_5s").Set(value.MHS5s)
		e.devsHashRateGauge.WithLabelValues(idStr, "MHS_1m").Set(value.MHS1m)
		e.devsHashRateGauge.WithLabelValues(idStr, "MHS_5m").Set(value.MHS5m)
		e.devsHashRateGauge.WithLabelValues(idStr, "MHS_15m").Set(value.MHS15m)
		e.devsErrorsGauge.WithLabelValues(idStr).Set(float64(value.HardwareErrors))
	}

	e.devsHashRateGauge.Collect(ch)
	e.devsHashCountGauge.Collect(ch)
	e.devsErrorsGauge.Collect(ch)
}

//
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.chipStatGauge.Describe(ch)
	e.devsHashRateGauge.Describe(ch)
	e.devsHashCountGauge.Describe(ch)
	e.devsErrorsGauge.Describe(ch)
}

//
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.Lock()
	defer e.Unlock()

	// Fetch data from the worker
	cgStats, err := e.fetchData()
	if err != nil {
		return
	}

	collectChipStat(e, cgStats, ch)
	collectDevs(e, cgStats, ch)
}
