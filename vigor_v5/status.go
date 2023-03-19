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

	"github.com/go-kit/log/level"
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

	ActualRateDownstream     int
	ActualRateUpstream       int
	AttainableRateDownstream int
	AttainableRateUpstream   int
	SNRMarginDownstream      float64
	SNRMarginUpstream        float64
}

func (v *Vigor) FetchStatus() (Status, error) {
	post := vigorForm{
		pid: "0MONITORING_DSL_GENERAL",
		op:  "501",
		ct:  dslStatusGeneral,
	}

	resp, err := v.postWithLogin(post)
	if err != nil {
		level.Debug(v.logger).Log("msg", "Got error from post", "err", err)
		return Status{}, err
	}

	return v.parseDSLStatusGeneralJson(resp)
}

func (v *Vigor) parseDSLStatusGeneralJson(respJson string) (Status, error) {
	value := gjson.Get(respJson, "ct.0.0MONITORING_DSL_GENERAL.#(Name==\"Setting\")")
	if !value.Exists() {
		level.Debug(v.logger).Log("msg", "Unable to get settings", "response_json", respJson)
		return Status{}, ErrParseFailed
	}

	level.Debug(v.logger).Log("msg", "Parsed DSL Status General json", "json", value.String())

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
		case "SNR Margin":
			status.SNRMarginDownstream = parsedB(v.Get("Downstream").String())
			status.SNRMarginUpstream = parsedB(v.Get("Upstream").String())
		}
	}

	return status, nil
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
