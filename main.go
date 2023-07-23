package main

import (
	"encoding/json"
	"log"
	"os"
	"tunProxy/client"
)

type Config struct {
	ServerAddr string
	ClientAddr string
	StaticKey  string
}

func main() {
	cfg := &Config{}
	fd, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("cannot open config file : %v", err)
	}
	jb := make([]byte, 1024)
	n, err1 := fd.Read(jb)
	fd.Close()
	if err1 != nil {
		log.Fatalln(err1)
	}

	err2 := json.Unmarshal(jb[:n], cfg)
	if err2 != nil {
		log.Fatalln(err2)
	}

	tunClient := client.NewTunClient(cfg.ServerAddr, cfg.ClientAddr, cfg.StaticKey)
	tunClient.Launch()
	tunClient.Accepting()
}
