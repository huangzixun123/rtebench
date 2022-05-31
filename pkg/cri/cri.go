package cri

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/huangzixun123/rtebench/pkg/util"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const maxRemovalAttempts = 10
const maxRemovalTimeout = 10

const defaultNamespace = "rtebench"

type Client struct {
	Runtime runtimeapi.RuntimeServiceClient
	Image   runtimeapi.ImageServiceClient
	conn    *grpc.ClientConn
}

var defaultLinuxPodLabels = map[string]string{}

func (api *Client) CreateContainer(sandbox *runtimeapi.PodSandboxConfig, pod, name, image string, command []string) (string, error) {
	return api.CreateContainerWithResources(sandbox, pod, name, image, command, nil)
}

func (api *Client) CreateContainerWithResources(sandbox *runtimeapi.PodSandboxConfig, pod, name, image string, command []string, resources *runtimeapi.LinuxContainerResources) (string, error) {
	container := &runtimeapi.ContainerConfig{
		Metadata: &runtimeapi.ContainerMetadata{
			Name:    name,
			Attempt: 0,
		},
		Image: &runtimeapi.ImageSpec{
			Image: image,
		},
		LogPath: "/var/log/pods/" + name + ".log",
		Mounts: []*runtimeapi.Mount{
			&runtimeapi.Mount{
				ContainerPath: "/var/log/pods/" + name + ".log",
				HostPath:      "/tmp/log",
			},
		},
		Linux: &runtimeapi.LinuxContainerConfig{
			Resources: resources,
		},
	}
	if command != nil {
		container.Command = command
	}
	req := &runtimeapi.CreateContainerRequest{
		PodSandboxId:  pod,
		Config:        container,
		SandboxConfig: sandbox,
	}
	resp, err := api.Runtime.CreateContainer(context.Background(), req)
	if err != nil {
		return "", err
	}
	return resp.ContainerId, nil
}

func (api *Client) StartContainer(container string) error {
	_, err := api.Runtime.StartContainer(context.Background(), &runtimeapi.StartContainerRequest{
		ContainerId: container,
	})
	if err != nil {
		return err
	}
	return nil
}

func (api *Client) StopContainer(container string, timeout int) error {
	_, err := api.Runtime.StopContainer(context.Background(), &runtimeapi.StopContainerRequest{
		ContainerId: container,
		Timeout:     int64(timeout),
	})
	if err != nil {
		return err
	}
	return nil
}

func (api *Client) RemoveContainer(container string) error {
	_, err := api.Runtime.RemoveContainer(context.Background(), &runtimeapi.RemoveContainerRequest{
		ContainerId: container,
	})
	if err != nil {
		return err
	}
	return nil
}

func (api *Client) StartSandbox(sandbox *runtimeapi.PodSandboxConfig, runtime string) (string, error) {
	resp, err := api.Runtime.RunPodSandbox(context.Background(), &runtimeapi.RunPodSandboxRequest{
		Config:         sandbox,
		RuntimeHandler: runtime,
	})
	if err != nil {
		return "", err
	}
	return resp.PodSandboxId, nil
}

func (api *Client) StopAndRemoveContainer(container string) (err error) {
	for attempt := 0; attempt < maxRemovalAttempts; attempt++ {
		err = api.StopContainer(container, maxRemovalTimeout)
		if err != nil {
			continue
		}
		err = api.RemoveContainer(container)
		if err != nil {
			continue
		}
		return nil
	}
	return errors.Errorf("stop-remove container failed: %v", err)
}

func (api *Client) StopSandbox(pod string) error {
	_, err := api.Runtime.StopPodSandbox(context.Background(), &runtimeapi.StopPodSandboxRequest{
		PodSandboxId: pod,
	})
	if err != nil {
		return err
	}
	return nil
}

func (api *Client) RemoveSandbox(pod string) error {
	_, err := api.Runtime.RemovePodSandbox(context.Background(), &runtimeapi.RemovePodSandboxRequest{
		PodSandboxId: pod,
	})
	if err != nil {
		return err
	}
	return nil
}

func (api *Client) StopAndRemoveSandbox(pod string) (err error) {
	for attempt := 0; attempt < maxRemovalAttempts; attempt++ {
		err = api.StopSandbox(pod)
		if err != nil {
			continue
		}
		err = api.RemoveSandbox(pod)
		if err != nil {
			continue
		}
		return nil
	}
	return errors.Errorf("stop-remove pod failed: %v", err)
}

func (api *Client) InitLinuxSandbox(name string) *runtimeapi.PodSandboxConfig {
	return &runtimeapi.PodSandboxConfig{
		Metadata: &runtimeapi.PodSandboxMetadata{
			Name:      name,
			Uid:       util.NewUUID(),
			Namespace: defaultNamespace,
			Attempt:   1,
		},
		Linux:  &runtimeapi.LinuxPodSandboxConfig{},
		Labels: defaultLinuxPodLabels,
	}
}

func (api *Client) PullImage(image string, sandbox *runtimeapi.PodSandboxConfig) error {
	if !strings.Contains(image, ":") {
		image = image + ":latest"
	}
	imageSpec := &runtimeapi.ImageSpec{
		Image: image,
	}
	_, err := api.Image.PullImage(context.Background(), &runtimeapi.PullImageRequest{
		Image: imageSpec,
	})
	if err != nil {
		return err
	}
	return nil
}

func (api *Client) Close() {
	api.conn.Close()
}

func (api *Client) WaitForLogs(container string) ([]byte, error) {
	for {
		resp, err := api.Runtime.ContainerStatus(context.Background(), &runtimeapi.ContainerStatusRequest{
			ContainerId: container,
		})
		status := resp.Status
		if err != nil {
			return nil, err
		}
		if status.State >= 2 {
			break
		}

		<-time.After(time.Second)
	}
	buf := &bytes.Buffer{}

	if err := api.Logs(container, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (api *Client) Logs(container string, writer io.Writer) error {
	resp, err := api.Runtime.ContainerStatus(context.Background(), &runtimeapi.ContainerStatusRequest{
		ContainerId: container,
	})
	status := resp.Status
	if err != nil {
		return err
	}
	logPath := status.GetLogPath()
	if logPath == "" {
		return errors.New("missing log path")
	}

	f, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %v", logPath, err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		items := bytes.SplitN(line, []byte(" "), 4)
		if len(items) == 4 {
			fmt.Fprintln(writer, string(items[3]))
		}
	}
	return nil
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial")
	}
	runtimeSvc := runtimeapi.NewRuntimeServiceClient(conn)
	imageSvc := runtimeapi.NewImageServiceClient(conn)
	runtimeClient := &Client{
		Runtime: runtimeSvc,
		Image:   imageSvc,
		conn:    conn,
	}
	return runtimeClient, nil
}
