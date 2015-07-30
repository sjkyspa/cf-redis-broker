package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/pivotal-cf/cf-redis-broker/metrics"
)

func main() {
	configFile := flag.String("config", "", "path to config file")
	flag.Parse()

	config, err := metrics.LoadConfig(*configFile)
	if err != nil {
		panic(err)
	}

	brokerMetrics := metrics.NewBrokerMetrics(config)
	metrics, err := brokerMetrics.FetchMetrics()
	if err != nil {
		panic(err)
	}

	output, err := json.Marshal(metrics)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)

	// config := &metrics.Config{}
	// candiedyaml.NewDecoder(file).Decode(config)
	// return nil, err
	// }

	// resp, _ := http.Get(config.BrokerEndpoint)
	// // defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	// var debugResponse metrics.DebugResponse
	// json.Unmarshal(body, &debugResponse)

	// fmt.Println(fmt.Sprintf(`[{"key":"dedicated_vm_total","value":%d}]`, debugResponse.Pool.Count))
}
