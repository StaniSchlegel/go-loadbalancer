package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"loadbalancer/api"
	"loadbalancer/data"
	"loadbalancer/loadbalancer"

	"gopkg.in/yaml.v2"
)

var loadBalancerConfig loadbalancer.LoadBalancerConfiguration

func main() {

	configFilePath := flag.String("config-file", "config.yml", "Local path to the config yaml file that contains the load balancer configuration")
	flag.Parse()

	configFileContent, err := ioutil.ReadFile(*configFilePath)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFileContent, &loadBalancerConfig)

	if err != nil {
		fmt.Println("Could not parse config yaml file")
		panic(err)
	}

	data.DbConnection, err = data.CreateDbConnection()

	if err != nil {
		panic(err)
	}

	apiServerAddress := fmt.Sprintf(
		"%s:%v",
		loadBalancerConfig.ApiServerConfiguration.Host,
		loadBalancerConfig.ApiServerConfiguration.Port,
	)
	apiServer := api.NewApiServer(apiServerAddress)
	go apiServer.StartListening()

	loadbalancer := loadbalancer.NewLoadBalancer(loadBalancerConfig)
	loadbalancer.StartListening()
}
