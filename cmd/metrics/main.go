package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cloudfoundry/dropsonde"
	"github.com/cloudfoundry/dropsonde/metrics"
	"github.com/pivotal-cf/cf-redis-broker/redis/client"
	"github.com/pivotal-cf/cf-redis-broker/redisconf"
)

var (
	metronHost    string
	metronPort    int
	redisHost     string
	redisPort     int
	redisConfPath string
)

func main() {
	flag.StringVar(&redisConfPath, "redis-config", "", "Path to redis config file")
	flag.StringVar(&redisHost, "redis-host", "localhost", "Redis host")
	flag.IntVar(&redisPort, "redis-port", 6379, "Redis port")
	flag.StringVar(&metronHost, "metron-host", "localhost", "Metron host")
	flag.IntVar(&metronPort, "metron-port", 3457, "Metron port")

	flag.Parse()

	dropsonde.Initialize(fmt.Sprintf("%s:%d", metronHost, metronPort), "redis")

	capture()
	for {
		select {
		case <-time.After(5 * time.Second):
			capture()
		}
	}
}

func password() (string, error) {
	conf, err := redisconf.Load(redisConfPath)
	if err != nil {
		return "", err
	}

	return conf.Get("requirepass"), nil
}

func capture() {
	pwd, err := password()
	if err != nil {
		log.Println("Error getting password from redis config: " + err.Error())
		return
	}

	client, err := client.Connect(
		client.Host(redisHost),
		client.Port(redisPort),
		client.Password(pwd),
	)
	if err != nil {
		log.Println("Error connecting to Redis: " + err.Error())
		return
	}

	info, err := client.Info()
	if err != nil {
		log.Println("Error getting info form Redis: " + err.Error())
		return
	}

	client.Disconnect()

	cpu, err := strconv.ParseFloat(info["used_cpu_sys"], 64)
	if err != nil {
		log.Println("Error getting info form Redis: " + err.Error())
		return
	}

	fmt.Printf("CPU: %.2f\n", cpu)

	if err = metrics.SendValue("cpu", cpu, "Load"); err != nil {
		log.Println("Error emitting metric: " + err.Error())
		return
	}
}
