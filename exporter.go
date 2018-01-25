package main

import (
	"log"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
	"time"
	"bw_exporter/cgminer"
)

func newGaugeMetric(metricName string, docString string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricName,
			Help: docString,
		},
		labels,
	)
}

type Exporter struct {
	hostname      string
	port          int64
	timeout       time.Duration
	chipStatGuage *prometheus.GaugeVec
	sync.Mutex
}

type CgminerStats struct {
	ChipStat *cgminer.ChipStat
}

func (e *Exporter) fetchData() (*CgminerStats, error) {
	log.Printf("fetchData(%s:%d, %s)", e.hostname, e.port, e.timeout)

	miner := cgminer.New(e.hostname, e.port, e.timeout)

	chipStat, err := miner.ChipStat()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &CgminerStats{chipStat}, nil
}

func NewExporter(cgHost string, cgPort int64, cgTimeout time.Duration) (*Exporter, error) {
	return &Exporter{
		hostname:      cgHost,
		port:          cgPort,
		timeout:       cgTimeout,
		chipStatGuage: newGaugeMetric("chip_total", "How many accept/reject on a chip", []string{"id", "state"}),
	}, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.chipStatGuage.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.Lock()
	defer e.Unlock()

	cgStats, err := e.fetchData()
	if err != nil {
		return
	}
	chipStat := cgStats.ChipStat

	for chip, value := range *chipStat {
		s := strings.Split(chip, "_")
		chipId, chipType := s[0], s[1]

		e.chipStatGuage.WithLabelValues(chipId, chipType).Set(value)
	}

	e.chipStatGuage.Collect(ch)
}

