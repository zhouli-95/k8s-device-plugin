package e2e

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeConfig string
var kubeClient kubernetes.Interface

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	suiteConfig, reporterConfig := ginkgo.GinkgoConfiguration()
	// Randomize specs as well as suites
	// suiteConfig.RandomizeAllSpecs = true

	initKubeConfig()

	ginkgo.RunSpecs(t, "k8s-device-plugin e2e suite", suiteConfig, reporterConfig)
}

func initKubeConfig() {
	kubeConfig = os.Getenv("KUBECONFIG")
	if kubeConfig == "" {
		kubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	csconfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("Failed to build kubeconfig, %v", err)
	}
	client, err := kubernetes.NewForConfig(csconfig)
	if err != nil {
		log.Fatalf("Failed to new config, %v", err)
	}

	fmt.Printf("Using kubeconfig: %s\n", kubeConfig)

	kubeClient = client
}
