// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v1"
)

var (
	configFile string
	iso        bool
	isoGenHost string
)

type Config struct {
	DNS      string `yaml:"dns"`
	Gateway  string `yaml:"gateway"`
	MasterIP string `yaml:"master_ip"`
	Node1IP  string `yaml:"node1_ip"`
	Node2IP  string `yaml:"node2_ip"`
	SSHKey   string `yaml:"sshkey"`
}

func init() {
	flag.StringVar(&configFile, "c", "kubernetes.yml", "config file to use")
	flag.BoolVar(&iso, "iso", false, "generate config-drive iso images")
	flag.StringVar(&isoGenHost, "isogenhost", "", "Host address for the iso generator")
}

func main() {
	flag.Parse()
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	var c Config
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	data := make(map[string]string)
	data["dns"] = c.DNS
	data["gateway"] = c.Gateway
	data["sshkey"] = c.SSHKey
	data["machines"] = fmt.Sprintf("%s,%s,%s", c.MasterIP, c.Node1IP, c.Node2IP)

	// Generate master.yml
	data["subnet"] = "10.244.0.1/24"
	data["hostname"] = "master"
	data["ip"] = c.MasterIP
	render(data)

	// Generate node1.yml
	data["subnet"] = "10.244.1.1/24"
	data["hostname"] = "node1"
	data["ip"] = c.Node1IP
	data["peers"] = fmt.Sprintf("%s:7001,%s:7001", c.MasterIP, c.Node2IP)
	render(data)

	// Generate node2.yml
	data["subnet"] = "10.244.2.1/24"
	data["hostname"] = "node2"
	data["ip"] = c.Node2IP
	data["peers"] = fmt.Sprintf("%s:7001,%s:7001", c.MasterIP, c.Node1IP)
	render(data)
}

func render(data map[string]string) {
	var buf bytes.Buffer
	f, err := os.Create(data["hostname"] + ".yml")
	if err != nil {
		log.Fatal(err.Error())
	}
	w := io.MultiWriter(f, &buf)
	if err := nodeTmpl.Execute(w, data); err != nil {
		log.Fatal(err.Error())
	}
	if iso {
		isoName := data["hostname"] + ".iso"
		url := "http://107.178.217.37:6500/genisoimage"
		if len(isoGenHost) > 0 {
			url = fmt.Sprintf("http://%s:6500/genisoimage", isoGenHost)
		}

		resp, err := http.Post(url, "application/yaml", &buf)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Fatal("non 200 exit code")
		}
		f, err := os.Create(isoName)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer f.Close()
		io.Copy(f, resp.Body)
	}
}
