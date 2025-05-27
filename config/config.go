package config

type Config struct {
	MQTT   MQTTConfig   `yaml:"mqtt"`
	Tests  TestsConfig  `yaml:"tests"`
	Global GlobalConfig `yaml:"global"`
}

type MQTTConfig struct {
	Broker   string `yaml:"broker" env:"MQTT_BROKER"`
	Topic    string `yaml:"topic" env:"MQTT_TOPIC"`
	ClientID string `yaml:"clientId" env:"MQTT_CLIENT_ID"`
	Username string `yaml:"username" env:"MQTT_USERNAME"`
	Password string `yaml:"password" env:"MQTT_PASSWORD"`
}

type TestsConfig struct {
	Ping  []PingTest  `yaml:"ping"`
	DNS   []DNSTest   `yaml:"dns"`
	HTTP  []HTTPTest  `yaml:"http"`
}

type GlobalConfig struct {
	Interval  int    `yaml:"interval" env:"CHECK_INTERVAL"`
	AgentName string `yaml:"agentName" env:"AGENT_NAME"`
}

type PingTest struct {
	Name        string `yaml:"name"`
	Target      string `yaml:"target"`
	Count       int    `yaml:"count"`
	TimeoutSec  int    `yaml:"timeout"`
}

type DNSTest struct {
	Name     string   `yaml:"name"`
	Server   string   `yaml:"server"`
	Queries  []string `yaml:"queries"`
	Type     string   `yaml:"type"`
}

type HTTPTest struct {
	Name        string            `yaml:"name"`
	URL         string            `yaml:"url"`
	Method      string            `yaml:"method"`
	Headers     map[string]string `yaml:"headers"`
	TimeoutSec  int              `yaml:"timeout"`
	ExpectCode  int              `yaml:"expectCode"`
}