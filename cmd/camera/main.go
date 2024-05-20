package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nenavizhuleto/genconf"
)

var ConfigPath = "config.json"

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("failed to get machine hostname:", err)
	}
	log.Println("hostname:", hostname)

	ConfigPath = fmt.Sprintf("%s.json", hostname)

	config := DefaultConfig()
	if err := genconf.NewJSON(ConfigPath).Load(&config); err != nil {
		log.Fatalln("failed to initialize config:", err)
	}

	log.Println("config initialized at", ConfigPath)
	log.Printf("%s", config)

	app, err := NewApp(config)
	if err != nil {
		log.Fatalln(err)
	}

	go ServeSadm(app.camera, config.Sadm.Port)

	retry := NewRetry(config.Camera.Maintenance.MaxRetryCount, time.Duration(config.Camera.Maintenance.MaxTimeoutSec)*time.Second)
	retry.Do(app.Loop)
}
