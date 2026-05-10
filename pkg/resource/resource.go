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
