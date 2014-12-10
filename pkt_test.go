package router

import (
	"testing"
)

type strAddr string

func (s strAddr) Distance(o Address) int {
	so, ok := o.(strAddr)
	if !ok {
		return -1
	}

	for i := range s {
		if s[i] != so[i] {
			return len(s) - i
		}
	}
	return 0
}

type mockNode struct {
	addr strAddr
	pkts []Packet
}

func (m *mockNode) Address() Address {
	return m.addr
}

func (m *mockNode) HandlePacket(p Packet, n Node) {
	m.pkts = append(m.pkts, p)
}

type mockPacket struct {
	a strAddr
}

func (m *mockPacket) Destination() Address {
	return m.a
}

func (m *mockPacket) Payload() interface{} {
	return m
}

func TestAddrs(t *testing.T) {

	ta := func(a, b Address, expect int) {
		actual := a.Distance(b)
		if actual != expect {
			t.Error("address distance error:", a, b, expect, actual)
		}
	}

	a := strAddr("abc")
	b := strAddr("abc")
	c := strAddr("abd")
	d := strAddr("add")
	e := strAddr("ddd")

	ta(a, a, 0)
	ta(a, b, 0)
	ta(a, c, 1)
	ta(a, d, 2)
	ta(a, e, 3)
}

func TestNodes(t *testing.T) {

	a := &mockPacket{"abc"}
	b := &mockPacket{"abc"}
	c := &mockPacket{"abd"}
	d := &mockPacket{"add"}
	e := &mockPacket{"ddd"}

	n := &mockNode{addr: "abc"}
	n2 := &mockNode{addr: "ddd"}
	n.HandlePacket(a, n2)
	n.HandlePacket(b, n2)
	n.HandlePacket(c, n2)
	n.HandlePacket(d, n2)
	n.HandlePacket(e, n2)

	pkts := []Packet{a, b, c, d, e}
	for i, p := range pkts {
		if n.pkts[i] != p {
			t.Error("pkts not handled in order.")
		}
	}
}
