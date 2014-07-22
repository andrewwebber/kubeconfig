// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v1"
)

var configFile string

type Config struct {
	DNS      string `yaml:"dns"`
	Gateway  string `yaml:"gateway"`
	MasterIP string `yaml:"master_ip"`
	Node1IP  string `yaml:"node1_ip"`
	Node2IP  string `yaml:"node2_ip"`
	SSHKey   string `yaml:"sshkey"`
	Token    string `yaml:"token"`
}


func init() {
	flag.StringVar(&configFile, "c", "kubernetes.yml", "config file to use")
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
	data["token"] = c.Token
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
	render(data)

	// Generate node2.yml
	data["subnet"] = "10.244.2.1/24"
	data["hostname"] = "node2"
	data["ip"] = c.Node2IP
	render(data)
}

func render(data map[string]string) {
	f, err := os.Create(data["hostname"] + ".yml")
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := nodeTmpl.Execute(f, data); err != nil {
		log.Fatal(err.Error())
	}
}
