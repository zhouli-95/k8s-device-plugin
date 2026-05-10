package pluginmanager

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func StartServer(name string, p pluginapi.DevicePluginServer) (*grpc.Server, error) {
	socketPath := filepath.Join(pluginapi.DevicePluginPath, getEndpoint(name))
	os.Remove(socketPath)
	sock, err := net.Listen("unix", socketPath)
	if err != nil {
		klog.Errorf("Failed to listen, %v", err)
		return nil, err
	}

	server := grpc.NewServer([]grpc.ServerOption{}...)
	pluginapi.RegisterDevicePluginServer(server, p)

	go func() {
		if err := server.Serve(sock); err != nil {
			klog.Errorf("Failed to serve, %v", err)
		}
	}()

	// Wait for server to start.
	conn, err := grpc.NewClient(pluginapi.KubeletSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", addr)
		}),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return server, nil
}

func RegisterPlugin(domain, name string) error {
	conn, err := grpc.DialContext(context.TODO(), pluginapi.KubeletSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", addr)
		}),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	req := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     getEndpoint(name),
		ResourceName: getResourceName(domain, name),
		Options:      &pluginapi.DevicePluginOptions{},
	}

	_, err = client.Register(context.TODO(), req)
	if err != nil {
		return err
	}
	return nil
}

func getEndpoint(name string) string {
	return fmt.Sprintf("%s-device-plugin.sock", name)
}

func getResourceName(domain, name string) string {
	return fmt.Sprintf("%s/%s", domain, name)
}
