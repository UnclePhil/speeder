package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"
	"speeder/config"
	"speeder/monitors"
)

func loadConfig() (*config.Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	return &cfg, nil
}

func createMQTTClient(cfg *config.Config) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.MQTT.Broker).
		SetClientID(cfg.MQTT.ClientID).
		SetUsername(cfg.MQTT.Username).
		SetPassword(cfg.MQTT.Password)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func publishResults(client mqtt.Client, topic string, results interface{}) error {
	data, err := json.Marshal(results)
	if err != nil {
		return err
	}

	token := client.Publish(topic, 0, false, data)
	token.Wait()
	return token.Error()
}

func runMonitors(cfg *config.Config, mqttClient mqtt.Client) {
	// Run ping tests
	for _, test := range cfg.Tests.Ping {
		result := monitors.RunPingTest(test)
		err := publishResults(mqttClient, cfg.MQTT.Topic+"/ping", result)
		if err != nil {
			log.Printf("Error publishing ping results: %v", err)
		}
	}

	// Run DNS tests
	for _, test := range cfg.Tests.DNS {
		results := monitors.RunDNSTest(test)
		err := publishResults(mqttClient, cfg.MQTT.Topic+"/dns", results)
		if err != nil {
			log.Printf("Error publishing DNS results: %v", err)
		}
	}

	// Run HTTP tests
	for _, test := range cfg.Tests.HTTP {
		result := monitors.RunHTTPTest(test)
		err := publishResults(mqttClient, cfg.MQTT.Topic+"/http", result)
		if err != nil {
			log.Printf("Error publishing HTTP results: %v", err)
		}
	}
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mqttClient, err := createMQTTClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer mqttClient.Disconnect(0)

	interval := time.Duration(cfg.Global.Interval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Starting monitoring with interval: %v", interval)

	for {
		runMonitors(cfg, mqttClient)
		<-ticker.C
	}
}