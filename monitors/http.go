package monitors

import (
	"net/http"
	"speeder/config"
	"time"
)

type HTTPResult struct {
	Type       string    `json:"type"`
	AgentName  string    `json:"agent_name"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Success    bool      `json:"success"`
	StatusCode int       `json:"status_code"`
	Duration   float64   `json:"duration_ms"`
	Timestamp  time.Time `json:"timestamp"`
	Error      string    `json:"error,omitempty"`
}

func RunHTTPTest(test config.HTTPTest, agentName string) HTTPResult {
	client := &http.Client{
		Timeout: time.Duration(test.TimeoutSec) * time.Second,
	}

	req, err := http.NewRequest(test.Method, test.URL, nil)
	if err != nil {
		return HTTPResult{
			Type:      "http",
			AgentName: agentName,
			Name:      test.Name,
			URL:       test.URL,
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		}
	}

	for key, value := range test.Headers {
		req.Header.Set(key, value)
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	result := HTTPResult{
		Type:      "http",
		AgentName: agentName,
		Name:      test.Name,
		URL:       test.URL,
		Duration:  float64(duration.Milliseconds()),
		Timestamp: time.Now(),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if test.ExpectCode != 0 {
		result.Success = resp.StatusCode == test.ExpectCode
	} else {
		result.Success = resp.StatusCode >= 200 && resp.StatusCode < 400
	}

	return result
}