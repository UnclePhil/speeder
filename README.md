# Speeder - Network Monitoring Agent

A cross-platform network monitoring agent written in Go that performs ping, DNS, and HTTP(S) tests and reports metrics via MQTT.

## Features

- Ping tests with customizable count and timeout
- DNS resolution tests with multiple queries support
- HTTP/HTTPS tests with custom headers and status code validation
- Metrics reporting via MQTT
- Configurable via YAML or environment variables
- Cross-platform support (Windows/Linux)

## Configuration

The agent can be configured using either a YAML file or environment variables.

### YAML Configuration

Default path: `config.yaml` (can be overridden with `CONFIG_PATH` environment variable)

See the provided `config.yaml` for a complete example.

### Environment Variables

- `CONFIG_PATH`: Path to the YAML config file
- `MQTT_BROKER`: MQTT broker URL
- `MQTT_TOPIC`: Base topic for publishing metrics
- `MQTT_CLIENT_ID`: MQTT client identifier
- `MQTT_USERNAME`: MQTT username
- `MQTT_PASSWORD`: MQTT password
- `CHECK_INTERVAL`: Interval between check cycles in seconds
- `AGENT_NAME`: Unique identifier for this agent instance

## Metrics Format

All metrics include these standardized fields to simplify aggregation and analysis:

- `type`: The type of test ("ping", "dns", or "http")
- `agent_name`: The unique identifier of the agent that performed the test
- `name`: The name of the test as configured
- `success`: Boolean indicating if the test passed
- `timestamp`: ISO 8601 timestamp of when the test was performed
- `error`: Error message if the test failed (omitted if successful)

The agent publishes JSON-formatted metrics to the following MQTT topics:

- `{base_topic}/{agentname}/ping`: Ping test results
- `{base_topic}/{agentname}/dns`: DNS test results
- `{base_topic}/{agentname}/http`: HTTP test results

### Sample Metrics

#### Ping Metrics
```json
{
  "type": "ping",
  "agent_name": "agent1",
  "name": "Google DNS",
  "target": "8.8.8.8",
  "success": true,
  "roundtrip_ms": 15.4,
  "packet_loss": 0,
  "timestamp": "2025-05-27T10:00:00Z"
}
```

#### DNS Metrics
```json
{
  "type": "dns",
  "agent_name": "agent1",
  "name": "Google DNS Test",
  "server": "8.8.8.8:53",
  "query": "google.com",
  "success": true,
  "response": ["142.250.180.78"],
  "duration_ms": 25.3,
  "timestamp": "2025-05-27T10:00:00Z"
}
```

#### HTTP Metrics
```json
{
  "type": "http",
  "agent_name": "agent1",
  "name": "GitHub API",
  "url": "https://api.github.com/status",
  "success": true,
  "status_code": 200,
  "duration_ms": 150.2,
  "timestamp": "2025-05-27T10:00:00Z"
}
```