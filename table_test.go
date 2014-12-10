package router

import (
	"testing"
)

func TestTable(t *testing.T) {

	na := &mockNode{addr: "abc"}
	nc := &mockNode{addr: "abd"}
	nd := &mockNode{addr: "add"}
	ne := &mockNode{addr: "ddd"}

	pa := &mockPacket{"abc"}
	pb := &mockPacket{"abc"}
	pc := &mockPacket{"abd"}
	pd := &mockPacket{"add"}
	pe := &mockPacket{"ddd"}

	table := &Table{
		Entries: []TableEntry{
			TableEntry{na.Address(), na},
			TableEntry{nc.Address(), nc},
			TableEntry{nd.Address(), nd},
			TableEntry{ne.Address(), ne},
		},
	}

	s := NewSwitch(strAddr("sss"), table, []Node{na, nc, nd, ne})
	s.HandlePacket(pa, na)
	s.HandlePacket(pb, na)
	s.HandlePacket(pc, na)
	s.HandlePacket(pd, na)
	s.HandlePacket(pe, na)

	tt := func(n *mockNode, pkts []Packet) {
		for i, p := range pkts {
			if len(n.pkts) <= i {
				t.Error("pkts not handled in order.", n, pkts)
				return
			}
			if n.pkts[i] != p {
				t.Error("pkts not handled in order.", n, pkts)
			}
		}
	}

	tt(na, []Packet{pa, pb})
	tt(nc, []Packet{pc})
	tt(nd, []Packet{pd})
	tt(ne, []Packet{pe})
}
