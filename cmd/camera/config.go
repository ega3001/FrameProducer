package main

import (
	"encoding/json"
)

type CameraConfig struct {
	// Name of the camera
	Name string `json:"name"`

	// VideoSource: RTSP | File | DeviceID
	Source string `json:"source"`

	// Camera Location
	Location string `json:"location"`

	// Time threshold between frames to dismiss motion event (in ms)
	Threshold string `json:"threshold"`

	// Max duration time of motion event (in ms)
	MaxDuration string `json:"max_duration"`

	// Maintenance
	Maintenance struct {
		MaxRetryCount int `json:"max_retry_count"`
		MaxTimeoutSec int `json:"max_timeout_sec"`
	} `json:"maintenance"`
}

type KafkaConfig struct {
	// Address of kafka instance
	Broker string `json:"broker"`

	// Name of topic to which frames should be sent
	Frame string `json:"frame"`

	// Name of topic to which motion events should be sent
	Detection string `json:"detection"`
}

type ZookeeperConfig struct {
	// Address of zookeeper
	Url string `json:"url"`

	// Group to listen
	GroupId string `json:"groupId"`
}

type SadmConfig struct {
	Port string `json:"port"`
}

type Config struct {
	// Namespace for configuring kafka
	Kafka KafkaConfig `json:"kafka"`

	// Namespace for configuring camera
	Camera CameraConfig `json:"camera"`

	// Namespace for configuring sadm
	Sadm SadmConfig `json:"sadm"`

	// Namespace for congifuring zookeeper
	Zookeeper ZookeeperConfig `json:"zookeeper"`
}

func (c Config) String() string {
	v, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic("unreachable")
	}
	return string(v)
}

func DefaultConfig() Config {
	camera := CameraConfig{
		Name:        "camera",
		Threshold:   "1s",
		MaxDuration: "5s",
	}

	// Zero for infinity
	camera.Maintenance.MaxRetryCount = 0
	camera.Maintenance.MaxTimeoutSec = 10

	return Config{
		Kafka: KafkaConfig{
			Broker:    "kafka:29092",
			Frame:     "frame",
			Detection: "detection",
		},
		Camera: camera,
		Sadm: SadmConfig{
			Port: ":5090",
		},
		Zookeeper: ZookeeperConfig{
			Url:     "zookeeper:2182",
			GroupId: "",
		},
	}
}
