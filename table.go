package router

// TableEntry is an (Address, Node) pair for a routing Table
type TableEntry struct {
	Address Address
	NextHop Node
}

// Table is a Router (Routing Table, really) based on a distance criterion.
type Table struct {
	Entries []TableEntry

	// Distance returns a measure of distance between two Addresses.
	Distance DistanceFunc
}

// Route decides how to route a Packet out of a list of Nodes.
// It returns the Node chosen to send the Packet to.
// Route may return nil, if no route is suitable at all (equivalent of drop).
func (t *Table) Route(p Packet, ns []Node) Node {
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
