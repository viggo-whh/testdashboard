package k8s

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeClient struct {
	clientset *kubernetes.Clientset
	//config *rest.Config
}

func NewNodeClient(clientset *kubernetes.Clientset) *NodeClient {
	return &NodeClient{clientset: clientset}
}

func (nc *NodeClient) List(labels string) ([]corev1.Node, error) {
	opts := metav1.ListOptions{}
	if labels != "" {
		opts.LabelSelector = labels
	}

	var nodelist, err = nc.clientset.CoreV1().Nodes().List(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	return nodelist.Items, nil

}
