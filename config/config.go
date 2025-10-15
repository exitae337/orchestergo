package config

import (
	"fmt"
	"os"

	"github.com/exitae337/orchestergo/types"
	"gopkg.in/yaml.v3"
)

// Загрузка конфигурации из .yaml файла (config.yaml)
func LoadConfig(configPath string) (*types.OrchestratorConfig, error) {
	// Чтение файла конфигурации оркестратора
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	// Анмаршалинг файла конфигурации
	var config types.OrchestratorConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data from config: %w", err)
	}
	// Значения по умолчанию
	if config.ListenAddr == "" {
		config.ListenAddr = ":8080"
	}
	if config.DataDir == "" { // Если путь к директории с данными не указан
		config.DataDir = "./orchestrator-data"
	}
	if config.ClusterName == "" { // Если имя кластера не указано
		config.ClusterName = "default-cluster-name"
	}

	// Создагние директории для данных оркестратора
	if err := os.Mkdir(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to make data dir for orchestrator: %w", err)
	}

	return &config, nil
}

func ValidateConfig(config *types.OrchestratorConfig) error {
	// Проверка api адреса оркестратора
	if config.ListenAddr == "" {
		return fmt.Errorf("listen_addr is required")
	}
	// Проверка каждого сревиса
	for i, service := range config.Services {
		if service.Name == "" {
			return fmt.Errorf("service[%d] name is required", i)
		}
		if service.Image == "" {
			return fmt.Errorf("service[%d] image name is required", i)
		}
	}
	return nil
}
