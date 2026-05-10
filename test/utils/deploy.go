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
