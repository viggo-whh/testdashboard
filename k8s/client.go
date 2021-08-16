package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type KubeClient struct {
	Pod  *PodClient
	Node *NodeClient
}

var Client *KubeClient

func InitK8sClientSet() (*kubernetes.Clientset, *rest.Config, error) {
	var err error
	var config *rest.Config
	var clientset *kubernetes.Clientset

	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return nil, nil, err
		}
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return clientset, config, nil
}

func NewKubeClient() error {
	client, _, err := InitK8sClientSet()
	if err != nil {
		return err
	}
	Client = &KubeClient{
		Pod:  NewPodClient(client),
		Node: NewNodeClient(client),
	}
	return nil
}
