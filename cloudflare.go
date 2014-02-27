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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const cfCloudflareUrl = "https://www.cloudflare.com/api_json.html"

type CloudflareClient struct {
	Email string
	Token string
}

type CloudflareRecObj struct {
	Zone        string `json:"zone_name"`
	ID          string `json:"rec_id"`
	Name        string
	Type        string
	Content     string
	TTL         string
	ServiceMode string `json:"service_mode"`
}

type CloudflareRecs struct {
	Has_more bool
	Count    int
	Objs     []CloudflareRecObj
}

type CloudflareResponse struct {
	Recs CloudflareRecs
}

type CloudflareRecLoadAll struct {
	Request  map[string]string
	Response CloudflareResponse
	Result   string
	Msg      string
}

func (c *CloudflareClient) RecLoadAll(z string) ([]CloudflareRecObj, error) {
	params := make(map[string]string)
	params["z"] = z

	var m CloudflareRecLoadAll
	err := c.apiCall("rec_load_all", params, &m)
	if err != nil {
		return nil, err
	}
	if m.Result != "success" {
		return nil, fmt.Errorf("API call returned error: %v", m.Msg)
	}

	return m.Response.Recs.Objs, nil
}

func (c *CloudflareClient) RecEdit(r *CloudflareRecObj) error {
	params := make(map[string]string)
	params["z"] = r.Zone
	params["type"] = r.Type
	params["id"] = r.ID
	params["name"] = r.Name
	params["content"] = r.Content
	params["ttl"] = r.TTL
	if r.Type == "A" {
		params["service_mode"] = r.ServiceMode
	}

	var m CloudflareRecLoadAll
	err := c.apiCall("rec_edit", params, &m)
	if err != nil {
		return err
	}
	if m.Result != "success" {
		return fmt.Errorf("API call returned error: %v", m.Msg)
	}

	return nil
}

func (c *CloudflareClient) apiCall(a string, args map[string]string, v interface{}) error {
	p := url.Values{}
	p.Set("email", c.Email)
	p.Set("tkn", c.Token)
	p.Set("a", a)

	for arg, val := range args {
		p.Set(arg, val)
	}

	resp, err := http.PostForm(cfCloudflareUrl, p)
	if err != nil {
		return fmt.Errorf("Error posting request: %v", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(v)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v", err)
	}

	return nil

}
