// Copyright 2018 Ben Kochie <superq@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"

	vigorv5 "github.com/SuperQ/draytek_exporter/vigor_v5"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "draytek"

// Exporter collects Vigor stats from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	v *vigorv5.Vigor
}

// NewExporter returns an initialized Exporter.
func NewExporter(v *vigorv5.Vigor) *Exporter {
	return &Exporter{v: v}
}

type metricPair struct {
	first  *prometheus.Desc
	second *prometheus.Desc
}

func (p metricPair) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.first
	ch <- p.second
}

func (p metricPair) Collect(ch chan<- prometheus.Metric, first float64, second float64) {
	ch <- prometheus.MustNewConstMetric(
		p.first, prometheus.GaugeValue, first,
	)
	ch <- prometheus.MustNewConstMetric(
		p.second, prometheus.GaugeValue, second,
	)
}

func streamTableMetricsPair(name string, helpTpl string) metricPair {
	down := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", name),
		fmt.Sprintf(helpTpl, "downstream"),
		nil, nil,
	)
	up := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", name),
		fmt.Sprintf(helpTpl, "upstream"),
		nil, nil,
	)

	return metricPair{down, up}
}

func endTableMetricsPair(name string, helpTpl string) metricPair {
	near := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", name),
		fmt.Sprintf(helpTpl, "near end"),
		nil, nil,
	)
	far := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", name),
		fmt.Sprintf(helpTpl, "far end"),
		nil, nil,
	)

	return metricPair{near, far}
}

var (
	draytekUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the draytek instance status successful?",
		nil, nil,
	)
	draytekInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "info"),
		"Info about the draytek router",
		[]string{"dsl_version", "mode", "profile"}, nil,
	)
)
var actualRate = streamTableMetricsPair("actual_bps", "The actual %s bits per second rate")
var attainableRate = streamTableMetricsPair("attainable_bps", "The attainable %s bits per second rate")
var interleaveDepth = streamTableMetricsPair("interleave_depth", "The amount of interleaving configured for the %s")
var actualPsd = streamTableMetricsPair("actual_psd_db", "The actual %s power spectrum density in dB")
var snrMargin = streamTableMetricsPair("snr_margin_db", "The %s SNR margin in dB")

var bitswapActive = endTableMetricsPair("bitswap_active", "Whether bitswap is active on the %s")
var reTx = endTableMetricsPair("retx_active", "Whether retransmission is active on the %s")
var attenuation = endTableMetricsPair("attenuation_db", "The attenuation on the %s in dB")
var crcCount = endTableMetricsPair("crc_count", "The number of CRC errors on the %s")
var erroredSeconds = endTableMetricsPair("errored_seconds", "The number of seconds the %s was errored")
var severelyErroredSeconds = endTableMetricsPair("severely_errored_seconds", "The number of seconds the %s was severely errored")
var unavailableSeconds = endTableMetricsPair("unavailable_seconds", "The number of seconds the %s was unavailable")
var hecErrorCount = endTableMetricsPair("hec_error_count", "The number of header errors on the %s")
var losFailureCount = endTableMetricsPair("los_failure_count", "The number of Loss of Signal failures")
var lofFailureCount = endTableMetricsPair("lof_failure_count", "The number of Loss of Frame failures")
var lprFailureCount = endTableMetricsPair("lpr_failure_count", "The number of Loss of Power failures")
var lcdFailureCount = endTableMetricsPair("lcd_failure_count", "The number of Loss of Cell Delineation failures")
var rfecCount = endTableMetricsPair("rfec_count", "The number of Reed–Solomon Forward Error Corrections")

// Describe describes all the metrics ever exported by the draytek_exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- draytekUpDesc
	ch <- draytekInfoDesc

	actualRate.Describe(ch)
	attainableRate.Describe(ch)
	interleaveDepth.Describe(ch)
	actualPsd.Describe(ch)
	snrMargin.Describe(ch)

	bitswapActive.Describe(ch)
	reTx.Describe(ch)
	attenuation.Describe(ch)
	crcCount.Describe(ch)
	erroredSeconds.Describe(ch)
	severelyErroredSeconds.Describe(ch)
	unavailableSeconds.Describe(ch)
	hecErrorCount.Describe(ch)
	losFailureCount.Describe(ch)
	lofFailureCount.Describe(ch)
	lprFailureCount.Describe(ch)
	lcdFailureCount.Describe(ch)
	rfecCount.Describe(ch)
}

// Collect fetches the stats from the draytek router and delivers them as
// Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	status, err := e.v.FetchStatus()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			draytekUpDesc, prometheus.GaugeValue, 0.0,
		)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		draytekUpDesc, prometheus.GaugeValue, 1.0,
	)

	ch <- prometheus.MustNewConstMetric(
		draytekInfoDesc, prometheus.GaugeValue, 1.0,
		status.DSLVersion, status.Mode, status.Profile,
	)

	actualRate.Collect(ch, float64(status.ActualRate.Downstream), float64(status.ActualRate.Upstream))
	attainableRate.Collect(ch, float64(status.AttainableRate.Downstream), float64(status.AttainableRate.Upstream))
	interleaveDepth.Collect(ch, float64(status.InterleaveDepth.Downstream), float64(status.InterleaveDepth.Upstream))
	actualPsd.Collect(ch, status.ActualPSD.Downstream, status.ActualPSD.Upstream)
	snrMargin.Collect(ch, status.SNRMargin.Downstream, status.SNRMargin.Upstream)

	bitswapActive.Collect(ch, boolToFloat(status.Bitswap.NearEnd), boolToFloat(status.Bitswap.FarEnd))
	reTx.Collect(ch, boolToFloat(status.ReTx.NearEnd), boolToFloat(status.ReTx.FarEnd))
	attenuation.Collect(ch, float64(status.Attenuation.NearEnd), float64(status.Attenuation.FarEnd))
	crcCount.Collect(ch, float64(status.Crc.NearEnd), float64(status.Crc.FarEnd))
	erroredSeconds.Collect(ch, float64(status.Es.NearEnd), float64(status.Es.FarEnd))
	severelyErroredSeconds.Collect(ch, float64(status.Ses.NearEnd), float64(status.Ses.FarEnd))
	unavailableSeconds.Collect(ch, float64(status.Uas.NearEnd), float64(status.Uas.FarEnd))
	hecErrorCount.Collect(ch, float64(status.HecError.NearEnd), float64(status.HecError.FarEnd))
	losFailureCount.Collect(ch, float64(status.LosFailure.NearEnd), float64(status.LosFailure.FarEnd))
	lofFailureCount.Collect(ch, float64(status.LofFailure.NearEnd), float64(status.LofFailure.FarEnd))
	lprFailureCount.Collect(ch, float64(status.LprFailure.NearEnd), float64(status.LprFailure.FarEnd))
	lcdFailureCount.Collect(ch, float64(status.LcdFailure.NearEnd), float64(status.LcdFailure.FarEnd))
	rfecCount.Collect(ch, float64(status.Rfec.NearEnd), float64(status.Rfec.FarEnd))
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
