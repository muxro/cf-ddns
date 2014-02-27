/*
cf-ddns - Dynamic DNS Client for CloudFlare
Copyright (C) 2014 Mux Systems SRL-D

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"

	yaml "launchpad.net/goyaml"
)

type Config struct {
	CheckInterval uint

	Zones map[string]ZoneConfig
}

type ZoneConfig struct {
	Email   string
	Token   string
	Domains []string
}

var configPath = flag.String("conf", "config.yaml", "Configuration file")

func loadConfig() (*Config, error) {
	var config Config

	cf, err := os.Open(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer cf.Close()

	cr := bufio.NewReader(cf)
	cd, err := ioutil.ReadAll(cr)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(cd, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
