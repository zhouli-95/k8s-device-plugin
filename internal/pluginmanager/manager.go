package pluginmanager

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/zhouli-95/k8s-device-plugin/internal/config"
	"github.com/zhouli-95/k8s-device-plugin/internal/native"
	"github.com/zhouli-95/k8s-device-plugin/internal/option"
	"github.com/zhouli-95/k8s-device-plugin/internal/plugin"
	"github.com/zhouli-95/k8s-device-plugin/internal/share"
	"github.com/zhouli-95/k8s-device-plugin/pkg/resource"
)

var HealthInterval = 10 * time.Second

type PluginManager struct {
	ctx        context.Context
	plugins    map[string]plugin.Plugin
	servers    map[string]*grpc.Server
	healthChan chan struct{}
	options    *option.Options
}

func NewPluginManager(ctx context.Context, opts *option.Options) *PluginManager {
	return &PluginManager{
		plugins:    make(map[string]plugin.Plugin),
		servers:    make(map[string]*grpc.Server),
		ctx:        ctx,
		healthChan: make(chan struct{}),
		options:    opts,
	}
}

func (p *PluginManager) Start() error {
	klog.Infof("Plugin manager start")
	p.initPlugins()

	errs := []error{}
	for name, plugin := range p.plugins {
		if err := plugin.Start(); err != nil {
			klog.Errorf("Failed to start plugin %s: %v", name, err)
			errs = append(errs, err)
			continue
		}

		server, err := StartServer(plugin.Name(), plugin)
		if err != nil {
			klog.Errorf("Failed to start server for plugin %s: %v", name, err)
			errs = append(errs, err)
			continue
		}
		p.servers[name] = server

		klog.Info("resource name", "domain", resource.Domain(), "name", plugin.Name())

		if err := RegisterPlugin(resource.Domain(), plugin.Name()); err != nil {
			klog.Errorf("Failed to register plugin %v", err)
			errs = append(errs, err)
			continue
		}
	}

	go func() {
		ticker := time.NewTicker(HealthInterval)
		defer ticker.Stop()
		klog.InfoS("Started health check ticker", "interval", HealthInterval)
		for {
			select {
			case <-ticker.C:
				p.healthChan <- struct{}{}
			case <-p.ctx.Done():
				klog.Info("Health ticker exited")
				return
			}
		}
	}()

	return errors.Join(errs...)
}

func (p *PluginManager) initPlugins() {
	mode := p.options.PluginMode
	switch mode {
	case config.PluginModeNative:
		plugin := native.NewDevicePlugin(
			p.ctx,
			p.healthChan,
			p.options,
		)
		p.plugins[mode] = plugin
	case config.PluginModeShare:
		plugin := share.NewDevicePlugin(
			p.ctx,
			p.healthChan,
			p.options,
		)
		p.plugins[mode] = plugin
	}
}

func (p *PluginManager) Stop() error {
	klog.Infof("Plugin manager stop")
	errs := []error{}
	for name, plugin := range p.plugins {
		if err := plugin.Stop(); err != nil {
			klog.Errorf("Failed to stop plugin %s: %v", name, err)
			errs = append(errs, err)
			continue
		}

		if server, ok := p.servers[name]; ok {
			server.GracefulStop()
		}

		klog.Infof("Plugin %s stopped", name)
	}
	return errors.Join(errs...)
}
