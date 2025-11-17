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
	interleaveDepthDownDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "interleave_depth"),
		"The amount of interleaving configured for the downstream",
		nil, nil,
	)
	interleaveDepthUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "interleave_depth"),
		"The amount of interleaving configured for the upstream",
		nil, nil,
	)
	actualPsdDownDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "actual_psd_db"),
		"The actual downstream power spectrum density in dB",
		nil, nil,
	)
	actualPsdUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "actual_psd_db"),
		"The actual upstream power spectrum density in dB",
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

	bitswapActiveNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "bitswap_active"),
		"Whether bitswap is active on the near end",
		nil, nil,
	)
	bitswapActiveFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "bitswap_active"),
		"Whether bitswap is active on the far end",
		nil, nil,
	)
	reTxActiveNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "retx_active"),
		"Whether retransmission is active on the near end",
		nil, nil,
	)
	reTxActiveFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "retx_active"),
		"Whether retransmission is active on the far end",
		nil, nil,
	)
	attenuationNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "attenuation_db"),
		"The attenuation on the near end in dB",
		nil, nil,
	)
	attenuationFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "attenuation_db"),
		"The attenuation on the far end in dB",
		nil, nil,
	)
	crcCountNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "crc_errors_total"),
		"The number of CRC errors on the near end",
		nil, nil,
	)
	crcCountFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "crc_errors_total"),
		"The number of CRC errors on the far end",
		nil, nil,
	)
	erroredSecondsNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "errored_seconds_total"),
		"The number of seconds the near end was errored",
		nil, nil,
	)
	erroredSecondsFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "errored_seconds_total"),
		"The number of seconds the far end was errored",
		nil, nil,
	)
	severelyErroredSecondsNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "severely_errored_seconds_total"),
		"The number of seconds the near end was severely errored",
		nil, nil,
	)
	severelyErroredSecondsFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "severely_errored_seconds_total"),
		"The number of seconds the far end was severely errored",
		nil, nil,
	)
	unavailableSecondsNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "unavailable_seconds_total"),
		"The number of seconds the near end was unavailable",
		nil, nil,
	)
	unavailableSecondsFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "unavailable_seconds_total"),
		"The number of seconds the far end was unavailable",
		nil, nil,
	)
	hecErrorCountNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "hec_errors_total"),
		"The number of header errors on the near end",
		nil, nil,
	)
	hecErrorCountFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "hec_errors_total"),
		"The number of header errors on the far end",
		nil, nil,
	)
	losFailureCountNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "los_failures_total"),
		"The number of Loss of Signal failures at the near end",
		nil, nil,
	)
	losFailureCountFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "los_failures_total"),
		"The number of Loss of Signal failures at the far end",
		nil, nil,
	)
	lofFailureCountNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "lof_failures_total"),
		"The number of Loss of Frame failures at the near end",
		nil, nil,
	)
	lofFailureCountFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "lof_failures_total"),
		"The number of Loss of Frame failures at the far end",
		nil, nil,
	)
	lprFailureCountNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "lpr_failures_total"),
		"The number of Loss of Power failures at the near end",
		nil, nil,
	)
	lprFailureCountFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "lpr_failures_total"),
		"The number of Loss of Power failures at the far end",
		nil, nil,
	)
	lcdFailureCountNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "lcd_failures_total"),
		"The number of Loss of Cell Delineation failures at the near end",
		nil, nil,
	)
	lcdFailureCountFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "lcd_failures_total"),
		"The number of Loss of Cell Delineation failures at the far end",
		nil, nil,
	)
	rfecCountNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "reed_solomon_forward_error_corrections_total"),
		"The number of Reed–Solomon Forward Error Corrections at the near end",
		nil, nil,
	)
	rfecCountFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "reed_solomon_forward_error_corrections_total"),
		"The number of Reed–Solomon Forward Error Corrections at the far end",
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
	ch <- interleaveDepthDownDesc
	ch <- interleaveDepthUpDesc
	ch <- actualPsdDownDesc
	ch <- actualPsdUpDesc
	ch <- snrMarginDownDesc
	ch <- snrMarginUpDesc

	ch <- bitswapActiveNearEndDesc
	ch <- bitswapActiveFarEndDesc
	ch <- reTxActiveNearEndDesc
	ch <- reTxActiveFarEndDesc
	ch <- attenuationNearEndDesc
	ch <- attenuationFarEndDesc
	ch <- crcCountNearEndDesc
	ch <- crcCountFarEndDesc
	ch <- erroredSecondsNearEndDesc
	ch <- erroredSecondsFarEndDesc
	ch <- severelyErroredSecondsNearEndDesc
	ch <- severelyErroredSecondsFarEndDesc
	ch <- unavailableSecondsNearEndDesc
	ch <- unavailableSecondsFarEndDesc
	ch <- hecErrorCountNearEndDesc
	ch <- hecErrorCountFarEndDesc
	ch <- losFailureCountNearEndDesc
	ch <- losFailureCountFarEndDesc
	ch <- lofFailureCountNearEndDesc
	ch <- lofFailureCountFarEndDesc
	ch <- lprFailureCountNearEndDesc
	ch <- lprFailureCountFarEndDesc
	ch <- lcdFailureCountNearEndDesc
	ch <- lcdFailureCountFarEndDesc
	ch <- rfecCountNearEndDesc
	ch <- rfecCountFarEndDesc
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
		interleaveDepthDownDesc, prometheus.GaugeValue, float64(status.InterleaveDepthDownstream),
	)
	ch <- prometheus.MustNewConstMetric(
		interleaveDepthUpDesc, prometheus.GaugeValue, float64(status.InterleaveDepthUpstream),
	)
	ch <- prometheus.MustNewConstMetric(
		actualPsdDownDesc, prometheus.GaugeValue, status.ActualPSDDownstream,
	)
	ch <- prometheus.MustNewConstMetric(
		actualPsdUpDesc, prometheus.GaugeValue, status.ActualPSDUpstream,
	)
	ch <- prometheus.MustNewConstMetric(
		snrMarginDownDesc, prometheus.GaugeValue, status.SNRMarginDownstream,
	)
	ch <- prometheus.MustNewConstMetric(
		snrMarginUpDesc, prometheus.GaugeValue, status.SNRMarginUpstream,
	)

	ch <- prometheus.MustNewConstMetric(
		bitswapActiveNearEndDesc, prometheus.GaugeValue, optionToFloat64(status.BitswapNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		bitswapActiveFarEndDesc, prometheus.GaugeValue, optionToFloat64(status.BitswapFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		reTxActiveNearEndDesc, prometheus.GaugeValue, optionToFloat64(status.ReTxNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		reTxActiveFarEndDesc, prometheus.GaugeValue, optionToFloat64(status.ReTxFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		attenuationNearEndDesc, prometheus.GaugeValue, status.AttenuationNearEnd,
	)
	ch <- prometheus.MustNewConstMetric(
		attenuationFarEndDesc, prometheus.GaugeValue, status.AttenuationFarEnd,
	)
	ch <- prometheus.MustNewConstMetric(
		crcCountNearEndDesc, prometheus.CounterValue, float64(status.CrcNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		crcCountFarEndDesc, prometheus.CounterValue, float64(status.CrcFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		erroredSecondsNearEndDesc, prometheus.CounterValue, float64(status.EsNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		erroredSecondsFarEndDesc, prometheus.CounterValue, float64(status.EsFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		severelyErroredSecondsNearEndDesc, prometheus.CounterValue, float64(status.SesNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		severelyErroredSecondsFarEndDesc, prometheus.CounterValue, float64(status.SesFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		unavailableSecondsNearEndDesc, prometheus.CounterValue, float64(status.UasNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		unavailableSecondsFarEndDesc, prometheus.CounterValue, float64(status.UasFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		hecErrorCountNearEndDesc, prometheus.CounterValue, float64(status.HecErrorsNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		hecErrorCountFarEndDesc, prometheus.CounterValue, float64(status.HecErrorsFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		losFailureCountNearEndDesc, prometheus.CounterValue, float64(status.LosFailureNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		losFailureCountFarEndDesc, prometheus.CounterValue, float64(status.LosFailureFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		lofFailureCountNearEndDesc, prometheus.CounterValue, float64(status.LofFailureNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		lofFailureCountFarEndDesc, prometheus.CounterValue, float64(status.LofFailureFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		lprFailureCountNearEndDesc, prometheus.CounterValue, float64(status.LprFailureNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		lprFailureCountFarEndDesc, prometheus.CounterValue, float64(status.LprFailureFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		lcdFailureCountNearEndDesc, prometheus.CounterValue, float64(status.LcdFailureNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		lcdFailureCountFarEndDesc, prometheus.CounterValue, float64(status.LcdFailureFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		rfecCountNearEndDesc, prometheus.CounterValue, float64(status.RfecNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		rfecCountFarEndDesc, prometheus.CounterValue, float64(status.RfecFarEnd),
	)
}

func optionToFloat64(option bool) float64 {
	if option {
		return 1
	}
	return 0
}
