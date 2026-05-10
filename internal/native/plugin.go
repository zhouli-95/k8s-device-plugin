package native

import (
	"context"
	"fmt"

	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"tags.cncf.io/container-device-interface/specs-go"

	dpcdi "github.com/zhouli-95/k8s-device-plugin/internal/cdi"
	"github.com/zhouli-95/k8s-device-plugin/internal/config"
	"github.com/zhouli-95/k8s-device-plugin/internal/option"
	"github.com/zhouli-95/k8s-device-plugin/internal/plugin"
	"github.com/zhouli-95/k8s-device-plugin/pkg/resource"
)

var _ pluginapi.DevicePluginServer = &DevicePlugin{}
var _ plugin.Plugin = &DevicePlugin{}

type DevicePlugin struct {
	name       string
	ctx        context.Context
	devList    []*resource.ComputeDevice
	devMap     map[string]int
	healthChan <-chan struct{}
	cdiHandler *dpcdi.CDIHandler
	options    *option.Options
}

func NewDevicePlugin(ctx context.Context, ch <-chan struct{}, opts *option.Options) *DevicePlugin {
	p := DevicePlugin{
		name:       resource.GetResourceName(),
		ctx:        ctx,
		devList:    []*resource.ComputeDevice{},
		devMap:     make(map[string]int),
		healthChan: ch,
		options:    opts,
	}

	p.cdiHandler = dpcdi.NewCDIHandler(
		dpcdi.WithVendor(resource.Domain()),
		dpcdi.WithClass(p.name),
	)
	return &p
}

func (p *DevicePlugin) Name() string {
	return p.name
}

func (p *DevicePlugin) Start() error {
	klog.InfoS("Plugin start", "options", p.options)

	p.devList = resource.GetComputeDeviceList()
	for i, dev := range p.devList {
		p.devMap[dev.ID] = i
	}

	if err := p.createCDISpec(); err != nil {
		klog.Errorf("Failed to create cdi spec, %v", err)
		return err
	}

	return nil
}

func (p *DevicePlugin) createCDISpec() error {
	handler := func(spec *specs.Spec) error {
		for _, dev := range p.devList {
			devNodes := []*specs.DeviceNode{}
			for _, ctrDev := range resource.GetControlDeviceList() {
				devNodes = append(devNodes, &specs.DeviceNode{
					Path:        ctrDev.ContainerPath,
					HostPath:    ctrDev.HostPath,
					Permissions: "rw",
				})
			}

			devNodes = append(devNodes, &specs.DeviceNode{
				Path:        dev.ContainerPath,
				HostPath:    dev.HostPath,
				Permissions: "rw",
			})

			spec.Devices = append(spec.Devices, specs.Device{
				Name: dev.ID,
				ContainerEdits: specs.ContainerEdits{
					DeviceNodes: devNodes,
				},
			})
		}
		spec.ContainerEdits = specs.ContainerEdits{
			Env: []string{fmt.Sprintf("%s=%s", config.EnvVarDeviceStrategy, p.options.DeviceStrategy)},
		}
		return nil
	}

	if err := p.cdiHandler.CreateSpec(handler); err != nil {
		return err
	}

	return nil
}

func (p *DevicePlugin) Stop() error {
	klog.Infof("Plugin stop")
	return nil
}

func (p *DevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	options := &pluginapi.DevicePluginOptions{
		GetPreferredAllocationAvailable: true,
		PreStartRequired:                true,
	}
	return options, nil
}

func (p *DevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	klog.Info("List and watch start")
	defer klog.Infof("List and watch exit")

	if err := s.Send(&pluginapi.ListAndWatchResponse{Devices: []*pluginapi.Device{}}); err != nil {
		return err
	}

	for {
		select {
		case <-s.Context().Done():
			klog.Infof("ListAndWatch stream close: %v", s.Context().Err())
			return nil
		case <-p.ctx.Done():
			s.Send(&pluginapi.ListAndWatchResponse{Devices: []*pluginapi.Device{}})
			return nil
		case <-p.healthChan:
			if err := s.Send(&pluginapi.ListAndWatchResponse{Devices: p.genDevList()}); err != nil {
				return err
			}
		}
	}
}

