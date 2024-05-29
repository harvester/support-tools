package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	longhorn "github.com/longhorn/longhorn-manager/k8s/pkg/apis/longhorn/v1beta2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	diag "github.com/harvester/support-tools/libs/diag-client"
	diagutils "github.com/harvester/support-tools/libs/utils"
)

const (
	LONGHORN   = diagutils.PRODUCT_LONGHORN
	HARVESTER  = diagutils.PRODUCT_HARVESTER
	KUBERNETES = diagutils.PRODUCT_KUBERNETES
)

var replicaLocationCmd = &cobra.Command{
	Use:   "replica-location",
	Short: "Show the replicas and their corresponding phyiscal block devices.",
	Long:  `Show the replicas and their corresponding phyiscal block devices.`,
	Run: func(_ *cobra.Command, args []string) {
		if err := run(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	replicaLocationCmd.Flags().StringVar(&volumeName, "volume", os.Getenv("LONGHORN_VOLUME_NAME"), "The target volume name for the replica helper")
	replicaLocationCmd.Flags().StringVar(&namespace, "namespace", os.Getenv("LONGHORN_VOLUME_NAMESPACE"), "The target volume namespace for the replica helper")

	rootCmd.AddCommand(replicaLocationCmd)
}

func run() error {
	if volumeName == "" {
		return fmt.Errorf("volume name is required")
	}
	if namespace == "" {
		namespace = "longhorn-system"
	}
	logrus.Infof("replica-location with volume %s in namespace %s", volumeName, namespace)

	diagClient := diag.NewDiagClient("Harvester")
	allResources := diagClient.GetAllResourceList()
	params := map[string]string{}
	params["volume"] = volumeName
	params["namespace"] = namespace

	replicas, err := allResources[LONGHORN]["replicas"](diagClient.GetClient(LONGHORN), params)
	if err != nil {
		logrus.Errorf("Failed to get nodes: %v", err)
	}
	replicasOjb := replicas.(*longhorn.ReplicaList)
	fmt.Printf("Volume %s status as below:\n", volumeName)
	for _, replica := range replicasOjb.Items {
		if replica.Spec.DesireState != "running" {
			continue
		}
		if replica.Spec.VolumeName == volumeName {
			diskIDRaw := replica.Spec.DiskPath
			diskID := getDiskID(diskIDRaw)
			params["node"] = replica.Spec.NodeID
			fmt.Printf(" - Replica %s on %s disk in node %s\n", replica.Name, diskID, replica.Spec.NodeID)
			node, err := allResources[LONGHORN]["node"](diagClient.GetClient(LONGHORN), params)
			if err != nil {
				logrus.Errorf("Failed to get nodes: %v", err)
			}
			nodeObj := node.(*longhorn.Node)
			if _, found := nodeObj.Status.DiskStatus[diskID]; found {
				diskStatus := nodeObj.Status.DiskStatus[diskID]
				diskAvail := diskStatus.StorageAvailable
				diskMax := diskStatus.StorageMaximum
				diskScheduled := diskStatus.StorageScheduled
				percent := float64(diskMax-diskAvail) / float64(diskMax) * 100
				percentStr := fmt.Sprintf("%.2f%s", percent, "%")
				fmt.Printf("   - Disk %s, Path: %s \n", diskID, replica.Spec.DiskPath)
				fmt.Printf("   - Usage: %d/%d Bytes, %s used, %d Bytes are scheduled\n", diskMax-diskAvail, diskMax, percentStr, diskScheduled)
			}
			if diskID == "defaultdisk" {
				for name, diskStatus := range nodeObj.Status.DiskStatus {
					if strings.HasPrefix(name, "default-disk") {
						diskAvail := diskStatus.StorageAvailable
						diskMax := diskStatus.StorageMaximum
						diskScheduled := diskStatus.StorageScheduled
						percent := float64(diskMax-diskAvail) / float64(diskMax) * 100
						percentStr := fmt.Sprintf("%.2f%s", percent, "%")
						fmt.Printf("   - Disk %s, Path: %s \n", diskID, replica.Spec.DiskPath)
						fmt.Printf("   - Usage: %d/%d Bytes, %s used, %d Bytes are scheduled\n", diskMax-diskAvail, diskMax, percentStr, diskScheduled)
					}
				}
			}
		}
	}

	return nil
}

// the diskIDRaw is the disk path, might be looks like below:
//   - /var/lib/harvester/defaultdisk
//   - /var/lib/harvester/extra-disks/786b8467faa1bade048c7c040a02f185

func getDiskID(diskIDRaw string) string {
	return filepath.Base(diskIDRaw)
}
