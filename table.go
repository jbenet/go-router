package router

// TableEntry is an (Address, Node) pair for a routing Table
type TableEntry struct {
	Address Address
	NextHop Node
}

// Table is a Router (Routing Table, really) based on a distance criterion.
//
// For example:
//
//   n1 := NewQueueNode("aaa", make(chan Packet, 10))
//   n2 := NewQueueNode("aba", make(chan Packet, 10))
//   n3 := NewQueueNode("abc", make(chan Packet, 10))
//
//   var t router.Table
//   t.Distance = router.HammingDistance
//   t.AddNodes(n1, n2)
//
//   p1 := NewPacket("aaa", "hello1")
//   p2 := NewPacket("aba", "hello2")
//   p3 := NewPacket("abc", "hello3")
//
//   t.Route(p1) // n1
//   t.Route(p2) // n2
//   t.Route(p3) // n2, because we don't have n3 and n2 is closet
//
//   t.AddNode(n3)
//   t.Route(p3) // n3
type Table struct {
	Entries []TableEntry

	// Distance returns a measure of distance between two Addresses.
	Distance DistanceFunc
}

// AddEntry adds an (Address, NextHop) entry to the Table
func (t *Table) AddEntry(addr Address, nextHop Node) {
	t.Entries = append(t.Entries, TableEntry{addr, nextHop})
}

// AddNode calls AddTableEntry for the given Node
func (t *Table) AddNode(n Node) {
	t.AddEntry(n.Address(), n)
}

// AddNodes calls AddTableEntry for the given Node
func (t *Table) AddNodes(ns ...Node) {
	for _, n := range ns {
		t.AddNode(n)
	}
}

// Route decides how to route a Packet out of a list of Nodes.
// It returns the Node chosen to send the Packet to.
// Route may return nil, if no route is suitable at all (equivalent of drop).
func (t *Table) Route(p Packet) Node {
	if t.Entries == nil {
		return nil
	}

	dist := t.Distance
	if dist == nil {
		dist = unitDistance
	}

	var best Node
	var bestDist int
	var addr = p.Destination()

	for _, e := range t.Entries {
		d := dist(e.Address, addr)
		if d < 0 {
			continue
		}
		if best == nil || d < bestDist {
			bestDist = d
			best = e.NextHop
		}
	}
	return best
}

func unitDistance(a, b Address) int {
	if a == b {
		return 0
	}
	return -1
}
