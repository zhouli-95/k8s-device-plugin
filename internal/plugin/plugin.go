package plugin

import (
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type Plugin interface {
	Name() string
	Start() error
	Stop() error
	pluginapi.DevicePluginServer
}
