package k8sutil

import "net"

// ACXPlusNodeConfig Configurations pertaining to each ACX Plus Instance
type ACXPlusNodeConfig struct {
	App            string
	GroupID        string
	Family         string
	Topic          string
	Instances      int32
	KafkaHosts     *[]net.TCPAddr
	CassandraHosts *[]net.TCPAddr
}

// Connection Returns the connection check results
type Connection struct {
	Open  bool
	Addr  net.TCPAddr
	Error error
}
