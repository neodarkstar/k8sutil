package k8sutil

import (
	"github.com/go-akka/configuration"
	"github.com/go-akka/configuration/hocon"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetConfigMap Retrieves a ConfigMap File
func (k *ACXK8sUtil) GetConfigMap(product string, fileName string) string {
	namespace := k.Namespace
	clientset := k.Clientset
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
