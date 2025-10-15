package main

import (
	"fmt"
	"log"

	"github.com/exitae337/orchestergo/config"
	"github.com/exitae337/orchestergo/types"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Валидация переданной конфигурации
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("Failed to validate config: %v", err)
	}
	// Вывод информации для тестирования
	fmt.Printf("Starting orchestrator:\n")
	fmt.Printf("  Cluster: %s\n", cfg.ClusterName)
	fmt.Printf("  API Address: %s\n", cfg.ListenAddr)
	fmt.Printf("  Data Directory: %s\n", cfg.DataDir)
	fmt.Printf("  Services: %d\n", len(cfg.Services))
	// Старт оркестратора
	startOrchestrator(cfg)
}

func startOrchestrator(cfg *types.OrchestratorConfig) {
	// TODO: Инициализировать компоненты оркестратора
	fmt.Printf("Orchestrtor is starting! API working on: %v", cfg.ListenAddr)
	// 1. Docker client
	// 2. API сервер
	// 3. Scheduler
	// 4. Health checker
}
