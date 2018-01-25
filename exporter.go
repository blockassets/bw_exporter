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

func newGaugeMetric(metricName string, docString string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricName,
			Help: docString,
		},
		labels,
	)
}

// Define the internal structure of the Exporter class
type Exporter struct {
	hostname      string
	port          int64
	timeout       time.Duration
	chipStatGauge *prometheus.GaugeVec
	devsHashrateGauge  *prometheus.GaugeVec
	devsAcceptedGauge  *prometheus.GaugeVec
	devsRejectedGauge  *prometheus.GaugeVec
	devsErrorsGauge    *prometheus.GaugeVec
	sync.Mutex
}

type CgminerStats struct {
	ChipStat *cgminer.ChipStat
	Devs *[]cgminer.Devs
}

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

func NewExporter(cgHost string, cgPort int64, cgTimeout time.Duration) (*Exporter, error) {
	return &Exporter{
		hostname:      cgHost,
		port:          cgPort,
		timeout:       cgTimeout,
		chipStatGauge: newGaugeMetric("chip_total", "How many accept/reject on a chip", []string{"id", "state"}),
		devsHashrateGauge:     newGaugeMetric("devs_hashrate", "Device Hashrate", []string{"id"}),
		devsAcceptedGauge:     newGaugeMetric("devs_accepted", "Device Hash Accepted", []string{"id"}),
		devsRejectedGauge:     newGaugeMetric("devs_rejected", "Device Hash Rejected", []string{"id"}),
		devsErrorsGauge:       newGaugeMetric("devs_errors", "Device Hardware Errors", []string{"id"}),
	}, nil
}

// This function is used to provide comments in the HTTP output
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.chipStatGauge.Describe(ch)
	e.devsHashrateGauge.Describe(ch)
	e.devsAcceptedGauge.Describe(ch)
	e.devsRejectedGauge.Describe(ch)
	e.devsErrorsGauge.Describe(ch)
}

// This is the main function that gets called on a HTTP request, to produce the collected
// output data.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.Lock()
	defer e.Unlock()

        // cgStats is the combined structure read from the device by the fetchData() function
        // It has a ChipStat and a Devs element.
	cgStats, err := e.fetchData()
	if err != nil {
		return
	}
	
	// Read the data from the ChipStat element and parse into the chipStatGauge array.
	chipStat := cgStats.ChipStat
	for chip, value := range *chipStat {
		s := strings.Split(chip, "_")
		chipId, chipType := s[0], s[1]
		e.chipStatGauge.WithLabelValues(chipId, chipType).Set(value)
	}
	
	// Read the data from the Devs element and parse into 4 separate arrays for
	// hashrate, accepted rate, rejected rate and hardware errors.
	devs := cgStats.Devs
	for id, value := range *devs {
		e.devsHashrateGauge.WithLabelValues(strconv.Itoa(id)).Set(value.MHS5s)
		e.devsAcceptedGauge.WithLabelValues(strconv.Itoa(id)).Set(float64(value.Accepted))
		e.devsRejectedGauge.WithLabelValues(strconv.Itoa(id)).Set(float64(value.Rejected))
		e.devsErrorsGauge.WithLabelValues(strconv.Itoa(id)).Set(float64(value.HardwareErrors))
		// log.Printf("dev %d, %3.2f", id, value.MHS5s)
	}

	// This causes the gauges to be emitted on the HTTP output.
	e.chipStatGauge.Collect(ch)
	e.devsHashrateGauge.Collect(ch)
	e.devsAcceptedGauge.Collect(ch)
	e.devsRejectedGauge.Collect(ch)
	e.devsErrorsGauge.Collect(ch)
}

