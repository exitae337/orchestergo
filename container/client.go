package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/exitae337/orchestergo/types"
)

type ContainerOptions struct {
	StopTimeout   int
	RemoveVolumes bool
	ForceRemove   bool
}

type ContainerClient struct {
	dockerClient *client.Client
	options      ContainerOptions
}

func NewContainerClient(opts ContainerOptions) (*ContainerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %v", err)
	}
	// Default values
	if opts.StopTimeout == 0 {
		opts.StopTimeout = 30 // in seconds: int -> seconds
	}
	return &ContainerClient{
		dockerClient: cli,
		options:      opts,
	}, nil
}

// DOCKER CONTAINER METHODS

// Create container
func (c *ContainerClient) CreateContainer(ctx context.Context, service types.ServiceConfig, instanceID string) (string, error) {
	resources := container.Resources{} // our resources into -> Docker format
	// NanoCPUs
	if service.Resources.CPU != "" {
		resources.NanoCPUs = parseCPU(service.Resources.CPU) // 0.5 (50%) -> 0.5 * 100000 * 10000
	}
	// Memory
	if service.Resources.Memory != "" {
		resources.Memory = parseMemory(service.Resources.Memory)
	}
	// Config
	config := &container.Config{
		Image: service.Image,
		Env:   convertEnvVars(service.Environment),
	}
	// Add port exposition
	updateContainerConfigs(config, service.Ports)

	hostConfig := &container.HostConfig{
		Resources:    resources,
		PortBindings: convertPorts(service.Ports),
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}

	containerName := fmt.Sprintf("%s-%s", service.Name, instanceID)

	resp, err := c.dockerClient.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil, // networking config
		nil, // platform
		containerName,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %v", err)
	}

	return resp.ID, nil
}

// Start container
func (c *ContainerClient) StartContainer(ctx context.Context, containerID string, checkpointID string) error {
	opts := container.StartOptions{}
	if checkpointID != "" {
		opts.CheckpointID = checkpointID
	}
	return c.dockerClient.ContainerStart(ctx, containerID, opts)
}

// Stop container
func (c *ContainerClient) StopContainer(ctx context.Context, containerID string, signal string) error {
	timeout := c.options.StopTimeout
	stopOpts := container.StopOptions{
		Timeout: &timeout,
	}
	if signal != "" {
		stopOpts.Signal = signal
	}
	return c.dockerClient.ContainerStop(ctx, containerID, stopOpts)
}

// Remove container
func (c *ContainerClient) RemoveContainer(ctx context.Context, containerID string) error {
	return c.dockerClient.ContainerRemove(ctx, containerID, container.RemoveOptions{
		RemoveVolumes: c.options.RemoveVolumes,
		Force:         c.options.ForceRemove,
	})
}

// GraceFul stop (correct container stopping)
func (c *ContainerClient) GracefulStopContainer(ctx context.Context, containerID string) error {
	if err := c.StopContainer(ctx, containerID, "SIGTERM"); err != nil {
		// If don't working:
		fmt.Printf("Soft stop failed, trying force stopping: %v/n", err)
		if err := c.StopContainer(ctx, containerID, "SIGKILL"); err != nil {
			return fmt.Errorf("force stop also failed due to: %v", err)
		}
	}
	return nil
}

// Get container status
func (c *ContainerClient) GetContainerOptions(ctx context.Context, containerID string) (string, error) {
	info, err := c.dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}
	return info.State.Status, nil
}

// Util functions

// Nano CPU parser for Docker config
func parseCPU(cpuStr string) int64 {
	if cpuStr == "" {
		return 0
	}

	var cpu float64
	_, err := fmt.Sscanf(cpuStr, "%f", &cpu)
	if err != nil {
		return 0
	}

	return int64(cpu * 1000000000)
}

// Memory parser for Docker config
func parseMemory(memoryStr string) int64 {
	if memoryStr == "" {
		return 0
	}

	var value int64
	var unit string

	_, err := fmt.Sscanf(memoryStr, "%d%s", &value, &unit)
	if err != nil {
		return 0
	}

	switch unit {
	case "g", "G", "gb", "GB":
		return value * 1024 * 1024 * 1024
	case "m", "M", "mb", "MB":
		return value * 1024 * 1024
	case "k", "K", "kb", "KB":
		return value * 1024
	default:
		return value
	}
}

func convertEnvVars(env map[string]string) []string {
	var envVars []string
	for key, value := range env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}
	return envVars
}
