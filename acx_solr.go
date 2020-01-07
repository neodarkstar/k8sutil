package k8sutil

import (
	"bytes"
	"strconv"

	"net"

	"github.com/spf13/viper"
)

// DefaultSolrFileName defaults to solr.properties
const DefaultSolrFileName = "solr.properties"

// Declare Types
type solrConfigError struct {
	Message string
}

func (s *solrConfigError) Error() string {
	return s.Message
}

// SolrConfig A shortened list of Solr Properties File Heirarchy
type SolrConfig struct {
	Config Solr
	Error  solrConfigError
}

// Solr A shortened list of Solr Properties File Heirarchy
type Solr struct {
	Ribbon Ribbon
}

// Ribbon A shortened list of Solr Properties File Heirarchy
type Ribbon struct {
	MaxAutoRestries           int
	MaxAutoRetriesNextServer  int
	OkToRetryOnAllOperations  bool
	ServerListRefreshInterval int
	ListOfServers             []string
	ClientClassName           string
}

// GetSolrConfig Returns a java properties file
func (k *ACXK8sUtil) GetSolrConfig(product string, fileName string) *SolrConfig {
	viper.SetConfigType("properties")
	configMap := k.GetConfigMap(product, fileName)

	var S SolrConfig = SolrConfig{}

	if configMap == "" {
		S.Error = solrConfigError{"Missing ConfigMap File : " + fileName}
		return &S
	}

	viper.ReadConfig(bytes.NewBufferString(configMap))
	err := viper.UnmarshalKey("solr", &S.Config)
	if err != nil {
		S.Error = solrConfigError{"unable to decode into struct"}
	}

	return &S
}

// ValidateSolrConnectivity checks if the port is open and accessible
func (k *ACXK8sUtil) ValidateSolrConnectivity(listOfServers []string) (bool, []Connection) {
	addrs := make([]net.TCPAddr, 0)

	for _, server := range listOfServers {
		h, p, _ := net.SplitHostPort(server)
		port, _ := strconv.Atoi(p)
		addrs = append(addrs, net.TCPAddr{
			IP:   net.ParseIP(h),
			Port: port,
		})
	}

	connections := ValidateConnectivity(&addrs)

	for _, conn := range connections {
		if conn.Error != nil {
			return false, connections
		}
	}

	return true, connections
}
