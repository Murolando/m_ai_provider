// Package config содержит конфигурацию для маппинга моделей между провайдерами.
package config

import (
	_ "embed"

	"github.com/Murolando/m_ai_provider/entities"
	"gopkg.in/yaml.v3"
)

//go:embed models.yaml
var modelsConfigData []byte

// Config представляет структуру конфигурационного файла.
type Config struct {
	ProviderMappings map[string]map[string]string `yaml:"provider_mappings"` // Маппинги моделей для каждого провайдера
}

// GlobalConfig содержит глобальную конфигурацию приложения.
var GlobalConfig *Config

// HydraNamesMap содержит маппинг внутренних названий моделей на названия в HydraAI.
var HydraNamesMap map[entities.ModelName]string

// OpenRouterNamesMap содержит маппинг внутренних названий моделей на названия в OpenRouter.
var OpenRouterNamesMap map[entities.ModelName]string

func init() {
	var config Config
	if err := yaml.Unmarshal(modelsConfigData, &config); err != nil {
		return
	}
	GlobalConfig = &config

	// Инициализируем маппинги для провайдеров
	HydraNamesMap = make(map[entities.ModelName]string)
	OpenRouterNamesMap = make(map[entities.ModelName]string)

	// Заполняем маппинг для Hydra
	if hydraMappings, exists := config.ProviderMappings["hydra"]; exists {
		for internalName, externalName := range hydraMappings {
			HydraNamesMap[entities.ModelName(internalName)] = externalName
		}
	}

	// Заполняем маппинг для OpenRouter
	if openrouterMappings, exists := config.ProviderMappings["openrouter"]; exists {
		for internalName, externalName := range openrouterMappings {
			OpenRouterNamesMap[entities.ModelName(internalName)] = externalName
		}
	}
}
