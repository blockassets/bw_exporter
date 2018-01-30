package main

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lookfirst/bw_exporter/cgminer"
	"github.com/prometheus/client_golang/prometheus"
)

//
var (
	idStateLabelNames = []string{"id", "state"}
	idRateLabelNames  = []string{"id", "rate"}
	idLabelNames      = []string{"id"}
)

//
func newGauge(metricName string, docString string, constLabels prometheus.Labels) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        metricName,
		Help:        docString,
		ConstLabels: constLabels,
	})
}

//
func newGaugeVec(metricName string, docString string, constLabels prometheus.Labels, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        metricName,
			Help:        docString,
			ConstLabels: constLabels,
		},
		labels,
	)
}

// Collector interface
type Exporter struct {
	hostname             string
	port                 int64
	timeout              time.Duration
	chipStatGauge        *prometheus.GaugeVec
	devsHashRateGauge    *prometheus.GaugeVec
	devsHashCountGauge   *prometheus.GaugeVec
	devsErrorsGauge      *prometheus.GaugeVec
	devsTemperatureGauge prometheus.Gauge
	sync.Mutex
}

//
type CgminerStats struct {
	ChipStat *cgminer.ChipStat
	Devs     *[]cgminer.Devs
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
func NewExporter(cgHost string, cgPort int64, cgTimeout time.Duration) *Exporter {
	versionLabel := prometheus.Labels{}
	version = cgminer.ReadVersionFile()
	if len(version) > 0 {
		versionLabel = prometheus.Labels{"version": version}
	}

	return &Exporter{
		hostname:             cgHost,
		port:                 cgPort,
		timeout:              cgTimeout,
		chipStatGauge:        newGaugeVec("bw_chipstat", "Chip accept/reject", versionLabel, idStateLabelNames),
		devsHashCountGauge:   newGaugeVec("bw_devs_hashcount", "Device hash accept/reject", versionLabel, idStateLabelNames),
		devsHashRateGauge:    newGaugeVec("bw_devs_hashrate", "Device hashrate", versionLabel, idRateLabelNames),
		devsErrorsGauge:      newGaugeVec("bw_devs_errors", "Device hardware errors", versionLabel, idLabelNames),
		devsTemperatureGauge: newGauge("bw_devs_temperature", "Device temperature", versionLabel),
	}
}

// Parses the map[string]float64 into a gauge with labels.
// {"1_accept": 1, "1_reject": 100}
func collectChipStat(e *Exporter, cgStats *CgminerStats) {
	chipStat := cgStats.ChipStat
	for chip, value := range *chipStat {
		s := strings.Split(chip, "_")
		chipId, chipType := s[0], s[1]
		e.chipStatGauge.WithLabelValues(chipId, chipType).Set(value)
	}
}

//
func collectDevs(e *Exporter, cgStats *CgminerStats) {
	devs := cgStats.Devs
	for id, value := range *devs {
		idStr := strconv.Itoa(id)
		e.devsHashCountGauge.WithLabelValues(idStr, "accept").Set(float64(value.Accepted))
		e.devsHashCountGauge.WithLabelValues(idStr, "reject").Set(float64(value.Rejected))
		e.devsHashRateGauge.WithLabelValues(idStr, "MHS_av").Set(value.MHSav)
		e.devsHashRateGauge.WithLabelValues(idStr, "MHS_5s").Set(value.MHS5s)
		e.devsHashRateGauge.WithLabelValues(idStr, "MHS_1m").Set(value.MHS1m)
		e.devsHashRateGauge.WithLabelValues(idStr, "MHS_5m").Set(value.MHS5m)
		e.devsHashRateGauge.WithLabelValues(idStr, "MHS_15m").Set(value.MHS15m)
		e.devsErrorsGauge.WithLabelValues(idStr).Set(float64(value.HardwareErrors))
		// All 4 devs report the same temperature
		e.devsTemperatureGauge.Set(value.Temperature)
	}
}

// Outputs the gauge values on the channel
func collectGauges(e *Exporter, ch chan<- prometheus.Metric) {
	e.chipStatGauge.Collect(ch)
	e.devsHashRateGauge.Collect(ch)
	e.devsHashCountGauge.Collect(ch)
	e.devsErrorsGauge.Collect(ch)
	e.devsTemperatureGauge.Collect(ch)
}

//
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.chipStatGauge.Describe(ch)
	e.devsHashCountGauge.Describe(ch)
	e.devsHashRateGauge.Describe(ch)
	e.devsErrorsGauge.Describe(ch)
	e.devsTemperatureGauge.Describe(ch)
}

//
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Prevents multiple concurrent calls
	e.Lock()
	defer e.Unlock()

	// Fetch data from the worker
	cgStats, err := e.fetchData()
	if err != nil {
		return
	}

	collectChipStat(e, cgStats)
	collectDevs(e, cgStats)

	collectGauges(e, ch)
}
