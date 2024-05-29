package resource

import (
	"context"

	ndmclientset "github.com/harvester/node-disk-manager/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func listBlockdevices(client *ndmclientset.Clientset, _ parameters) (interface{}, error) {
	return client.HarvesterhciV1beta1().BlockDevices("longhorn-system").List(context.TODO(), metav1.ListOptions{})
}
