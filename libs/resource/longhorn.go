package resource

import (
	"context"
	"fmt"

	lhclientset "github.com/longhorn/longhorn-manager/k8s/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func listReplicas(client *lhclientset.Clientset, params parameters) (interface{}, error) {
	ns := "longhorn-system"
	if _, found := params["namespace"]; found {
		ns = params["namespace"]
	}
	return client.LonghornV1beta2().Replicas(ns).List(context.TODO(), metav1.ListOptions{})
}

func getVolume(client *lhclientset.Clientset, params parameters) (interface{}, error) {
	ns := "longhorn-system"
	if _, found := params["namespace"]; found {
		ns = params["namespace"]
	}
	if _, found := params["volume"]; !found {
		return nil, fmt.Errorf("volume name is required")
	}
	name := params["volume"]

	return client.LonghornV1beta2().Volumes(ns).Get(context.TODO(), name, metav1.GetOptions{})
}

func getNode(client *lhclientset.Clientset, params parameters) (interface{}, error) {
	ns := "longhorn-system"
	if _, found := params["namespace"]; found {
		ns = params["namespace"]
	}
	if _, found := params["node"]; !found {
		return nil, fmt.Errorf("node name is required")
	}
	name := params["node"]
	return client.LonghornV1beta2().Nodes(ns).Get(context.TODO(), name, metav1.GetOptions{})
}
