package utils

import (
	"fmt"
	"os"

	ndmclientset "github.com/harvester/node-disk-manager/pkg/generated/clientset/versioned"
	lhclientset "github.com/longhorn/longhorn-manager/k8s/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getK8sClientsetInCluster() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getK8sClientsetLocal() (*kubernetes.Clientset, error) {
	configStr := os.Getenv("KUBECONFIG")
	if configStr == "" {
		return nil, fmt.Errorf("unable to get local config, please make sure the kubeconfig is set")
	}
	config, err := clientcmd.BuildConfigFromFlags("", configStr)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func GetK8sClientset() (*kubernetes.Clientset, error) {
	if clientset, err := getK8sClientsetInCluster(); err == nil {
		return clientset, nil
	}
	logrus.Debugf("Unable to get in cluster config, trying local config ...")
	if clientset, err := getK8sClientsetLocal(); err == nil {
		return clientset, nil
	}
	return nil, fmt.Errorf("unable to get local config, please make sure the kubeconfig is set")
}

func getLonghornClientsetInCluster() (*lhclientset.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return lhclientset.NewForConfig(config)
}

func getLonghornClientsetLocal() (*lhclientset.Clientset, error) {
	configStr := os.Getenv("KUBECONFIG")
	if configStr == "" {
		return nil, fmt.Errorf("unable to get local config, please make sure the kubeconfig is set")
	}
	config, err := clientcmd.BuildConfigFromFlags("", configStr)
	if err != nil {
		return nil, err
	}
	return lhclientset.NewForConfig(config)
}

func GetLonghornClientset() (*lhclientset.Clientset, error) {
	if clientset, err := getLonghornClientsetInCluster(); err == nil {
		return clientset, nil
	}
	logrus.Debugf("Unable to get in cluster config, trying local config ...")
	if clientset, err := getLonghornClientsetLocal(); err == nil {
		return clientset, nil
	}
	return nil, fmt.Errorf("unable to get local config, please make sure the kubeconfig is set")
}

func getNDMClientsetInCluster() (*ndmclientset.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return ndmclientset.NewForConfig(config)
}

func getNDMClientsetLocal() (*ndmclientset.Clientset, error) {
	configStr := os.Getenv("KUBECONFIG")
	if configStr == "" {
		return nil, fmt.Errorf("unable to get local config, please make sure the kubeconfig is set")
	}
	config, err := clientcmd.BuildConfigFromFlags("", configStr)
	if err != nil {
		return nil, err
	}
	return ndmclientset.NewForConfig(config)
}

func GetNDMClientset() (*ndmclientset.Clientset, error) {
	if clientset, err := getNDMClientsetInCluster(); err == nil {
		return clientset, nil
	}
	logrus.Debugf("Unable to get in cluster config, trying local config ...")
	if clientset, err := getNDMClientsetLocal(); err == nil {
		return clientset, nil
	}
	return nil, fmt.Errorf("unable to get local config, please make sure the kubeconfig is set")
}
