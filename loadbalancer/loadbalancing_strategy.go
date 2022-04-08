package loadbalancer

type LoadBalancingStrategy interface {
	ChooseNode([]string) string
}
