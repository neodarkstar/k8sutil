package k8sutil

import (
	"fmt"

	"github.com/go-akka/configuration"
	"github.com/go-akka/configuration/hocon"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetConfigMap Retrieves a ConfigMap File
func (k *ACXK8sUtil) GetConfigMap(product string, fileName string) string {
	namespace := k.Namespace
	clientset := k.clientset
	configmap, _ := clientset.CoreV1().ConfigMaps(namespace).Get(product, metav1.GetOptions{})
	return configmap.Data[fileName]
}

// GetHOCONConfig Returns a HOCON based Configurations from a Kubernetes ConfigMap
func (k *ACXK8sUtil) GetHOCONConfig(product string, fileName string) *configuration.Config {
	configMap := k.GetConfigMap(product, fileName)

	conf := configuration.ParseString(configMap, func(internalFileName string) *hocon.HoconRoot {
		includeConfigMap := k.GetConfigMap(product, internalFileName)

		return configuration.ParseString(includeConfigMap).Root().AtKey(internalFileName)
	})

	return conf
}

// GetJavaPropertiesConfig Returns a java properties file
func (k *ACXK8sUtil) GetJavaPropertiesConfig(product string, fileName string) *SolrConfig {
	viper.SetConfigName(fileName)
	viper.SetConfigType("properties")
	viper.AddConfigPath("./acx-plus/")

	viper.ReadInConfig()

	c := viper.AllSettings()
	bs, _ := yaml.Marshal(c)

	fmt.Println(string(bs))

	var S SolrConfig

	err := viper.UnmarshalKey("solr", &S)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}

	fmt.Println(S.Ribbon.ListOfServers)

	return &S
}

type SolrConfig struct {
	Ribbon Ribbon
}

type Ribbon struct {
	MaxAutoRestries           int
	MaxAutoRetriesNextServer  int
	OkToRetryOnAllOperations  bool
	ServerListRefreshInterval int
	ListOfServers             string
	ClientClassName           string
	Eureka                    Eureka
}

type Eureka struct {
	enabled bool
}
