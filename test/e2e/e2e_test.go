/*
Copyright 2026 zhouli-95

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
