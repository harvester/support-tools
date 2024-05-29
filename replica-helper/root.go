package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logDebug bool

	kubeConfigPath string

	volumeName string
	namespace  string
)

var rootCmd = &cobra.Command{
	Use:     "replica-helper",
	Short:   "Replica Helper",
	Long:    "A tool to show the replicas and their corresponding phyiscal block devices.",
	Version: AppVersion,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		logrus.SetOutput(os.Stdout)
		if logDebug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if kubeConfigPath != "" {
			os.Setenv("KUBECONFIG", kubeConfigPath)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&logDebug, "debug", false, "set logging level to debug")

	rootCmd.Flags().StringVar(&kubeConfigPath, "kubeconfig", os.Getenv("KUBECONFIG"), "Path to the kubeconfig file")
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}
