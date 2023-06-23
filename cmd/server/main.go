// Copyright 2019 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/hmetrics"

	"github.com/lynkdb/kvgo-server/config"
	"github.com/lynkdb/kvgo-server/data"
)

var (
	version = "git"
	release = "1"
)

func main() {

	if err := config.Setup(version, release); err != nil {
		hlog.Printf("warn", "kvgo-server/config/setup err %s", err.Error())
		hlog.Flush()
		os.Exit(1)
	}

	if err := data.Setup(); err != nil {
		hlog.Printf("warn", "kvgo-server/data/setup err %s", err.Error())
		hlog.Flush()
		os.Exit(1)
	}

	config.Flush()

	if config.Config.Server.HttpPort > 0 &&
		(config.Config.Server.PprofEnable || config.Config.Server.MetricsEnable) {

		if config.Config.Server.MetricsEnable {
			http.HandleFunc("/metrics", hmetrics.HttpHandler)
		}

		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.HttpPort))
		if err != nil {
			hlog.Printf("warn", "http listen err %s", err.Error())
			hlog.Flush()
			os.Exit(1)
		}
		go http.Serve(ln, nil)

		hlog.Printf("info", "kvgo-server http listen :%d ok",
			config.Config.Server.HttpPort)
	}

	quit := make(chan os.Signal, 2)

	//
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)
	sg := <-quit

	data.Data.Close()
	time.Sleep(1e9)

	hlog.Printf("warn", "kvgo-server signal quit %s", sg.String())
	hlog.Flush()
}
