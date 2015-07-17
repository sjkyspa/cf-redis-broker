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
)

var (
	metronHost    string
	metronPort    int
	redisHost     string
	redisPort     int
	redisPassword string
)

func main() {
	flag.StringVar(&redisHost, "redis-host", "localhost", "Redis host")
	flag.IntVar(&redisPort, "redis-port", 6379, "Redis port")
	flag.StringVar(&redisPassword, "redis-password", "", "Redis password")
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

func capture() {
	client, err := client.Connect(
		client.Host(redisHost),
		client.Port(redisPort),
		client.Password(redisPassword),
	)
	if err != nil {
		log.Println("Error connecting to Redis: " + err.Error())
	}

	info, err := client.Info()
	if err != nil {
		log.Println("Error getting info form Redis: " + err.Error())
	}

	client.Disconnect()

	cpu, err := strconv.ParseFloat(info["used_cpu_sys"], 64)
	if err != nil {
		log.Println("Error getting info form Redis: " + err.Error())
	}

	fmt.Printf("CPU: %.2f\n", cpu)

	if err = metrics.SendValue("cpu", cpu, "Load"); err != nil {
		log.Println("Error emitting metric: " + err.Error())
	}
}
