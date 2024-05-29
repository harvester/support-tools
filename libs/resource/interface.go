package resource

import (
	ndmclientset "github.com/harvester/node-disk-manager/pkg/generated/clientset/versioned"
	lhclientset "github.com/longhorn/longhorn-manager/k8s/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type parameters = map[string]string

type GatherFunc func(client any, args parameters) (interface{}, error)

// GetNodes returns the node list of the cluster
// Parameters:
//
//	client: k8s client
//	args[0]: envType
func GetNodeLists(client any, args parameters) (interface{}, error) {
	logrus.Debugf("Calling GetNodes with args: %v", args)
	k8sClient := client.(*kubernetes.Clientset)
	return listNodes(k8sClient, args)
}

func GetReplicaLists(client any, args parameters) (interface{}, error) {
	logrus.Debugf("Calling GetReplicas with args: %v", args)
	lhClient := client.(*lhclientset.Clientset)
	return listReplicas(lhClient, args)
}

func GetVolume(client any, args parameters) (interface{}, error) {
	logrus.Infof("Calling GetVolume with args: %v", args)
	lhClient := client.(*lhclientset.Clientset)
	return getVolume(lhClient, args)
}

func GetNode(client any, args parameters) (interface{}, error) {
	logrus.Debugf("Calling GetNode with args: %v", args)
	lhClient := client.(*lhclientset.Clientset)
	return getNode(lhClient, args)
}

func GetBlockdeviceLists(client any, args parameters) (interface{}, error) {
	logrus.Debugf("Calling GetBlockdevices with args: %v", args)
	ndmClient := client.(*ndmclientset.Clientset)
	return listBlockdevices(ndmClient, args)
}
