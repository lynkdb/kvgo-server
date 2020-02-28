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

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hooto/hconf4g/hconf"
	"github.com/lessos/lessgo/crypto/idhash"

	"github.com/lynkdb/kvgo"
)

var (
	version    = "0.9.0"
	release    = "1"
	AppName    = "kvgo-server"
	Prefix     = ""
	err        error
	ConfigData kvgo.Config
	confFile   = ""
)

func Setup(ver, rel string) error {

	version = ver
	release = rel

	if Prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/.."); err != nil {
		Prefix = "/opt/lynkdb/" + AppName
	}

	confFile = Prefix + "/etc/kvgo-server.conf"

	optErr := hconf.DecodeFromFile(&ConfigData, confFile)

	if os.IsNotExist(optErr) {

		fp, err := os.Open(Prefix + "/misc/conf/kvgo-server.conf.default")
		if err != nil {
			return err
		}
		defer fp.Close()

		bs, err := ioutil.ReadAll(fp)
		if err != nil {
			return err
		}

		var (
			authKey = idhash.RandBase64String(40)
			cfgStr  = string(bs)
		)
		cfgStr = strings.Replace(cfgStr,
			"change_this_server_auth_secret_key", authKey, -1)
		cfgStr = strings.Replace(cfgStr,
			"change_this_cluster_auth_secret_key", authKey, -1)

		//
		os.MkdirAll(Prefix+"/etc", 0755)
		fpo, err := os.OpenFile(confFile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		fpo.Seek(0, 0)
		fpo.Truncate(0)

		if _, err = fpo.Write([]byte(cfgStr)); err == nil {
			err = fpo.Sync()
		}
		fpo.Close()

		if err == nil {
			optErr = hconf.DecodeFromFile(&ConfigData, confFile)
		}

		optErr = err
	}

	if optErr != nil {
		return optErr
	}

	if ConfigData.Storage.DataDirectory == "" {
		ConfigData.Storage.DataDirectory = Prefix + "/var/data"
	}

	return nil
}
