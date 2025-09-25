package types

// Main configuration
type ServiceConfig struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Replicas    int               `yaml:"replicas"`
	Ports       []PortMapping     `yaml:"ports"`
	Environment map[string]string `yaml:"environment"`
	Resources   ResourceLimits    `yaml:"resources"`
	ScalePolicy ScalePolicy       `yaml:"scale_policy"`
	HealthCheck HealthCheck       `yaml:"health_check"`
}

// Type for ports
type PortMapping struct {
	HostPort      int    `yaml:"host_port"`
	ContainerPort int    `yaml:"container_port"`
	Protocol      string `yaml:"protocol"`
}

// Type for server resource limits
type ResourceLimits struct {
	CPU    string `yaml:"cpu"`    // 0.5..1.0
	Memory string `yaml:"memory"` // "512m, 1g"
}

// Type for scale policy
type ScalePolicy struct {
	MinReplicas     int     `yaml:"min_replicas"`
	MaxReplicas     int     `yaml:"max_replicas"`
	TargetCPU       float64 `yaml:"target_cpu"`    // 70.0 = 70%
	TargetMemory    float64 `yaml:"target_memory"` // same
	CooldownSeconds int     `yaml:"cooldown_seconds"`
}

// Type for node health checking
type HealthCheck struct {
	Command  []string `yaml:"command"`
	Interval int      `yaml:"interval"`
	Timeout  int      `yaml:"timeout"`
	Retries  int      `yaml:"retries"`
}

// Node type
type Node struct {
	ID       string            `json:"id"`
	Address  string            `json:"address"`
	Capacity NodeCapacity      `json:"capacity"`
	Usage    NodeUsage         `json:"usage"`
	Labels   map[string]string `json:"labels"`
	Status   string            `json:"status"` // "ready", "draining", "down"
}

// Node capacity struct
type NodeCapacity struct {
	CPU    float64 `json:"cpu"`    // in CPU cores
	Memory int64   `json:"memory"` // in bytes
}

// Node usage struct
type NodeUsage struct {
	CPU    float64 `json:"cpu_usage"`
	Memory int64   `json:"memory_usage"`
}
