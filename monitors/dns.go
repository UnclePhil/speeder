package monitors

import (
	"context"
	"net"
	"speeder/config"
	"time"
)

type DNSResult struct {
	Type      string    `json:"type"`
	AgentName string    `json:"agent_name"`
	Name      string    `json:"name"`
	Server    string    `json:"server"`
	Query     string    `json:"query"`
	Success   bool      `json:"success"`
	Response  []string  `json:"response,omitempty"`
	Duration  float64   `json:"duration_ms"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

func RunDNSTest(test config.DNSTest, agentName string) []DNSResult {
	var results []DNSResult
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, "udp", test.Server)
		},
	}

	for _, query := range test.Queries {
		start := time.Now()
		ips, err := r.LookupHost(context.Background(), query)
		duration := time.Since(start)

		result := DNSResult{
			Type:      "dns",
			AgentName: agentName,
			Name:      test.Name,
			Server:    test.Server,
			Query:     query,
			Duration:  float64(duration.Milliseconds()),
			Timestamp: time.Now(),
		}

		if err != nil {
			result.Success = false
			result.Error = err.Error()
		} else {
			result.Success = true
			result.Response = ips
		}

		results = append(results, result)
	}

	return results
}