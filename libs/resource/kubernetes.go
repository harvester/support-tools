package resource

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func listNodes(client *kubernetes.Clientset, _ parameters) (interface{}, error) {
	return client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}
