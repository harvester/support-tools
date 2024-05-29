package diag

import (
	"github.com/sirupsen/logrus"

	"github.com/harvester/support-tools/libs/resource"
	"github.com/harvester/support-tools/libs/utils"
)

type HarvesterDiagClient struct {
	harvesterClients   map[string]any
	LonghornDiagClient DiagClient
	K8sDiagClient      DiagClient
	product            string
	resourceList       map[string]resource.GatherFunc
	DiagClients        []DiagClient
}

// NewHarvesterDiagClient init then return the new HarvesterDiagClient
func NewHarvesterDiagClient(client DiagClient) DiagClient {
	harvclient := &HarvesterDiagClient{}
	// Harvester should not have more than 8 clientsets
	harvclient.harvesterClients = make(map[string]any, 8)
	harvclient.LonghornDiagClient = client
	harvclient.K8sDiagClient = client.(*LonghornDiagClient).K8sDiagClient
	harvclient.product = utils.PRODUCT_HARVESTER
	harvclient.resourceList = map[string]resource.GatherFunc{}
	harvclient.DiagClients = []DiagClient{harvclient.LonghornDiagClient, harvclient.K8sDiagClient}
	harvclient.Init()
	return harvclient
}

// GetProduct return the product name
func (h *HarvesterDiagClient) GetProduct() string {
	return h.product
}

// GetClient return the clientset based on the accepted type
func (h *HarvesterDiagClient) GetClient(clientType string) any {
	switch clientType {
	case utils.PRODUCT_KUBERNETES:
		return h.K8sDiagClient.GetClient(clientType)
	case utils.PRODUCT_LONGHORN:
		return h.LonghornDiagClient.GetClient(clientType)
	case utils.PRODUCT_HARVESTER:
		return h.harvesterClients[clientType]
	case utils.PRODUCT_HARV_NDM:
		return h.harvesterClients[clientType]
	}
	logrus.Errorf("Unsupported client type: %s", clientType)
	return nil
}

// Init of the k8SAnalyzer get the k8s clientset for other analyzer
func (h *HarvesterDiagClient) Init() {
	logrus.Debugf("Init HarvesterDiagClient")
	h.initClients()
	h.generateResourceList()
}

// GetResourceList return the internal resource list
func (h *HarvesterDiagClient) GetResourceList() map[string]resource.GatherFunc {
	return h.resourceList
}

// GetAllResourceList return the all resource list from all diag clients
func (h *HarvesterDiagClient) GetAllResourceList() map[string]map[string]resource.GatherFunc {
	allResources := make(map[string]map[string]resource.GatherFunc)
	for _, client := range h.DiagClients {
		allResources[client.GetProduct()] = client.GetResourceList()
	}
	return allResources
}

func (h *HarvesterDiagClient) generateResourceList() {
	logrus.Debugf("Generate Harvester resource list (%p)", h.resourceList)
	h.resourceList["blockdevices"] = resource.GetBlockdeviceLists
}

func (h *HarvesterDiagClient) initClients() {
	var err error
	h.harvesterClients[utils.PRODUCT_HARV_NDM], err = utils.GetNDMClientset()
	if err != nil {
		logrus.Fatalf("Failed to get k8s clientset: %v", err)
		return
	}
}
