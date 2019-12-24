package k8sutil

import (
	"testing"
)

func TestRead(t *testing.T) {
	namespace := "acx"
	err := readConfigMap(&namespace)

	if err != nil {
		t.Error("Got an Error")
	}
}

// func TestDropKeyspace(t *testing.T) {
// 	_, err := dropKeyspace("unit_tests")
// 	if err != nil {
// 		t.Error("Got an Error", err)
// 	}
// }
