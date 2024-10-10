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
	"net/http"
	_ "net/http/pprof"
	"os"

	vigorv5 "github.com/SuperQ/draytek_exporter/vigor_v5"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const exporterName = "draytek_exporter"

func init() {
	prometheus.MustRegister(versioncollector.NewCollector(exporterName))
}

func main() {
	var (
		toolkitFlags = kingpinflag.AddFlags(kingpin.CommandLine, ":9103")
		metricsPath  = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()

		username    = kingpin.Flag("username", "username to authenticate to the target").Default("monitor").String()
		passwordEnv = kingpin.Flag("password-env", "Env var that contains password to authenticate to the target").Default("DRAYTEK_PASSWORD").String()
		target      = kingpin.Flag("target", "target host/ip the router/modem is reachable on").Default("192.168.1.1").String()
	)
	promslogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.Version(version.Print("mysqld_exporter"))
	kingpin.Parse()
	logger := promslog.New(promslogConfig)

	logger.Info("Starting "+exporterName, "version", version.Info())
	logger.Info("Build context", "build_context", version.BuildContext())

	password := os.Getenv(*passwordEnv)
	if password == "" {
		logger.Error("Missing password from env", "env", *passwordEnv)
		os.Exit(1)
	}

	var err error
	v, err := vigorv5.New(logger, *target, *username, password)
	if err != nil {
		logger.Error("Unable to create target", "err", err)
		os.Exit(1)
	}

	err = v.Login()
	if err != nil {
		logger.Error("Failed initial login attempt", "err", err)
		os.Exit(1)
	}
	logger.Info("Initial Login on DrayTek device successful")

	http.Handle(*metricsPath, promhttp.Handler())
	if *metricsPath != "/" && *metricsPath != "" {
		landingConfig := web.LandingConfig{
			Name:        "DrayTek Exporter",
			Description: "Prometheus Exporter for DrayTek modems/routers",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}

	prometheus.MustRegister(NewExporter(v))

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
		logger.Error("Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
