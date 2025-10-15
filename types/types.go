package types

import "time"

// Main configuration
type ServiceConfig struct {
	Name        string            `yaml:"name"`     // Имя сервиса
	Image       string            `yaml:"image"`    // Image для сервиса
	Replicas    int               `yaml:"replicas"` // Количество реплик данного сервиса
	Ports       []PortMapping     `yaml:"ports"`    // Порты для сервиса -> PortMapping
	Environment map[string]string `yaml:"environment"`
	Resources   ResourceLimits    `yaml:"resources"`    // Выделенные ресурсы компьютера под сервис
	ScalePolicy ScalePolicy       `yaml:"scale_policy"` // Политика расширения (верт./горизонт.)
	HealthCheck HealthCheck       `yaml:"health_check"`
}

// ServiceStatus -> состояние сервиса
type ServiceStatus struct {
	DesiredReplicas int       `json:"desired_replicas"` // Желанное количество реплик
	RunningReplicas int       `json:"running_replicas"` // Количество запущенных реплик
	Health          string    `json:"health"`           // 'healthy', 'unhealthy', 'degraded'
	LastUpdated     time.Time `json:"last_updated"`     // Последнее обновление статуса
}

// PortMapping -> Порт хоста, Порт контейнера, Протокол для обращения
// Protocol - типизированная константа. Значения -> 'udp', 'tcp'
type Protocol string

const (
	TCP Protocol = "tcp" // для tcp протокола
	UDP Protocol = "udp" // для udp
)

type PortMapping struct {
	HostPort      int      `yaml:"host_port"`      // порт хоста
	ContainerPort int      `yaml:"container_port"` // порт контейнера
	Protocol      Protocol `yaml:"protocol"`       // типизированная константа (tcp/udp)
}

// ResourceLimits -> Ресурсы, выделяемые под конкретный сервис (контейнер)
type ResourceLimits struct {
	CPUMillicores int64 `yaml:"cpu_millicores"` // 1000 -> 1 CPU ядро
	MemoryBytes   int64 `yaml:"memory_bytes"`   // 536870912 = 512 MB
}

// ScalePolicy -> автоскейлинг
type ScalePolicy struct {
	MinReplicas     int     `yaml:"min_replicas"`     // Минимальное количесвто реплик
	MaxReplicas     int     `yaml:"max_replicas"`     // Максимальное количество реплик
	TargetCPU       float64 `yaml:"target_cpu"`       // 70.0 = 70% (значение для автоскейлинга)
	TargetMemory    float64 `yaml:"target_memory"`    // Также как и для CPU
	CooldownSeconds int     `yaml:"cooldown_seconds"` //
}

// Проверка состояния контейнера (сервиса)
type HealthCheck struct {
	Type     string   `yaml:"type"`                // "command", "http", "tcp"
	Command  []string `yaml:"command,omitempty"`   // Команда для проверки состояния
	HTTPPath string   `yaml:"http_path,omitempty"` // HTTP путь (адрес)
	Port     int      `yaml:"port,omitempty"`      // Порт для обращения
	Interval int      `yaml:"interval"`            // Интервал для проверки
	Timeout  int      `yaml:"timeout"`             // Таймаут
	Retries  int      `yaml:"retries"`             // Количество повторов отправки запроса провреки состояния
}

// Node -> тип ноды
type Node struct {
	ID       string            `json:"id"`       // ID ноды
	Address  string            `json:"address"`  // Адресс ноды
	Capacity NodeCapacity      `json:"capacity"` // Вместимость (ядра CPU и память)
	Usage    NodeUsage         `json:"usage"`    // Использование (текущее использование)
	Labels   map[string]string `json:"labels"`   // Заголовки
	Status   string            `json:"status"`   // "ready", "draining", "down"
}

// NodeCapacity
type NodeCapacity struct {
	CPU    float64 `json:"cpu"`    // in CPU cores
	Memory int64   `json:"memory"` // in bytes
}

// NodeUsage
type NodeUsage struct {
	CPU    float64 `json:"cpu_usage"`
	Memory int64   `json:"memory_usage"`
}

// Task -> структура для задачи
type Task struct {
	ID          string            `json:"id"`          // ID задачи
	ServiceID   string            `json:"service_id"`  // ID сервиса
	NodeID      string            `json:"node_id"`     // ID ноды
	Image       string            `json:"image"`       // Имя image
	Status      string            `json:"task_status"` // Статус задачи: 'running, stopped, failed'
	Ports       []PortMapping     `json:"ports"`
	Environment map[string]string `json:"environment"`
}

// Event -> событие
type Event struct {
	ID        string                 `json:"id"`             // ID события
	Type      string                 `json:"type"`           // Тип события: 'service_scaled, node_added, task_failed'
	Message   string                 `json:"message"`        // Сообщение
	Timestamp time.Time              `json:"timestamp"`      // Время собятия
	Data      map[string]interface{} `json:"data,omitempty"` // Дата события
}
