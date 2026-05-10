package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
	"k8s.io/klog/v2"

	"github.com/zhouli-95/k8s-device-plugin/internal/config"
	"github.com/zhouli-95/k8s-device-plugin/internal/option"
	"github.com/zhouli-95/k8s-device-plugin/internal/pluginmanager"
	"github.com/zhouli-95/k8s-device-plugin/pkg/resource"
)

var (
	version = "0.0.1"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	c := cli.NewApp()
	o := option.New()
	c.Name = "K8S Device Plugin"
	c.Version = versionString()
	c.Action = func(ctx *cli.Context) error {
		// Initialize configuration before starting the plugin.
		if err := resource.LoadConfig(o.ConfigPath); err != nil {
			return err
		}
		return start(o)
	}

	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "plugin-mode",
			Value:       config.PluginModeNative,
			Destination: &o.PluginMode,
			EnvVars:     []string{"PLUGIN_MODE"},
		},
		&cli.StringFlag{
			Name:        "config-path",
			Value:       config.DefaultConfigPath,
			Destination: &o.ConfigPath,
			EnvVars:     []string{"CONFIG_PATH"},
		},
		&cli.StringFlag{
			Name:        "device-strategy",
			Value:       config.DeviceStrategyNative,
			Destination: &o.DeviceStrategy,
			EnvVars:     []string{"DEVICE_STRATEGY"},
		},
		&cli.IntFlag{
			Name:        "share-number",
			Value:       option.DefaultShareNumber,
			Destination: &o.ShareOption.ShareNumber,
			EnvVars:     []string{"SHARE_NUMBER"},
		},
	}

	err := c.Run(os.Args)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

func start(opt *option.Options) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		klog.Info("receive stop signal")
		cancel()
	}()

	pm := pluginmanager.NewPluginManager(ctx, opt)
	if err := pm.Start(); err != nil {
		return fmt.Errorf("failed to start plugin manager %w", err)
	}
	defer pm.Stop()

	<-ctx.Done()
	klog.Info("device plugin exited")
	return nil
}

func versionString() string {
	v := []string{version}
	v = append(v, "commit: "+commit)
	v = append(v, "date: "+date)
	return strings.Join(v, "\n")
}
