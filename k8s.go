package k8sutil

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-akka/configuration"
	"github.com/go-akka/configuration/hocon"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ACXK8sUtil contains the methods to access kubernetes configurations and verify cross suite
type ACXK8sUtil struct {
	clientset *kubernetes.Clientset
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

// GetConfigMap Retrieves a ConfigMap File
func (k *ACXK8sUtil) GetConfigMap(product string, fileName string) string {
	namespace := k.Namespace
	clientset := k.clientset
	configmap, _ := clientset.CoreV1().ConfigMaps(namespace).Get(product, metav1.GetOptions{})
	return configmap.Data[fileName]
}

// GetConfig Returns a HOCON based Configurations from a Kubernetes ConfigMap
func (k *ACXK8sUtil) GetConfig(product string, fileName string) *configuration.Config {
	configMap := k.GetConfigMap(product, fileName)

	conf := configuration.ParseString(configMap, func(internalFileName string) *hocon.HoconRoot {
		includeConfigMap := k.GetConfigMap(product, internalFileName)

		return configuration.ParseString(includeConfigMap).Root().AtKey(internalFileName)
	})

	return conf
}

func (k *ACXK8sUtil) getKafkaConfig(product string) *hocon.HoconObject {
	const nodeFamilyPath = "acx-plus-tm.node-group.node-families"
	conf := k.GetConfig(product, "application.conf")
	hconvalue := conf.GetConfig(nodeFamilyPath).Root().GetArray()[0]
	kafkaConfig := hconvalue.GetChildObject("kafka").GetObject()

	return kafkaConfig
}

// GetGroupID retrieves the KafkaGroupID associated with the passed in product
func (k *ACXK8sUtil) GetGroupID(product string) string {
	kafkaConfig := k.getKafkaConfig(product)

	return kafkaConfig.GetKey("topic").GetString()
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
	conf := k.GetConfig(product, "core.conf")

	keyspace := conf.GetString(nodePath)

	return keyspace
}

// GetMerkleKeyspace Retrieves the Configured Keyspace
func (k *ACXK8sUtil) GetMerkleKeyspace() string {
	const nodePath = "services.merkleTree.properties.keyspace"
	fileName := "application.conf"
	product := "acx-merkle-rest"
	conf := k.GetConfig(product, fileName)

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

	fmt.Println(topics)

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
