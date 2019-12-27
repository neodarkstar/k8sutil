package k8sutil

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-akka/configuration"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func buildClientSet() *kubernetes.Clientset {
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

func downloadConfigMaps(namespace string, product string) error {
	clientset := buildClientSet()

	configmap, _ := clientset.CoreV1().ConfigMaps(namespace).Get(product, metav1.GetOptions{})

	for fileName := range configmap.Data {
		writeToFile(fileName, configmap.Data[fileName])
	}

	return nil
}

func writeToFile(fileName string, config string) {
	d1 := []byte(config)
	err := ioutil.WriteFile("./acx-plus/"+fileName, d1, 0644)
	check(err)
}

func parseHOCON(config string) {
	conf := configuration.ParseString(config)

	hconvalue := conf.GetConfig("acx-plus-tm.node-group.node-families").Root().GetArray()[0]

	kafkaConfig := hconvalue.GetChildObject("kafka").GetObject()

	fmt.Println(kafkaConfig.GetKey("topic").GetString())
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
