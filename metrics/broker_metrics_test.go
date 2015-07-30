package metrics_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf/cf-redis-broker/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BrokerMetrics", func() {
	var (
		config        metrics.Config
		brokerMetrics *metrics.BrokerMetrics
		brokerServer  *httptest.Server
	)

	Describe(".Metrics", func() {
		BeforeEach(func() {
			brokerServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `{"pool":{"count":3,"clusters":[["10.10.32.9"]]},"allocated":{"count":2,"clusters":null}}`)
			}))

			config = metrics.Config{
				BrokerEndpoint: brokerServer.URL,
			}

			brokerMetrics = metrics.NewBrokerMetrics(config)
		})

		AfterEach(func() {
			brokerServer.Close()
		})

		It("returns the correct metrics", func() {
			fetchedMetrics, err := brokerMetrics.FetchMetrics()
			Expect(err).ToNot(HaveOccurred())

			Expect(fetchedMetrics).To(Equal(
				metrics.Metrics{
					metrics.Metric{Key: "dedicated_vm_total", Value: 3},
					metrics.Metric{Key: "dedicated_vm_available", Value: 1},
				},
			))
		})
	})
})
