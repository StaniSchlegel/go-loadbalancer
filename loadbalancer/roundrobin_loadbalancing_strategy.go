package loadbalancer

type RoundRobinLoadBalancingStrategy struct {
	counter int
}

func (r *RoundRobinLoadBalancingStrategy) ChooseNode(nodeAddresses []string) string {

	r.counter += 1

	if r.counter > len(nodeAddresses)-1 {
		r.counter = 0
	}

	return nodeAddresses[r.counter]
}
