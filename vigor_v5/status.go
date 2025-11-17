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
package vigorv5

import (
	"errors"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

var ErrUpdateFailed = errors.New("dsl status update failed")
var ErrParseFailed = errors.New("dsl status parse failed")

const (
	dslStatusGeneral = `{"param":[],"ct":[{"0MONITORING_DSL_GENERAL":[]},{"1MON_DSL_STREAM_TABLE":[]},{"1MON_DSL_END_TABLE":[]}]}`
)

type Status struct {
	Status     string
	Mode       string
	Profile    string
	Annex      string
	DSLVersion string

	ActualRateDownstream      int
	ActualRateUpstream        int
	AttainableRateDownstream  int
	AttainableRateUpstream    int
	InterleaveDepthDownstream int
	InterleaveDepthUpstream   int
	ActualPSDDownstream       float64
	ActualPSDUpstream         float64
	SNRMarginDownstream       float64
	SNRMarginUpstream         float64

	BitswapNearEnd     bool
	BitswapFarEnd      bool
	ReTxNearEnd        bool
	ReTxFarEnd         bool
	AttenuationNearEnd float64
	AttenuationFarEnd  float64
	CrcNearEnd         int
	CrcFarEnd          int
	EsNearEnd          int
	EsFarEnd           int
	SesNearEnd         int
	SesFarEnd          int
	UasNearEnd         int
	UasFarEnd          int
	HecErrorsNearEnd   int
	HecErrorsFarEnd    int
	LosFailureNearEnd  int
	LosFailureFarEnd   int
	LofFailureNearEnd  int
	LofFailureFarEnd   int
	LprFailureNearEnd  int
	LprFailureFarEnd   int
	LcdFailureNearEnd  int
	LcdFailureFarEnd   int
	RfecNearEnd        int
	RfecFarEnd         int
}

func (v *Vigor) FetchStatus() (Status, error) {
	post := vigorForm{
		pid: "0MONITORING_DSL_GENERAL",
		op:  "501",
		ct:  dslStatusGeneral,
	}

	resp, err := v.postWithLogin(post)
	if err != nil {
		v.logger.Debug("Got error from post", "err", err)
		return Status{}, err
	}

	return v.parseDSLStatusGeneralJSON(resp)
}

func (v *Vigor) parseDSLStatusGeneralJSON(respJSON string) (Status, error) {
	value := gjson.Get(respJSON, "ct.0.0MONITORING_DSL_GENERAL.#(Name==\"Setting\")")
	if !value.Exists() {
		v.logger.Debug("Unable to get settings", "response_json", respJSON)
		return Status{}, ErrParseFailed
	}

	v.logger.Debug("Parsed DSL Status General json", "json", value.String())

	status := Status{
		Status:     value.Get("Status").String(),
		Mode:       value.Get("Mode").String(),
		Profile:    value.Get("Profile").String(),
		Annex:      value.Get("Annex").String(),
		DSLVersion: value.Get("DSL_Version").String(),
	}

	streamTable := value.Get("Stream_Table").Array()
	for _, v := range streamTable {
		switch v.Get("Name").String() {
		case "Actual Rate":
			status.ActualRateDownstream = parseKbps(v.Get("Downstream").String())
			status.ActualRateUpstream = parseKbps(v.Get("Upstream").String())
		case "Attainable Rate":
			status.AttainableRateDownstream = parseKbps(v.Get("Downstream").String())
			status.AttainableRateUpstream = parseKbps(v.Get("Upstream").String())
		case "Interleave Depth":
			status.InterleaveDepthDownstream = parseKbps(v.Get("Downstream").String())
			status.InterleaveDepthUpstream = parseKbps(v.Get("Upstream").String())
		case "Actual PSD":
			status.ActualPSDDownstream = parsedB(v.Get("Downstream").String())
			status.ActualPSDUpstream = parsedB(v.Get("Upstream").String())
		case "SNR Margin":
			status.SNRMarginDownstream = parsedB(v.Get("Downstream").String())
			status.SNRMarginUpstream = parsedB(v.Get("Upstream").String())
		}
	}

	endTable := value.Get("End_Table").Array()
	for _, v := range endTable {
		switch v.Get("Name").String() {
		case "Bitswap":
			status.BitswapNearEnd = parseOption(v.Get("Near_End").String())
			status.BitswapFarEnd = parseOption(v.Get("Far_End").String())
		case "ReTx":
			status.ReTxNearEnd = parseOption(v.Get("Near_End").String())
			status.ReTxFarEnd = parseOption(v.Get("Far_End").String())
		case "Attenuation":
			status.AttenuationNearEnd = parsedB(v.Get("Near_End").String())
			status.AttenuationFarEnd = parsedB(v.Get("Far_End").String())
		case "CRC":
			status.CrcNearEnd = parseCount(v.Get("Near_End").String())
			status.CrcFarEnd = parseCount(v.Get("Far_End").String())
		case "ES":
			status.EsNearEnd = parseSeconds(v.Get("Near_End").String())
			status.EsFarEnd = parseSeconds(v.Get("Far_End").String())
		case "SES":
			status.SesNearEnd = parseSeconds(v.Get("Near_End").String())
			status.SesFarEnd = parseSeconds(v.Get("Far_End").String())
		case "UAS":
			status.UasNearEnd = parseSeconds(v.Get("Near_End").String())
			status.UasFarEnd = parseSeconds(v.Get("Far_End").String())
		case "HEC Errors":
			status.HecErrorsNearEnd = parseCount(v.Get("Near_End").String())
			status.HecErrorsFarEnd = parseCount(v.Get("Far_End").String())
		case "LOS Failure":
			status.LosFailureNearEnd = parseCount(v.Get("Near_End").String())
			status.LosFailureFarEnd = parseCount(v.Get("Far_End").String())
		case "LOF Failure":
			status.LofFailureNearEnd = parseCount(v.Get("Near_End").String())
			status.LofFailureFarEnd = parseCount(v.Get("Far_End").String())
		case "LPR Failure":
			status.LprFailureNearEnd = parseCount(v.Get("Near_End").String())
			status.LprFailureFarEnd = parseCount(v.Get("Far_End").String())
		case "LCD Failure":
			status.LcdFailureNearEnd = parseCount(v.Get("Near_End").String())
			status.LcdFailureFarEnd = parseCount(v.Get("Far_End").String())
		case "RFEC":
			status.RfecNearEnd = parseCount(v.Get("Near_End").String())
			status.RfecFarEnd = parseCount(v.Get("Far_End").String())
		}
	}

	return status, nil
}

func parseCount(s string) int {
	count, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return count
}

func parseOption(s string) bool {
	return s == "ON"
}

func parseSeconds(s string) int {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0
	}
	if parts[1] != "s" {
		return 0
	}
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return x
}

func parseKbps(s string) int {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0
	}
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return x * 1000
}

func parsedB(s string) float64 {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0
	}
	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}
	return x
}
