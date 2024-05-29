package diag

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	"github.com/harvester/support-tools/libs/resource"
	"github.com/harvester/support-tools/libs/utils"
)

type k8sDiagClient struct {
	client       *kubernetes.Clientset
	product      string
	resourceList map[string]resource.GatherFunc
}

// NewK8sDiagClient init then return the new K8sDiagClient
func NewK8sDiagClient() DiagClient {
	client := &k8sDiagClient{}
	client.product = utils.PRODUCT_KUBERNETES
	client.resourceList = map[string]resource.GatherFunc{}
	client.Init()
	return client
}

// GetProduct return the product name
func (k *k8sDiagClient) GetProduct() string {
	return k.product
}

// GetClient return the clientset based on the accepted type
func (k *k8sDiagClient) GetClient(clientType string) any {
	switch clientType {
	case utils.PRODUCT_KUBERNETES:
		return k.client
	}
	logrus.Errorf("Unsupported client type: %s", clientType)
	return nil
}

// Init of the k8SAnalyzer get the k8s clientset for other analyzer
func (k *k8sDiagClient) Init() {
	var err error
	k.client, err = utils.GetK8sClientset()
	if err != nil {
		logrus.Fatalf("Failed to get k8s clientset: %v", err)
		return
	}
	k.generateResourceList()
}

// GetResourceList return the internal resource list
func (k *k8sDiagClient) GetResourceList() map[string]resource.GatherFunc {
	return k.resourceList
}

// GetAllResourceList return the all resource list from all diag clients
func (h *k8sDiagClient) GetAllResourceList() map[string]map[string]resource.GatherFunc {
	allResources := make(map[string]map[string]resource.GatherFunc)
	allResources[h.product] = h.resourceList
	return allResources
}

func (k *k8sDiagClient) generateResourceList() {
	logrus.Debugf("Generate Kubernetes resource list (%p)", k.resourceList)
	k.resourceList["nodes"] = resource.GetNodeLists
}
