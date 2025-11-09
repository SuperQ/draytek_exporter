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

type gaugePair struct {
	first  *prometheus.Desc
	second *prometheus.Desc
	values func(status vigorv5.Status) (float64, float64)
}

func (p gaugePair) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.first
	ch <- p.second
}

func (p gaugePair) Collect(ch chan<- prometheus.Metric, status vigorv5.Status) {
	first, second := p.values(status)
	ch <- prometheus.MustNewConstMetric(
		p.first, prometheus.GaugeValue, first,
	)
	ch <- prometheus.MustNewConstMetric(
		p.second, prometheus.GaugeValue, second,
	)
}

func streamTableGaugePair(name string, helpTpl string, values func(status vigorv5.Status) (float64, float64)) gaugePair {
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

	return gaugePair{down, up, values}
}

func streamTableFloat64GaugePair(name string, helpTpl string, values func(status vigorv5.Status) vigorv5.StreamTableValue[float64]) gaugePair {
	return streamTableGaugePair(name, helpTpl, func(status vigorv5.Status) (float64, float64) {
		v := values(status)
		return v.Downstream, v.Upstream
	})
}

func streamTableIntGaugePair(name string, helpTpl string, values func(status vigorv5.Status) vigorv5.StreamTableValue[int]) gaugePair {
	return streamTableGaugePair(name, helpTpl, func(status vigorv5.Status) (float64, float64) {
		v := values(status)
		return float64(v.Downstream), float64(v.Upstream)
	})
}

func endTableGaugePair(name string, helpTpl string, values func(status vigorv5.Status) (float64, float64)) gaugePair {
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

	return gaugePair{near, far, values}
}

func endTableFloat64GaugePair(name string, helpTpl string, values func(status vigorv5.Status) vigorv5.EndTableValue[float64]) gaugePair {
	return endTableGaugePair(name, helpTpl, func(status vigorv5.Status) (float64, float64) {
		v := values(status)
		return v.NearEnd, v.FarEnd
	})
}

func endTableIntGaugePair(name string, helpTpl string, values func(status vigorv5.Status) vigorv5.EndTableValue[int]) gaugePair {
	return endTableGaugePair(name, helpTpl, func(status vigorv5.Status) (float64, float64) {
		v := values(status)
		return float64(v.NearEnd), float64(v.FarEnd)
	})
}

func endTableBoolGaugePair(name string, helpTpl string, values func(status vigorv5.Status) vigorv5.EndTableValue[bool]) gaugePair {
	return endTableGaugePair(name, helpTpl, func(status vigorv5.Status) (float64, float64) {
		v := values(status)
		return boolToFloat(v.NearEnd), boolToFloat(v.FarEnd)
	})
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

var actualRate = streamTableIntGaugePair(
	"actual_bps",
	"The actual %s bits per second rate",
	func(status vigorv5.Status) vigorv5.StreamTableValue[int] {
		return status.ActualRate
	},
)
var attainableRate = streamTableIntGaugePair(
	"attainable_bps",
	"The attainable %s bits per second rate",
	func(status vigorv5.Status) vigorv5.StreamTableValue[int] {
		return status.AttainableRate
	},
)
var interleaveDepth = streamTableIntGaugePair(
	"interleave_depth",
	"The amount of interleaving configured for the %s",
	func(status vigorv5.Status) vigorv5.StreamTableValue[int] {
		return status.InterleaveDepth
	},
)
var actualPsd = streamTableFloat64GaugePair(
	"actual_psd_db",
	"The actual %s power spectrum density in dB",
	func(status vigorv5.Status) vigorv5.StreamTableValue[float64] {
		return status.ActualPSD
	},
)
var snrMargin = streamTableFloat64GaugePair(
	"snr_margin_db",
	"The %s SNR margin in dB",
	func(status vigorv5.Status) vigorv5.StreamTableValue[float64] {
		return status.SNRMargin
	},
)

var bitswapActive = endTableBoolGaugePair(
	"bitswap_active",
	"Whether bitswap is active on the %s",
	func(status vigorv5.Status) vigorv5.EndTableValue[bool] {
		return status.Bitswap
	},
)
var reTx = endTableBoolGaugePair(
	"retx_active",
	"Whether retransmission is active on the %s",
	func(status vigorv5.Status) vigorv5.EndTableValue[bool] {
		return status.ReTx
	},
)
var attenuation = endTableFloat64GaugePair(
	"attenuation_db",
	"The attenuation on the %s in dB",
	func(status vigorv5.Status) vigorv5.EndTableValue[float64] {
		return status.Attenuation
	},
)
var crcCount = endTableIntGaugePair(
	"crc_count",
	"The number of CRC errors on the %s",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.Crc
	},
)
var erroredSeconds = endTableIntGaugePair(
	"errored_seconds",
	"The number of seconds the %s was errored",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.Es
	},
)
var severelyErroredSeconds = endTableIntGaugePair(
	"severely_errored_seconds",
	"The number of seconds the %s was severely errored",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.Ses
	},
)
var unavailableSeconds = endTableIntGaugePair(
	"unavailable_seconds",
	"The number of seconds the %s was unavailable",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.Uas
	},
)
var hecErrorCount = endTableIntGaugePair(
	"hec_error_count",
	"The number of header errors on the %s",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.HecError
	},
)
var losFailureCount = endTableIntGaugePair(
	"los_failure_count",
	"The number of Loss of Signal failures",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.LosFailure
	},
)
var lofFailureCount = endTableIntGaugePair(
	"lof_failure_count",
	"The number of Loss of Frame failures",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.LofFailure
	},
)
var lprFailureCount = endTableIntGaugePair(
	"lpr_failure_count",
	"The number of Loss of Power failures",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.LprFailure
	},
)
var lcdFailureCount = endTableIntGaugePair(
	"lcd_failure_count",
	"The number of Loss of Cell Delineation failures",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.LcdFailure
	},
)
var rfecCount = endTableIntGaugePair(
	"rfec_count",
	"The number of Reed–Solomon Forward Error Corrections",
	func(status vigorv5.Status) vigorv5.EndTableValue[int] {
		return status.Rfec
	},
)

var gauges = []gaugePair{
	actualRate,
	attainableRate,
	interleaveDepth,
	actualPsd,
	snrMargin,

	bitswapActive,
	reTx,
	attenuation,
	crcCount,
	erroredSeconds,
	severelyErroredSeconds,
	unavailableSeconds,
	hecErrorCount,
	losFailureCount,
	lofFailureCount,
	lprFailureCount,
	lcdFailureCount,
	rfecCount,
}

// Describe describes all the metrics ever exported by the draytek_exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- draytekUpDesc
	ch <- draytekInfoDesc

	for _, gauge := range gauges {
		gauge.Describe(ch)
	}
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

	for _, gauge := range gauges {
		gauge.Collect(ch, status)
	}
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
