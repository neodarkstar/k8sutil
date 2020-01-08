package k8sutil

import (
	"flag"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-akka/configuration/hocon"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ACXK8sUtil contains the methods to access kubernetes configurations and verify cross suite
type ACXK8sUtil struct {
	Clientset *kubernetes.Clientset
	Namespace string
}

// BuildClientSet Connects to Kubernetes
func BuildClientSet() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	check(err)

	return clientset
}

func (k *ACXK8sUtil) getKafkaConfig(product string) *hocon.HoconObject {
	const nodeFamilyPath = "acx-plus-tm.node-group.node-families"
	conf := k.GetHOCONConfig(product, "application.conf")
	hconvalue := conf.GetConfig(nodeFamilyPath).Root().GetArray()[0]
	kafkaConfig := hconvalue.GetChildObject("kafka").GetObject()

	return kafkaConfig
}

// GetACXPlusConfig Retrieves th ACXPlusConfig from Kubernetes
func (k *ACXK8sUtil) GetACXPlusConfig(product string) *ACXPlusNodeConfig {
	const nodeFamilyPath = "acx-plus-tm.node-group.node-families"
	conf := k.GetHOCONConfig(product, "application.conf")
	nodeObjects := conf.GetConfig(nodeFamilyPath).Root().GetArray()

	if len(nodeObjects) == 0 {
		panic("Missing Node Family Configuration")
	}

	nodeObject := nodeObjects[0]

	kafkaObject := nodeObject.GetChildObject("kafka")
	KafkaBootStrapServers := strings.Split(kafkaObject.GetChildObject("bootstrap").GetChildObject("servers").GetString(), ",")

	kafkaHosts := make([]net.TCPAddr, 0)

	for _, k := range KafkaBootStrapServers {
		host, p, splitError := net.SplitHostPort(k)

		if splitError != nil {
			panic("Could Not Split Kafka Host")
		}

		port, strconvError := strconv.Atoi(p)

		if strconvError != nil {
			panic("Invalid Port Number. Failed to Parse")
		}

		kafkaHosts = append(kafkaHosts, net.TCPAddr{
			IP:   net.ParseIP(host),
			Port: port,
		})
	}

	// Cassandra Hosts
	const cassandraPort string = "9042"
	coreConfig := k.GetHOCONConfig(product, "core.conf")
	cassandraConfHosts := strings.Split(coreConfig.GetString("CassandraConf.hosts"), ",")

	cassandraHosts := make([]net.TCPAddr, 0)

	port, _ := strconv.Atoi(cassandraPort)

	for _, host := range cassandraConfHosts {
		cassandraHosts = append(cassandraHosts, net.TCPAddr{
			IP:   net.ParseIP(host),
			Port: port,
		})
	}

	// SolrHosts

	return &ACXPlusNodeConfig{
		App:            product,
		GroupID:        kafkaObject.GetChildObject("group").GetChildObject("id").GetString(),
		Family:         nodeObject.GetChildObject("family").GetString(),
		Topic:          kafkaObject.GetChildObject("topic").GetString(),
		Instances:      nodeObject.GetChildObject("instances").GetInt32(),
		KafkaHosts:     &kafkaHosts,
		CassandraHosts: &cassandraHosts,
	}
}

// GetGroupID retrieves the KafkaGroupID associated with the passed in product
func (k *ACXK8sUtil) GetGroupID(product string) string {
	kafkaConfig := k.getKafkaConfig(product)

	return kafkaConfig.GetKey("group").GetObject().GetKey("id").GetString()
}

// ValidateGroupID Retrieves GroupID from ACX Products and verifies they are the same and in correct mode
// supported modes are init|update
func (k *ACXK8sUtil) ValidateGroupID(mode string) bool {
	products := []string{"acx-plus-metadata", "acx-plus-txn", "acx-plus-derived"}

	groups := make([]string, 0)

	for _, p := range products {
		groupID := k.GetGroupID(p)
		groups = append(groups, groupID)
	}

	return false
}

// GetKeyspace Retrieves the configurated Keyspace for a Product
func (k *ACXK8sUtil) GetKeyspace(product string) string {
	const nodePath = "CassandraConf.keyspace"
	conf := k.GetHOCONConfig(product, "core.conf")

	keyspace := conf.GetString(nodePath)

	return keyspace
}

// GetMerkleKeyspace Retrieves the Configured Keyspace
func (k *ACXK8sUtil) GetMerkleKeyspace() string {
	const nodePath = "services.merkleTree.properties.keyspace"
	fileName := "application.conf"
	product := "acx-merkle-rest"
	conf := k.GetHOCONConfig(product, fileName)

	keyspace := conf.GetString(nodePath)

	return keyspace
}

// ValidateKeyspace Checks all products for a consistent Keyspace and Returns the found keyspace
func (k *ACXK8sUtil) ValidateKeyspace() (bool, string) {
	products := []string{"acx-plus-metadata", "acx-plus-txn", "acx-plus-derived", "acx-rest", "acx-streaming-server"}
	keyspaces := make([]string, 0)

	for _, p := range products {
		keyspace := k.GetKeyspace(p)

		if keyspace == "" {
			return false, keyspace
		}
		keyspaces = append(keyspaces, keyspace)
	}

	return true, keyspaces[0]
}

// GetTopic Retrieves the Topic
func (k *ACXK8sUtil) GetTopic(product string) string {
	kafkaConfig := k.getKafkaConfig(product)
	return kafkaConfig.GetKey("topic").GetString()
}

// GetTopics Returns the Topics for All ACX Products
func (k *ACXK8sUtil) GetTopics() *[]string {
	products := []string{"acx-plus-metadata", "acx-plus-txn", "acx-plus-derived"}
	topics := make([]string, 0)

	for _, product := range products {
		topic := k.GetTopic(product)
		topics = append(topics, topic)
	}

	return &topics
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// ValidateConnectivity Checks for Open Ports
func ValidateConnectivity(addrs *[]net.TCPAddr) []Connection {
	connections := make([]Connection, 0)

	for _, addr := range *addrs {
		open := true
		network := addr.Network()
		host := addr.String()
		conn, error := net.Dial(network, host)
		defer func() {
			if conn != nil {
				conn.Close()
			}
		}()

		if error != nil {
			open = false
		}

		connections = append(connections, Connection{
			Open:  open,
			Addr:  addr,
			Error: error,
		})
	}

	return connections
}
