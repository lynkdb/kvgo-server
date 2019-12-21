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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hooto/hini4g/hini"
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
)

func Setup(ver, rel string) error {

	version = ver
	release = rel

	if Prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/.."); err != nil {
		Prefix = "/opt/lynkdb/" + AppName
	}

	file := Prefix + "/etc/kvgo-server.conf"
	opts, optErr := hini.ParseFile(file)

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
		fpo, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
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
			opts, err = hini.ParseFile(file)
		}

		optErr = err
	}

	if optErr != nil {
		return optErr
	}

	return dataSetup(opts)
}

func dataSetup(opts *hini.Options) error {

	{
		if v, ok := opts.ValueOK("storage/data_directory"); ok {
			ConfigData.StorageDataDirectory = v.String()
		} else {
			ConfigData.StorageDataDirectory = Prefix + "/var/data"
		}
	}

	{
		if v, ok := opts.ValueOK("server/bind"); ok {
			ConfigData.ServerBind = v.String()
		} else {
			return errors.New("no server/bind found")
		}

		if v, ok := opts.ValueOK("server/auth_secret_key"); ok {
			ConfigData.ServerAuthSecretKey = v.String()
		}
	}

	{
		if v, ok := opts.ValueOK("performance/write_buffer_size"); ok {
			ConfigData.PerformanceWriteBufferSize = v.Int()
		}

		if v, ok := opts.ValueOK("performance/block_cache_size"); ok {
			ConfigData.PerformanceBlockCacheSize = v.Int()
		}

		if v, ok := opts.ValueOK("performance/max_table_size"); ok {
			ConfigData.PerformanceMaxTableSize = v.Int()
		}

		if v, ok := opts.ValueOK("performance/max_open_files"); ok {
			ConfigData.PerformanceMaxOpenFiles = v.Int()
		}
	}

	{
		if v, ok := opts.ValueOK("cluster/masters"); ok {
			ConfigData.ClusterMasters = strings.Split(v.String(), ",")
		}

		if v, ok := opts.ValueOK("cluster/auth_secret_key"); ok {
			ConfigData.ClusterAuthSecretKey = v.String()
		}
	}

	return nil
}
