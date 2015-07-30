package metrics_test

import (
	"github.com/pivotal-cf/cf-redis-broker/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe(".LoadConfig", func() {
		Context("when config file is missing", func() {
			It("returns an error", func() {
				_, err := metrics.LoadConfig("assets/missing_config.yml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("config: Failed to open file"))
			})
		})

		Context("when broker endpoint property is missing", func() {
			It("returns an error", func() {
				_, err := metrics.LoadConfig("assets/invalid_config.yml")
				Expect(err).To(MatchError("config: Missing BrokerEndpoint property"))
			})
		})

		It("loads the config file correctly", func() {
			config, err := metrics.LoadConfig("assets/correct_config.yml")
			Expect(err).ToNot(HaveOccurred())
			Expect(config.BrokerEndpoint).To(Equal("http://localhost:5555/debug"))
		})
	})
})
