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

package utils

import (
	"fmt"
	"path"
)

var defaultKustomizeDir = "../../deployments/kustomize/base"
var kustomizeBaseDir = "../kustomize/overlays"

type DeployOption func(*DeployConfig)

type DeployConfig struct {
	KustomizeDir string
}

func WithKustomizeDir(dir string) DeployOption {
	return func(c *DeployConfig) {
		c.KustomizeDir = dir
	}
}

func Deploy(options ...DeployOption) error {
	config := &DeployConfig{}
	for _, option := range options {
		option(config)
	}

	kustomizeDir := defaultKustomizeDir
	if config.KustomizeDir != "" {
		kustomizeDir = path.Join(kustomizeBaseDir, config.KustomizeDir)
	}

	cmd := fmt.Sprintf("kubectl apply -k %s", kustomizeDir)
	return BashRun(cmd)
}

func UnDeploy(options ...DeployOption) error {
	config := &DeployConfig{}
	for _, option := range options {
		option(config)
	}

	kustomizeDir := defaultKustomizeDir
	if config.KustomizeDir != "" {
		kustomizeDir = path.Join(kustomizeBaseDir, config.KustomizeDir)
	}

	cmd := fmt.Sprintf("kubectl delete -k %s", kustomizeDir)
	return BashRun(cmd)
}
