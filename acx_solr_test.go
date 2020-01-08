package k8sutil

import (
	"testing"
)

func TestGetSolrConfig(t *testing.T) {
	if util.Clientset == nil {
		clientset := BuildClientSet()

		util = ACXK8sUtil{
			Clientset: clientset,
			Namespace: namespace,
		}
	}

	product := "acx-common"
	config := util.GetSolrConfig(product, DefaultSolrFileName)

	if config.Config.Ribbon.ListOfServers[0] == "" {
		t.Error("Solr Hosts Not Found")
	}
}

func TestGetSolrConfigMissingFile(t *testing.T) {
	product := "acx-common"
	config := util.GetSolrConfig(product, "missing.properties")

	if config.Error.Error() == "" {
		t.Error("Incorrect File Found")
	}
}

func TestValidateSolrConnectivity(t *testing.T) {
	product := "acx-common"
	config := util.GetSolrConfig(product, DefaultSolrFileName)

	listOfServers := config.Config.Ribbon.ListOfServers

	result, connections := util.ValidateSolrConnectivity(listOfServers)

	if result != true {
		t.Error("Connection is Closed")
	}

	if len(connections) != 3 {
		t.Error("Invalid Connection Count")
	}
}

func TestValidateSolrConnectivityInvalid(t *testing.T) {
	product := "acx-common"
	config := util.GetSolrConfig(product, DefaultSolrFileName)

	listOfServers := config.Config.Ribbon.ListOfServers

	listOfServers = append(listOfServers, "172.22.4.11:6666")

	result, connections := util.ValidateSolrConnectivity(listOfServers)

	if result != false {
		t.Error("Connection is Closed")
	}

	if len(connections) != 4 {
		t.Error("Invalid Connection Count")
	}
}
