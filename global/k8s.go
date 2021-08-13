package global

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var (
	client *kubernetes.Clientset
)

func InitK8sClientSet() error  {
	var err error
	var config *rest.Config

	kubeconfig := filepath.Join(homedir.HomeDir(),".kube","config")
	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("",kubeconfig); err != nil {
			return err
		}
	}

	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func K8sClient() *kubernetes.Clientset{
	return client
}

