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

package resource

import (
	"encoding/json"
	"fmt"
	"os"

	"k8s.io/klog/v2"
)

var config *ResourceConfig

type ResourceConfig struct {
	Domain   string   `json:"domain"`
	Resource Resource `json:"resource"`
}

type Resource struct {
	Name           string           `json:"name"`
	ControlDevices []*ControlDevice `json:"controlDevices,omitempty"`
	ComputeDevices []*ComputeDevice `json:"computeDevices,omitempty"`
}

type ControlDevice struct {
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath,omitempty"`
}

type ComputeDevice struct {
	ID            string `json:"id,omitempty"`
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath,omitempty"`
}

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func Domain() string {
	return config.Domain
}

func GetResourceName() string {
	return config.Resource.Name
}

func GetControlDeviceList() []*ControlDevice {
	return config.Resource.ControlDevices
}

func GetComputeDeviceList() []*ComputeDevice {
	return config.Resource.ComputeDevices
}

func (d *ComputeDevice) Healthy() bool {
	return true
}

func (d *ComputeDevice) Reset() error {
	klog.Infof("Device %s reset", d.ID)
	return nil
}
