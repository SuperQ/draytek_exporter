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
	vigorUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the vigor instance status successful?",
		nil, nil,
	)
	vigorInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "info"),
		"Info about the vigor router",
		[]string{"dsl_version", "mode", "profile"}, nil,
	)
)

// Describe describes all the metrics ever exported by the vigor_exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- vigorUpDesc
	ch <- vigorInfoDesc
}

// Collect fetches the stats from the vigor router and delivers them as
// Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	status, err := e.v.FetchStatus()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			vigorUpDesc, prometheus.GaugeValue, 0.0,
		)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		vigorUpDesc, prometheus.GaugeValue, 1.0,
	)

	ch <- prometheus.MustNewConstMetric(
		vigorInfoDesc, prometheus.GaugeValue, 1.0,
		status.DSLVersion, status.Mode, status.Profile,
	)
}
