package config

const (
	PluginModeNative = "native"
	PluginModeShare  = "share"
)

const (
	DeviceStrategyNative        = "native"
	DeviceStrategyCDICRI        = "cdi-cri"
	DeviceStrategyCDIAnnotation = "cdi-annotation"
)

const EnvVarDeviceStrategy = "DEVICE_STRATEGY"

const DefaultConfigPath = "/config/config.json"
