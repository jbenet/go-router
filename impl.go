package router

type switchh struct {
	addr   Address
	router Router
	nodes  []Node
}

// NewSwitch constructs a switch with given Router and list of adjacent Nodes.
func NewSwitch(a Address, r Router, adj []Node) Switch {
	return &switchh{a, r, adj}
}

func (s *switchh) Address() Address {
	return s.addr
}

func (s *switchh) Router() Router {
	return s.router
}

func (s *switchh) HandlePacket(p Packet, n Node) {
	next := s.router.Route(p, s.nodes)
	if next != nil {
		next.HandlePacket(p, s)
	}
}
