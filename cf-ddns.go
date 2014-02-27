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
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)

	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ipChan := make(chan string)
	zChans := make([]chan string, len(config.Zones))

	// This goroutine sends a message whenever the IP has changed
	go IPChangeCheck(ipChan, time.Duration(config.CheckInterval)*time.Second)

	// Create an update channel and spawn a goroutine that handles updating, for each zone
	i := 0
	for z, zc := range config.Zones {
		zChan := make(chan string)
		zChans[i] = zChan
		i += 1

		go ZoneUpdate(zChan, z, &zc)
	}

	// Update forever
	for {
		IP := <-ipChan
		for _, zChan := range zChans {
			zChan <- IP
		}
	}
}

func ZoneUpdate(ipc chan string, z string, zc *ZoneConfig) {
	client := &CloudflareClient{
		Email: zc.Email,
		Token: zc.Token,
	}

	for {
		IP := <-ipc

		// Load all current records
		recs, err := client.RecLoadAll(z)
		if err != nil {
			log.Printf("[ERROR] Error getting records for %s: %s\n", z, err)
			continue
		}

		for _, rec := range recs {
			// Only update ipv4 records, which are not already up to date
			if rec.Type != "A" || rec.Content == IP {
				continue
			}

			// Check if the domain is in our update list
			found := false
			for _, domain := range zc.Domains {
				if rec.Name == domain {
					found = true
					break
				}
			}
			if !found {
				continue
			}

			rec.Content = IP
			err = client.RecEdit(&rec)
			if err != nil {
				log.Printf("[ERROR] Error editing domain %s: %s\n", rec.Name, err)
				continue
			}
			log.Printf("%s now points to %s\n", rec.Name, IP)
		}
	}
}
