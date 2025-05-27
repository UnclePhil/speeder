package monitors

import (
	"speeder/config"
	"time"
)

type PingResult struct {
	Type       string    `json:"type"`
	AgentName  string    `json:"agent_name"`
	Name       string    `json:"name"`
	Target     string    `json:"target"`
	Success    bool      `json:"success"`
	RoundTrip  float64   `json:"roundtrip_ms"`
	PacketLoss float64   `json:"packet_loss"`
	Timestamp  time.Time `json:"timestamp"`
	Error      string    `json:"error,omitempty"`
}

func RunPingTest(test config.PingTest, agentName string) PingResult {
	pinger, err := ping.NewPinger(test.Target)
	if err != nil {
		return PingResult{
			Type:      "ping",
			AgentName: agentName,
			Name:      test.Name,
			Target:    test.Target,
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		}
	}

	pinger.Count = test.Count
	pinger.Timeout = time.Duration(test.TimeoutSec) * time.Second
	
	err = pinger.Run()
	if err != nil {
		return PingResult{
			Type:      "ping",
			AgentName: agentName,
			Name:      test.Name,
			Target:    test.Target,
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		}
	}

	stats := pinger.Statistics()
	return PingResult{
		Type:       "ping",
		AgentName:  agentName,
		Name:       test.Name,
		Target:     test.Target,
		Success:    true,
		RoundTrip:  float64(stats.AvgRtt.Milliseconds()),
		PacketLoss: stats.PacketLoss,
		Timestamp:  time.Now(),
	}
}