package k8sutil

// PlusConfig ACX Plus ConfigMap Wrapper
type PlusConfig struct {
	Topic   string
	GroupID string
	Family  string
}

// CoreConfig ACX Common ConfigMap Wrapper
type CoreConfig struct {
}

func (Config *PlusConfig) parseHOCON(conf string) {

}
