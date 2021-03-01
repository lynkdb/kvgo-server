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
	"os"
	"path/filepath"

	"github.com/hooto/htoml4g/htoml"

	"github.com/lynkdb/kvgo"
)

var (
	version  = "0.9.2"
	release  = "1"
	AppName  = "kvgo-server"
	Prefix   = ""
	confFile = ""
	Config   kvgo.Config
	err      error
)

func Setup(ver, rel string) error {

	version = ver
	release = rel

	if Prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/.."); err != nil {
		Prefix = "/opt/lynkdb/" + AppName
	}

	confFile = Prefix + "/etc/kvgo-server.conf"

	err = htoml.DecodeFromFile(&Config, confFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if Config.Server.Bind == "" {
		Config.Server.Bind = "127.0.0.1:9200"
	}

	if Config.Storage.DataDirectory == "" {
		Config.Storage.DataDirectory = Prefix + "/var/data"
	}

	Config.Reset()

	return Flush()
}

func Flush() error {
	return htoml.EncodeToFile(Config, confFile, nil)
}
