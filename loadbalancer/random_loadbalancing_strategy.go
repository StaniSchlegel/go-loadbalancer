package loadbalancer

import "math/rand"

type RandomLoadBalancingStrategy struct {
}

func (r *RandomLoadBalancingStrategy) ChooseNode(nodeAddresses []string) string {
	idx := rand.Intn(len(nodeAddresses))
	return nodeAddresses[idx]
}
