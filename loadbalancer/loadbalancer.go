package loadbalancer

import (
	"errors"
	"fmt"
	"loadbalancer/data/repositories"
	"net"
	"time"
)

type LoadBalancerConfiguration struct {
	LoadBalancerServerConfiguration LoadBalancerServerConfiguration `yaml:"loadbalancer"`
	ApiServerConfiguration          ApiServerConfiguration          `yaml:"apiServer"`
	Rules                           []LoadBalancingRule             `yaml:"rules"`
}

type LoadBalancerServerConfiguration struct {
	Host                  string `yaml:"host"`
	Port                  int    `yaml:"port"`
	LoadbalancingStrategy string `yaml:"loadbalancingStrategy"`
}

type ApiServerConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LoadBalancingRule struct {
	name string `yaml:"name"`
}

type LoadBalancer struct {
	IpAddress string
	Port      int

	tcpListener net.Listener

	LoadBalancingStrategy LoadBalancingStrategy
}

func NewLoadBalancer(config LoadBalancerConfiguration) LoadBalancer {
	l := LoadBalancer{
		IpAddress: config.LoadBalancerServerConfiguration.Host,
		Port:      config.LoadBalancerServerConfiguration.Port,
	}

	switch config.LoadBalancerServerConfiguration.LoadbalancingStrategy {
	case "roundrobin":
		{
			l.LoadBalancingStrategy = &RoundRobinLoadBalancingStrategy{}
		}
	case "random":
		{
			l.LoadBalancingStrategy = &RandomLoadBalancingStrategy{}
		}
	default:
		{
			l.LoadBalancingStrategy = &RandomLoadBalancingStrategy{}
		}
	}

	return l
}

func (l *LoadBalancer) StartListening() error {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf("%s:%v", l.IpAddress, l.Port))

	if err != nil {
		return err
	}

	l.tcpListener = tcpListener

	defer l.tcpListener.Close()

	for {
		conn, err := l.tcpListener.Accept()

		if err != nil {
			return err
		}

		go l.handleRequest(conn)
	}

}

func (l *LoadBalancer) StopListening() {
	if l.tcpListener != nil {
		l.tcpListener.Close()
	}
}

func (l *LoadBalancer) handleRequest(conn net.Conn) error {

	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	bufferSize := 1024
	buffer := make([]byte, bufferSize)

	var fullBuffer []byte
	totalReadBytes := 0

	for {
		readBytes, err := conn.Read(buffer)
		totalReadBytes += readBytes

		netErr, ok := err.(net.Error)

		if ok && netErr.Timeout() {
			fmt.Println("read timeout:", err)
			return nil
		} else if err != nil {
			fmt.Println("read error:", err)
			return err
		}

		fullBuffer = append(fullBuffer, buffer...)

		if readBytes < bufferSize {
			break
		}
	}

	forwardingAddresses, err := repositories.GetAllNodeAddresses()

	if err != nil {
		return err
	}

	responseSent := false

	for responseSent == false && len(forwardingAddresses) > 0 {
		nodeAddress := l.LoadBalancingStrategy.ChooseNode(forwardingAddresses)
		response, err := l.sendDataToNode(nodeAddress, fullBuffer[0:totalReadBytes])

		if err != nil {
			responseSent = false
		} else {
			_, err = conn.Write(response)
			responseSent = (err == nil)
		}

		if !responseSent {
			forwardingAddresses = removeItemFromSlice(forwardingAddresses, nodeAddress)
		}

	}

	if !responseSent {
		fmt.Println("Could not forward data to any forwaring address")
		return errors.New("Data could not be forwarded to any address")
	}

	return nil
}

func (l LoadBalancer) sendDataToNode(nodeAddress string, data []byte) ([]byte, error) {

	bufferSize := 1024
	buffer := make([]byte, bufferSize)

	fullBuffer := make([]byte, 0)
	totalReadBytes := 0

	clientConn, err := net.Dial("tcp", nodeAddress)

	if err != nil {
		return nil, err
	}

	clientConn.Write(data)

	for {
		readBytes, err := clientConn.Read(buffer)
		totalReadBytes += readBytes

		netErr, ok := err.(net.Error)

		if ok && netErr.Timeout() {
			fmt.Println("read timeout:", err)
			return nil, netErr
		} else if err != nil {
			fmt.Println("read error:", err)
			return nil, err
		}

		fullBuffer = append(fullBuffer, buffer...)

		if readBytes < bufferSize {
			break
		}
	}

	return fullBuffer[0:totalReadBytes], nil
}

func removeItemFromSlice(haystack []string, needle string) []string {
	needleIdx := -1

	for idx, item := range haystack {
		if item == needle {
			needleIdx = idx
		}
	}

	if needleIdx != -1 {
		return append(haystack[:needleIdx], haystack[needleIdx+1:]...)
	}

	return nil
}
