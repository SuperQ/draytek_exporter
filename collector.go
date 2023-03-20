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
	"github.com/SuperQ/draytek_exporter/vigor_v5"
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

	actualRateDownDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "actual_bps"),
		"The actual downstream bits per second rate",
		nil, nil,
	)
	actualRateUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "actual_bps"),
		"The actual upstream bits per second rate",
		nil, nil,
	)
	attainableRateDownDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "attainable_bps"),
		"The attainable downstream bits per second rate",
		nil, nil,
	)
	attainableRateUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "attainable_bps"),
		"The attainable upstream bits per second rate",
		nil, nil,
	)
	snrMarginDownDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "snr_margin_db"),
		"The downstream SNR margin in dB",
		nil, nil,
	)
	snrMarginUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "snr_margin_db"),
		"The downstream SNR margin in dB",
		nil, nil,
	)
)

// Describe describes all the metrics ever exported by the draytek_exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- draytekUpDesc
	ch <- draytekInfoDesc
	ch <- actualRateDownDesc
	ch <- actualRateUpDesc
	ch <- attainableRateDownDesc
	ch <- attainableRateUpDesc
	ch <- snrMarginDownDesc
	ch <- snrMarginUpDesc
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
	ch <- prometheus.MustNewConstMetric(
		actualRateDownDesc, prometheus.GaugeValue, float64(status.ActualRateDownstream),
	)
	ch <- prometheus.MustNewConstMetric(
		actualRateUpDesc, prometheus.GaugeValue, float64(status.ActualRateUpstream),
	)
	ch <- prometheus.MustNewConstMetric(
		attainableRateDownDesc, prometheus.GaugeValue, float64(status.AttainableRateDownstream),
	)
	ch <- prometheus.MustNewConstMetric(
		attainableRateUpDesc, prometheus.GaugeValue, float64(status.AttainableRateUpstream),
	)
	ch <- prometheus.MustNewConstMetric(
		snrMarginDownDesc, prometheus.GaugeValue, status.SNRMarginDownstream,
	)
	ch <- prometheus.MustNewConstMetric(
		snrMarginUpDesc, prometheus.GaugeValue, status.SNRMarginUpstream,
	)
}
