package k8sutil

import (
	"net"
	"strconv"
	"testing"
)

const namespace = "uat01"
const fileName = "application.conf"
const metaGroupID = "santander50k_115rc1_2-MD"
const metaTopic = "ac.init.metadata.santander50k_115rc1_2"

var util ACXK8sUtil

func TestGetConfigMap(t *testing.T) {
	if util.clientset == nil {
		clientset := BuildClientSet()

		util = ACXK8sUtil{
			clientset: clientset,
			Namespace: namespace,
		}
	}

	product := "acx-plus-metadata"
	config := util.GetConfigMap(product, fileName)

	if config == "" {
		t.Error("ConfigMap Not Retrieved")
	}
}

func TestGetGroupID(t *testing.T) {
	product := "acx-plus-metadata"

	groupID := util.GetGroupID(product)

	if groupID != metaGroupID {
		t.Error("Incorrect Group ID or Not Found")
	}
}

func TestGetKeyspace(t *testing.T) {
	product := "acx-plus-metadata"

	keyspace := util.GetKeyspace(product)

	if keyspace == "" {
		t.Error("Keyspace not found")
	}
}

func TestGetMerkleKeyspace(t *testing.T) {
	keyspace := util.GetMerkleKeyspace()

	if keyspace == "" {
		t.Error("Keyspace not found")
	}
}

func TestValidateKeyspace(t *testing.T) {
	valid, _ := util.ValidateKeyspace()

	if valid == false {
		t.Error("Invalid Keyspace Configuration")
	}
}

func TestGetTopics(t *testing.T) {
	topics := util.GetTopics()

	if len(*topics) < 3 {
		t.Error("Missing Topics")
	}
}

func TestGetACXPlusConfig(t *testing.T) {
	config := util.GetACXPlusConfig("acx-plus-metadata")

	if config.GroupID != metaGroupID {
		t.Errorf("Invalid Group ID %s", config.GroupID)
	}

	if config.Family != "metadata" {
		t.Errorf("Invalid Family %s", config.Family)
	}

	if config.Topic != metaTopic {
		t.Errorf("Invalid Topic %s", config.Topic)
	}

	if config.Instances != 1 {
		t.Errorf("Invalid Instance Count of %d", config.Instances)
	}

	if len(*config.KafkaHosts) != 3 {
		t.Error("Invalid Kafka Hosts")
	}

	if len(*config.CassandraHosts) != 3 {
		t.Error("Invalid Cassandra Hosts")
	}

	if config == nil {
		t.Error("Invalid or not found config")
	}
}

func TestValidateConnectivityOpen(t *testing.T) {
	port, _ := strconv.Atoi("9092")

	addrs := []net.TCPAddr{net.TCPAddr{
		IP:   net.ParseIP("172.22.4.61"),
		Port: port,
	}, net.TCPAddr{
		IP:   net.ParseIP("172.22.5.61"),
		Port: port,
	}, net.TCPAddr{
		IP:   net.ParseIP("172.22.6.61"),
		Port: port,
	}}

	conn := ValidateConnectivity(&addrs)

	if conn[0].Open == false {
		t.Error("Port is not open")
	}
}

func TestValidateConnectivityOpenCass(t *testing.T) {
	port, _ := strconv.Atoi("9042")

	addrs := []net.TCPAddr{net.TCPAddr{
		IP:   net.ParseIP("172.22.4.11"),
		Port: port,
	}, net.TCPAddr{
		IP:   net.ParseIP("172.22.5.11"),
		Port: port,
	}, net.TCPAddr{
		IP:   net.ParseIP("172.22.6.11"),
		Port: port,
	}}

	conn := ValidateConnectivity(&addrs)

	if conn[0].Open == false {
		t.Error("Port is not open")
	}
}
