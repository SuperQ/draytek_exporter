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
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/tidwall/gjson"
)

var ErrLoginFailed = errors.New("login failed")

const (
	loginJSONTemplate = `{"param":[],"ct":[{"Name":"%s","Password":"%s","locales":"en"}]}`
)

func (v *Vigor) Login() error {
	// Rotate the login token.
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		v.logger.Error("Unable to generate new CSRF token", "err", err)
		return err
	}
	v.csrf = hex.EncodeToString(token)

	v.logger.Debug("Attempting login", "username", v.username)
	post := vigorForm{
		pid: "event",
		op:  "552",
		ct:  encodeLogin(v.username, v.password),
	}
	resp, err := v.postForm(post)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		v.logger.Debug("Server returned non-ok http status", "status", resp.Status)
		return ErrLoginFailed
	}

	respJSON, err := decodeVigorJSON(resp)
	if err != nil {
		v.logger.Debug("Decoding response failed", "err", err)
		return ErrLoginFailed
	}

	rid := gjson.Get(respJSON, "rid").String()
	if rid != "0000" {
		v.logger.Debug("Got invalid response ID", "rid", rid)
		return ErrLoginFailed
	}

	cookies := resp.Header.Get("Set-Cookie")
	if cookies == "" {
		return ErrLoginFailed
	}

	for _, cookie := range v.jar.Cookies(v.cgiURL) {
		v.logger.Debug("Got Cookie", "name", cookie.Name, "value", cookie.Value)
	}

	v.logger.Debug("Login OK")

	return nil
}

func encodeLogin(username string, password string) string {
	h := sha512.New()
	h.Write([]byte(password))
	encodedPassword := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf(loginJSONTemplate, username, encodedPassword)
}
