package metrics_integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("RedisBrokerMetrics", func() {

	Context("Nodes available", func() {
		It("prints the number of nodes available", func() {

			redisBrokerStub, err := gexec.Build("github.com/pivotal-cf/cf-redis-broker/integration/metrics/redis_broker_stub")
			Expect(err).ToNot(HaveOccurred())
			redis_broker_metrics, err := gexec.Build("github.com/pivotal-cf/cf-redis-broker/cmd/metrics")
			Expect(err).ToNot(HaveOccurred())

			stub := exec.Command(redisBrokerStub)
			redisBrokerStubSession, err := gexec.Start(stub, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			command := exec.Command(redis_broker_metrics, "--config", "assets/redis_broker_stub.yml")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			Eventually(session.Out).Should(
				gbytes.Say(`\[{"key":"dedicated_vm_total","value":1}\]`))

			Expect(redisBrokerStubSession.Kill().ExitCode()).To(Equal(-1))
		})
	})
})
