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
// This package provides a library to acces DrayTek Vigor "v5" firmware.
//
// # Supported Devices
//
// * DrayTek Vigor 167
//
// # Notes
//
// ## The Vigor v5 API has some abnormal encoding.
//
// * The request is sent as a POST with a base64 encoded json.
// * The API always returns 200 OK
// * The response json is also base64 encoded.
// * The base64 encoding uses prefix notation for padding.
// * The actual response code is in the returned json payload as `{"rid":"0000"}`
//
// ## Prefix padding base64 encoding.
//
// The base64 encoding sent and received from the firmware does not use the
// standard trailing `=` padding. Instead, it places a number at the beginning
// to denote how much encoding padding is needed.
//
// Example:
// * Original json: `{"foo": "bar"}`
// * Normal base64: `e2ZvbzogYmFyfQo=`
// * DrayTek base64: `1e2ZvbzogYmFyfQo`

package vigorv5
