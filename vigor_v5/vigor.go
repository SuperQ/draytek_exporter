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
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/tidwall/gjson"
)

var ErrJsonDecodeFailed = errors.New("json decode failed")
var ErrRequestFailed = errors.New("failed to request with login")

type Vigor struct {
	jar    *cookiejar.Jar
	client *http.Client
	cgiURL *url.URL
	csrf   string

	host     string
	username string
	password string

	logger log.Logger
}

type vigorForm struct {
	pid string
	op  string
	ct  string
}

func New(logger log.Logger, host string, username string, password string) (*Vigor, error) {
	var err error

	v := Vigor{
		host:     host,
		username: username,
		password: password,
		logger:   logger,
	}
	v.jar, err = cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	v.cgiURL, err = url.Parse(fmt.Sprintf("http://%s/cgi-bin/webproc.cgi", v.host))
	if err != nil {
		return &v, err
	}

	v.client = &http.Client{
		Jar: v.jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &v, nil
}

func (v *Vigor) postForm(p vigorForm) (*http.Response, error) {
	urlValues := url.Values{
		"pid":    {p.pid},
		"op":     {p.op},
		"ct":     {encodeVigorJson(p.ct)},
		"_token": {v.csrf},
	}

	level.Debug(v.logger).Log("msg", "Posting pid", "pid", p.pid)

	for _, cookie := range v.client.Jar.Cookies(v.cgiURL) {
		level.Debug(v.logger).Log("msg", "Post Cookie", "name", cookie.Name, "value", cookie.Value)
	}

	return v.client.PostForm(v.cgiURL.String(), urlValues)
}

func (v *Vigor) postWithLogin(p vigorForm) (string, error) {
	var rid string
	for attempts := 0; attempts < 3; attempts++ {
		resp, err := v.postForm(p)
		if err == nil && resp.StatusCode == http.StatusOK {
			respJson, err := decodeVigorJson(resp)
			if err == nil {
				rid = gjson.Get(respJson, "rid").String()
				if rid == "0000" {
					resp.Body.Close()
					return respJson, nil
				}
			}
		}
		defer resp.Body.Close()
		level.Debug(v.logger).Log("msg", "Post failed, attempting login", "status", resp.Status, "err", err, "rid", rid)
		err = v.Login()
		if err != nil {
			level.Debug(v.logger).Log("msg", "Login failed", "err", err)
		}
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	return "", ErrRequestFailed
}

func decodeVigorJson(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ErrJsonDecodeFailed
	}
	respPadding, err := strconv.Atoi(string(body[0]))
	if err != nil {
		return "", ErrJsonDecodeFailed
	}
	if respPadding > 2 {
		return "", ErrJsonDecodeFailed
	}
	body = body[1:]
	for i := 0; i < respPadding; i++ {
		body = append(body, []byte("=")...)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return "", ErrJsonDecodeFailed
	}

	return string(decoded), nil
}

func encodeVigorJson(j string) string {
	j = base64.StdEncoding.EncodeToString([]byte(j))

	paddingLength := len(j)
	j, _ = strings.CutSuffix(j, "=")
	padding := strconv.Itoa(paddingLength - len(j))

	return padding + j
}
