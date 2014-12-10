// Package router implements a router that routes messages across a set
// of interfaces. It is reminiscent of computer networks (internet), but
// meant to be generalized.
package router

// Address is our way of knowing where we're headed. Traditionally, addresses
// are things like IP Addresses, or email addresses. But filepaths, or URLs
// can be seen as addresses too. We really leave it up to you. Our routing
// is general, and forces you to pick an addressing scheme, and some logic to
// discriminate addresses that you'll plug into Switches and Routers.
type Address interface{}

// Packet is the unit of moving things in our network. Anything can be routed
// in our network, as long as it has a Destination.
type Packet interface {

	// Destination is the Address of the endpoint this Packet is headed to.
	// They could use the same Addressing throughout (recommended) like the
	// internet, or face the pain of translating addresses at inter-network
	// gateways.
	Destination() Address

	// Payload is here for completeness. This can be the Packet itself, but this
	// function encourages the client to think through their packet design.
	Payload() interface{}
}

// // Interface is an entry point to a Node in our network. It represents
// // the way to communicate to another Node, but note that the interface has
// // its own address. This is so Nodes can have multiple Addresses, and
// // potentially be part of different networks altogether. This means that
// // Interfaces could actually carry some logic to them. For example, they
// // could modify Packets (e.g. serialization, Address mapping, etc).
// type Interface interface {
//
// 	// Address returns the Address of this Interface. Note that this is the
// 	// Address of the Device owning the Interface, at this Interface. Sometimes
// 	// devices have the same Address in all their interfaces. Sometimes they
// 	// have different addresses.
// 	Address() Address
//
// 	// WritePacket puts a Packet into the Interface. Depending on the interface
// 	// implementation, this may be a blocking call.
// 	WritePacket(Packet) error
//
// 	// ReadPacket receives the next Packet from the Interface. Depending on the
// 	// interface implementation, this may be a blocking call.
// 	ReadPacket() (Packet, error)
// }
//
// // A Connector is an object that connects two (or more) Interfaces together.
// // It's like a physical wire between two devices, or wireless medium. Depending
// // on the implementation, Connectors may be running goroutines and need to be
// // stopped.
// type Connector interface {
//
// 	// Interfaces returns all the Interfaces Connector connects.
// 	Interfaces() []Interface
//
// 	// Disconnect stops this connector. This is usually a destructive call.
// 	Disconnect()
// }

// Node is an object which has interfaces to connect to networks. This
// is an "endpoint" object.
type Node interface {

	// // Interfaces returns all the existing interfaces of the Node.
	// // Note that -- depending on the Node implementation -- Interfaces may
	// // be added or removed at any time.
	// Interfaces() []Interface

	// HandlePacket receives a packet sent by another node.
	HandlePacket(Packet, Node)
}

// Switch is the basic forwarding device. It listens to all its interfaces (it
// is a Node in our network) and Forwards all Packets received by any interface
// out another interface, according to its ForwardingTable.
type Switch interface {
	Node

	// Router returns the Router used to decide where to forward packets.
	Router() Router
}

// Router is an object that decides how a Packet should be Routed.
//
// Note that this is a break from traditional networking systems. Instead of
// having the abstractions of FIB, RIB, Routing/Forwarding Tables, Routers,
// and Switches, we only have the last two:
// - Router -- the things that "route" (decide where things go)
// - Switch -- connecting Nodes, "switch" Packets according to a Router.
//
// Our Router captures the entire Control Plane, meaning that we can implement:
// - Static Routing - forwarding table only
// - Dynamic Routing - routing table computed with an algorithm or protocol
// - "SDN" Routing - Control Plane separated from Data Plane
// And even:
// - URL Routers (like gorilla.Muxer)
// - Protocol Muxers
// entirely within different Router implementations.
type Router interface {

	// Route decides how to route a Packet out of a list of Nodes.
	// It returns the Node chosen to send the Packet to.
	// Route may return nil, if no route is suitable at all (equivalent of drop).
	Route(Packet, []Node) Node
}