func (p *DevicePlugin) genDevList() []*pluginapi.Device {
	devs := []*pluginapi.Device{}
	for _, dev := range p.devList {
		health := pluginapi.Healthy
		if !dev.Healthy() {
			health = pluginapi.Unhealthy
		}
		devs = append(devs, &pluginapi.Device{
			ID:     dev.ID,
			Health: health,
		})
	}
	return devs
}

func (p *DevicePlugin) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := pluginapi.AllocateResponse{}

	for _, req := range reqs.ContainerRequests {
		switch p.options.DeviceStrategy {
		case config.DeviceStrategyNative:
			resp, err := p.allocateByNative(req)
			if err != nil {
				return nil, err
			}
			responses.ContainerResponses = append(responses.ContainerResponses, resp)
		case config.DeviceStrategyCDICRI:
			resp, err := p.allocateByCDICRI(req)
			if err != nil {
				return nil, err
			}
			responses.ContainerResponses = append(responses.ContainerResponses, resp)

		case config.DeviceStrategyCDIAnnotation:
			resp, err := p.allocateByCDIAnnotation(req)
			if err != nil {
				return nil, err
			}
			responses.ContainerResponses = append(responses.ContainerResponses, resp)
		}
		klog.Infof("Allocated device %v by %s", req.DevicesIDs, p.options.DeviceStrategy)
	}

	return &responses, nil
}

func (p *DevicePlugin) allocateByNative(req *pluginapi.ContainerAllocateRequest) (*pluginapi.ContainerAllocateResponse, error) {
	resp := pluginapi.ContainerAllocateResponse{}

	resp.Envs = make(map[string]string)
	resp.Envs[config.EnvVarDeviceStrategy] = p.options.DeviceStrategy

	resp.Devices = []*pluginapi.DeviceSpec{}
	for _, ctrDev := range resource.GetControlDeviceList() {
		resp.Devices = append(resp.Devices, &pluginapi.DeviceSpec{
			ContainerPath: ctrDev.ContainerPath,
			HostPath:      ctrDev.HostPath,
			Permissions:   "rw",
		})
	}

	for _, id := range req.DevicesIDs {
		for _, dev := range p.devList {
			if dev.ID == id {
				resp.Devices = append(resp.Devices, &pluginapi.DeviceSpec{
					ContainerPath: dev.ContainerPath,
					HostPath:      dev.HostPath,
					Permissions:   "rw",
				})
			}
		}
	}

	return &resp, nil
}

func (p *DevicePlugin) allocateByCDIAnnotation(req *pluginapi.ContainerAllocateRequest) (*pluginapi.ContainerAllocateResponse, error) {
	resp := pluginapi.ContainerAllocateResponse{}

	var devices []string
	for _, id := range req.DevicesIDs {
		devices = append(devices, p.cdiHandler.QualifiedName(p.name, id))
	}

	annotations, err := p.cdiHandler.GetCDIAnnotation(devices)
	if err != nil {
		return nil, fmt.Errorf("failed to add cdi annotations: %w", err)
	}

	resp.Annotations = annotations
	return &resp, nil
}

func (p *DevicePlugin) allocateByCDICRI(req *pluginapi.ContainerAllocateRequest) (*pluginapi.ContainerAllocateResponse, error) {
	resp := pluginapi.ContainerAllocateResponse{}

	var devices []string
	for _, id := range req.DevicesIDs {
		devices = append(devices, p.cdiHandler.QualifiedName(p.name, id))
	}

	for _, id := range devices {
		dev := pluginapi.CDIDevice{
			Name: id,
		}
		resp.CDIDevices = append(resp.CDIDevices, &dev)
	}
	return &resp, nil
}

func (p *DevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	response := &pluginapi.PreferredAllocationResponse{}

	return response, nil
}

func (p *DevicePlugin) PreStartContainer(ctx context.Context, req *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	klog.Infof("Device %v reset", req.DevicesIDs)

	for _, id := range req.DevicesIDs {
		klog.Infof("Reset device %s", id)
		dev := p.devList[p.devMap[id]]
		dev.Reset()
	}

	return &pluginapi.PreStartContainerResponse{}, nil
}
