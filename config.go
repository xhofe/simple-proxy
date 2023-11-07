package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Config struct {
	Proxies map[string]string `json:"proxies"`
}

var config Config

func init() {
	confPath := flag.String("conf", "./config.json", "config file path")
	flag.Parse()
	f, err := os.OpenFile(*confPath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}
