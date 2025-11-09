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

type StreamTableValue[T any] struct {
	Downstream T
	Upstream   T
}

type EndTableValue[T any] struct {
	NearEnd T
	FarEnd  T
}

type Status struct {
	Status     string
	Mode       string
	Profile    string
	Annex      string
	DSLVersion string

	ActualRate      StreamTableValue[int]
	AttainableRate  StreamTableValue[int]
	InterleaveDepth StreamTableValue[int]
	ActualPSD       StreamTableValue[float64]
	SNRMargin       StreamTableValue[float64]

	Bitswap     EndTableValue[bool]
	ReTx        EndTableValue[bool]
	Attenuation EndTableValue[float64]
	Crc         EndTableValue[int]
	Es          EndTableValue[int]
	Ses         EndTableValue[int]
	Uas         EndTableValue[int]
	HecError    EndTableValue[int]
	LosFailure  EndTableValue[int]
	LofFailure  EndTableValue[int]
	LprFailure  EndTableValue[int]
	LcdFailure  EndTableValue[int]
	Rfec        EndTableValue[int]
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

func parseStreamTableValue[T any](value *StreamTableValue[T], json gjson.Result, parser func(string) T) {
	value.Downstream = parser(json.Get("Downstream").String())
	value.Upstream = parser(json.Get("Upstream").String())
}

func parseEndTableValue[T any](value *EndTableValue[T], json gjson.Result, parser func(string) T) {
	value.NearEnd = parser(json.Get("Near_End").String())
	value.FarEnd = parser(json.Get("Far_End").String())
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
			parseStreamTableValue(&status.ActualRate, v, parseKbps)
		case "Attainable Rate":
			parseStreamTableValue(&status.AttainableRate, v, parseKbps)
		case "Interleave Depth":
			parseStreamTableValue(&status.InterleaveDepth, v, parseCount)
		case "Actual PSD":
			parseStreamTableValue(&status.ActualPSD, v, parsedB)
		case "SNR Margin":
			parseStreamTableValue(&status.SNRMargin, v, parsedB)
		}
	}

	endTable := value.Get("End_Table").Array()
	for _, v := range endTable {
		switch v.Get("Name").String() {
		case "Bitswap":
			parseEndTableValue(&status.Bitswap, v, parseOption)
		case "ReTx":
			parseEndTableValue(&status.ReTx, v, parseOption)
		case "Attenuation":
			parseEndTableValue(&status.Attenuation, v, parsedB)
		case "CRC":
			parseEndTableValue(&status.Crc, v, parseCount)
		case "ES":
			parseEndTableValue(&status.Es, v, parseSeconds)
		case "SES":
			parseEndTableValue(&status.Ses, v, parseSeconds)
		case "UAS":
			parseEndTableValue(&status.Uas, v, parseSeconds)
		case "HEC Errors":
			parseEndTableValue(&status.HecError, v, parseCount)
		case "LOS Failure":
			parseEndTableValue(&status.LosFailure, v, parseCount)
		case "LOF Failure":
			parseEndTableValue(&status.LofFailure, v, parseCount)
		case "LPR Failure":
			parseEndTableValue(&status.LprFailure, v, parseCount)
		case "LCD Failure":
			parseEndTableValue(&status.LcdFailure, v, parseCount)
		case "RFEC":
			parseEndTableValue(&status.Rfec, v, parseCount)
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
