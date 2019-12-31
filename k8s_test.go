package k8sutil

import (
	"testing"
)

const namespace = "acx"

func TestParseConfig(t *testing.T) {

}

// func TestGetConfigMap(t *testing.T) {
// 	product := "acx-plus-metadata"
// 	getConfigMap(namespace, product)
// }

// func TestParseHOCON(t *testing.T) {
// 	product := "acx-plus-metadata"
// 	config := getConfigMap(namespace, product)

// 	parseHOCON(&config)
// }

func TestGetGroupID(t *testing.T) {
	product := "acx-plus-metadata"
	config := getConfigMap(namespace, product)

	conf := parseHOCON(&config)

	groupID := getGroupID(product, conf)

	if groupID != "santander50k_1_john_MD" {
		t.Error("Incorrect Group ID")
	}
}

// func TestDropKeyspace(t *testing.T) {
// 	_, err := dropKeyspace("unit_tests")
// 	if err != nil {
// 		t.Error("Got an Error", err)
// 	}
// }
