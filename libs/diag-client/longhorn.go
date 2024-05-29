package diag

import (
	lhclientset "github.com/longhorn/longhorn-manager/k8s/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"

	"github.com/harvester/support-tools/libs/resource"
	"github.com/harvester/support-tools/libs/utils"
)

type LonghornDiagClient struct {
	client       *lhclientset.Clientset
	product      string
	resourceList map[string]resource.GatherFunc

	K8sDiagClient DiagClient
	DiagClients   []DiagClient
}

// NewLonghornDiagClient init then return the new LonghornDiagClient
func NewLonghornDiagClient(client DiagClient) DiagClient {
	lhclient := &LonghornDiagClient{}
	lhclient.K8sDiagClient = client
	lhclient.product = utils.PRODUCT_LONGHORN
	lhclient.resourceList = map[string]resource.GatherFunc{}
	lhclient.DiagClients = []DiagClient{lhclient.K8sDiagClient}
	lhclient.Init()
	return lhclient
}

// GetProduct return the product name
func (l *LonghornDiagClient) GetProduct() string {
	return l.product
}

// GetClient return the clientset based on the accepted type
func (l *LonghornDiagClient) GetClient(clientType string) any {
	switch clientType {
	case utils.PRODUCT_KUBERNETES:
		return l.K8sDiagClient.GetClient(clientType)
	case utils.PRODUCT_LONGHORN:
		return l.client
	}
	logrus.Errorf("Unsupported client type: %s", clientType)
	return nil
}

// Init of the k8SAnalyzer get the k8s clientset for other analyzer
func (l *LonghornDiagClient) Init() {
	logrus.Debugf("Init LonghornDiagClient")
	var err error
	l.client, err = utils.GetLonghornClientset()
	if err != nil {
		logrus.Fatalf("Failed to get k8s clientset: %v", err)
		return
	}
	l.generateResourceList()
}

// GetResourceList return the internal resource list
func (l *LonghornDiagClient) GetResourceList() map[string]resource.GatherFunc {
	return l.resourceList
}

// GetAllResourceList return the all resource list from all diag clients
func (h *LonghornDiagClient) GetAllResourceList() map[string]map[string]resource.GatherFunc {
	allResources := make(map[string]map[string]resource.GatherFunc)
	for _, client := range h.DiagClients {
		allResources[client.GetProduct()] = client.GetResourceList()
	}
	return allResources
}

func (l *LonghornDiagClient) generateResourceList() {
	logrus.Debugf("Generate Longhorn resource list (%p)", l.resourceList)
	l.resourceList["replicas"] = resource.GetReplicaLists
	l.resourceList["volume"] = resource.GetVolume
	l.resourceList["node"] = resource.GetNode
}
