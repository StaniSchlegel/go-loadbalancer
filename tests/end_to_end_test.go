package tests

import (
	"bytes"
	"errors"
	"fmt"
	"loadbalancer/api"
	"loadbalancer/data"
	"loadbalancer/data/repositories"
	"loadbalancer/loadbalancer"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

var receivedMessages map[string][]string

func TestLoadBalancingBetweenTwoNodes(t *testing.T) {

	// Delete testing sqlite db if existing
	os.Remove("loadbalancer.db")

	receivedMessages = make(map[string][]string)

	loadBalancerConfig := loadbalancer.LoadBalancerConfiguration{
		LoadBalancerServerConfiguration: loadbalancer.LoadBalancerServerConfiguration{
			Host:                  "127.0.0.1",
			Port:                  8000,
			LoadbalancingStrategy: "roundrobin",
		},
		ApiServerConfiguration: loadbalancer.ApiServerConfiguration{
			Host: "127.0.0.1",
			Port: 3000,
		},
		Rules: []loadbalancer.LoadBalancingRule{},
	}

	err := errors.New("")
	data.DbConnection, err = data.CreateDbConnection()

	if err != nil {
		t.Fatalf(err.Error())
	}

	apiServerAddress := fmt.Sprintf(
		"%s:%v",
		loadBalancerConfig.ApiServerConfiguration.Host,
		loadBalancerConfig.ApiServerConfiguration.Port,
	)
	apiServer := api.NewApiServer(apiServerAddress)
	go apiServer.StartListening()

	loadbalancer := loadbalancer.NewLoadBalancer(loadBalancerConfig)
	go loadbalancer.StartListening()

	time.Sleep(2 * time.Second)

	// Setup virtual "nodes"
	firstNodeApiServer := fiber.New()
	firstNodeApiServer.Post("/test", func(c *fiber.Ctx) error {
		receivedMessages["firstNode"] = append(receivedMessages["firstNode"], string(c.Body()))
		return c.Status(200).SendString("Received message")
	})
	go firstNodeApiServer.Listen(":3001")

	secondNodeApiServer := fiber.New()
	secondNodeApiServer.Post("/test", func(c *fiber.Ctx) error {
		receivedMessages["secondNode"] = append(receivedMessages["secondNode"], string(c.Body()))
		return c.Status(200).SendString("Received message")
	})
	go secondNodeApiServer.Listen(":3002")

	time.Sleep(2 * time.Second)

	// Add nodes to list of forwarding addresses of loadbalancer
	_, err = http.Post("http://127.0.0.1:3000/nodes", "application/json", bytes.NewBuffer([]byte(`{ "address": "127.0.0.1:3001" }`)))
	_, err = http.Post("http://127.0.0.1:3000/nodes", "application/json", bytes.NewBuffer([]byte(`{ "address": "127.0.0.1:3002" }`)))

	time.Sleep(2 * time.Second)

	// Check if nodes were added to db
	allNodeAddresses, err := repositories.GetAllNodeAddresses()

	if len(allNodeAddresses) < 2 {
		t.Fatalf("Failed to add node addresses to loadbalancer db via the api server")
	}

	_, err = http.Post("http://127.0.0.1:8000/test", "text/plain", bytes.NewBuffer([]byte("testing")))

	if err != nil {
		fmt.Println(err)
		t.Fatalf("Could not send testing message to loadbalancer")
	}

	_, err = http.Post("http://127.0.0.1:8000/test", "text/plain", bytes.NewBuffer([]byte("testing")))

	if err != nil {
		fmt.Println(err)
		t.Fatalf("Could not send testing message to loadbalancer")
	}

	if messages, ok := receivedMessages["firstNode"]; ok {
		if len(messages) != 1 {
			t.Fatalf("Loadbalancer failed to equally balance load between first and second node")
		}
	}

	if messages, ok := receivedMessages["secondNode"]; ok {
		if len(messages) != 1 {
			t.Fatalf("Loadbalancer failed to equally balance load between first and second node")
		}
	}

	os.Remove("loadbalancer.db")
}
