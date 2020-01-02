package k8sutil

import (
	"testing"
)

const namespace = "uat01"
const fileName = "application.conf"

var util ACXK8sUtil

func TestGetConfigMap(t *testing.T) {
	clientset := BuildClientSet()

	util = ACXK8sUtil{
		clientset: clientset,
		Namespace: namespace,
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

	if groupID == "" {
		t.Error("GroupID not found")
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

// func TestValidateGroupID(t *testing.T) {
// 	mode := "init"
// 	valid := util.ValidateGroupID(mode)

// 	if valid == false {
// 		t.Error("Incorrect GroupID configuration")
// 	}
// }
