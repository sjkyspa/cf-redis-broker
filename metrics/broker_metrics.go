package metrics

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type BrokerMetrics struct {
	config Config
}

func NewBrokerMetrics(config *Config) *BrokerMetrics {
	return &BrokerMetrics{
		config: *config,
	}
}

func (b *BrokerMetrics) FetchMetrics() (Metrics, error) {
	resp, err := http.Get(b.config.BrokerEndpoint)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var debugResponse DebugResponse
	json.Unmarshal(body, &debugResponse)

	metrics := Metrics{}
	metrics = append(metrics, b.createMetric("dedicated_vm_total", debugResponse.Pool.Count))
	metrics = append(metrics, b.createMetric("dedicated_vm_available", debugResponse.Pool.Count-debugResponse.Allocated.Count))

	return metrics, nil
}

func (b *BrokerMetrics) createMetric(name string, value int) Metric {
	return Metric{
		Key:   name,
		Value: value,
	}
}
