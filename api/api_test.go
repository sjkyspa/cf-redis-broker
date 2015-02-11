package api_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf/cf-redis-broker/api"
	"github.com/pivotal-cf/cf-redis-broker/credentials"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeRedisResetter struct {
	deleteAllData func() error
}

func (client *fakeRedisResetter) ResetRedis() error {
	return client.deleteAllData()
}

var _ = Describe("redis agent HTTP API", func() {
	var server *httptest.Server
	var redisClient *fakeRedisResetter
	var deleteCount int
	var configPath string
	var response *http.Response

	var parseCredentials func(string) (credentials.Credentials, error)

	BeforeEach(func() {
		parseCredentials = func(path string) (credentials.Credentials, error) {
			Ω(path).Should(Equal(configPath))
			return credentials.Credentials{
				Port:     123345,
				Password: "secret",
			}, nil
		}
		configPath = "/some/Config/Path"
		redisClient = &fakeRedisResetter{}
		deleteCount = 0
	})

	JustBeforeEach(func() {
		handler := api.New(redisClient, configPath, parseCredentials)
		server = httptest.NewServer(handler)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("GET /", func() {
		JustBeforeEach(func() {
			response = makeRequest("GET", server.URL)
		})

		Context("When it can read the conf file successfully", func() {
			It("returns the correct credentials", func() {
				body, err := ioutil.ReadAll(response.Body)
				Ω(err).ShouldNot(HaveOccurred())

				creds := credentials.Credentials{}
				err = json.Unmarshal(body, &creds)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(creds).Should(Equal(credentials.Credentials{
					Port:     123345,
					Password: "secret",
				}))
			})
		})

		Context("When it is unable to read the conf file", func() {
			BeforeEach(func() {
				parseCredentials = func(path string) (credentials.Credentials, error) {
					Ω(path).Should(Equal(configPath))
					return credentials.Credentials{}, errors.New("unable to open config file")
				}
			})

			It("returns an 500", func() {
				Ω(response.StatusCode).Should(Equal(500))
			})

			It("returns the correct error in the body", func() {
				body, err := ioutil.ReadAll(response.Body)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(string(body)).Should(Equal("unable to open config file\n"))
			})
		})
	})

	Describe("DELETE /", func() {
		Context("When it can connect to Redis successfully", func() {
			JustBeforeEach(func() {
				redisClient.deleteAllData = func() error {
					deleteCount++
					return nil
				}

				response = makeRequest("DELETE", server.URL)
			})

			It("deletes all data from redis", func() {
				Ω(deleteCount).To(Equal(1))
			})

			It("returns HTTP 200 OK", func() {
				Ω(response.StatusCode).Should(Equal(200))
			})
		})

		Context("when deleting all data from redis goes wrong", func() {
			JustBeforeEach(func() {
				redisClient.deleteAllData = func() error {
					return errors.New("redis burned down")
				}
				response = makeRequest("DELETE", server.URL)
			})

			It("returns 500", func() {
				Ω(response.StatusCode).Should(Equal(500))
			})

			It("returns the correct error in the body", func() {
				body, err := ioutil.ReadAll(response.Body)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(string(body)).Should(Equal("redis burned down\n"))
			})
		})
	})

	Describe("All other HTTP methods", func() {
		for _, method := range []string{"POST", "PUT"} {
			requestMethod := method
			var response *http.Response

			JustBeforeEach(func() {
				response = makeRequest(requestMethod, server.URL)
			})

			It(method+" returns an http error", func() {
				Ω(response.StatusCode).Should(Equal(http.StatusNotFound))
			})
		}
	})
})

func makeRequest(method string, url string) *http.Response {
	request, err := http.NewRequest(method, url, nil)
	Ω(err).ShouldNot(HaveOccurred())

	response, err := http.DefaultClient.Do(request)
	Ω(err).ShouldNot(HaveOccurred())

	return response
}
