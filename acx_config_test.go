package k8sutil

import "testing"

func TestGetJavaPropertesConfig(t *testing.T) {
	product := "acx-plus-metadata"
	config := util.GetJavaPropertiesConfig(product, "solr")

	if config.Ribbon.ListOfServers == "" {
		t.Error("Solr Hosts Not Found")
	}
}
