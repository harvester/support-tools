package diag

import (
	"github.com/sirupsen/logrus"

	"github.com/harvester/support-tools/libs/resource"
)

type DiagClient interface {
	Init()
	GetResourceList() map[string]resource.GatherFunc
	GetAllResourceList() map[string]map[string]resource.GatherFunc
	GetClient(string) any
	GetProduct() string

	generateResourceList()
}

func NewDiagClient(prod string) DiagClient {
	switch prod {
	case "Harvester":
		return NewHarvesterDiagClient(NewLonghornDiagClient(NewK8sDiagClient()))
	case "Longhorn":
		return NewLonghornDiagClient(NewK8sDiagClient())
	default:
		logrus.Errorf("Unsupported product: %s", prod)
		return nil
	}
}
