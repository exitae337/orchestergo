package container

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/exitae337/orchestergo/types"
)

func convertPorts(portMappings []types.PortMapping) map[nat.Port][]nat.PortBinding {
	if len(portMappings) == 0 {
		return nil
	}
	portBindings := make(map[nat.Port][]nat.PortBinding)
	for _, pm := range portMappings {
		portKey := nat.Port(fmt.Sprintf("%d/%s", pm.ContainerPort, getProtocol(pm.Protocol))) // example for docker: "50/tcp"
		hostPort := ""
		if pm.HostPort > 0 {
			hostPort = fmt.Sprintf("%d", pm.HostPort)
		}
		portBinding := nat.PortBinding{
			HostIP:   "0.0.0.0", // Listen on all interfaces
			HostPort: hostPort,
		}
		portBindings[portKey] = []nat.PortBinding{portBinding}
	}
	return portBindings
}

// Get protocol function: default value -> tcp
func getProtocol(protocol string) string {
	if protocol == "" {
		return "tcp" // tcp for default value
	}
	return protocol
}

// Port masssive for exposition
func createExposedPorts(portMappings []types.PortMapping) map[nat.Port]struct{} {
	if len(portMappings) == 0 {
		return nil
	}
	exposedPorts := make(map[nat.Port]struct{})
	for _, pm := range portMappings {
		portKey := nat.Port(fmt.Sprintf("%d/%s", pm.ContainerPort, getProtocol(pm.Protocol)))
		exposedPorts[portKey] = struct{}{}
	}
	return exposedPorts
}

func updateContainerConfigs(config *container.Config, portMappings []types.PortMapping) {
	if len(portMappings) > 0 {
		config.ExposedPorts = createExposedPorts(portMappings)
	}
}
