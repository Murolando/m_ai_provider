// Package provider содержит интерфейсы и реализации для работы с различными AI провайдерами.
package provider

import (
	"context"

	"github.com/Murolando/m_ai_provider/entities"
	int_entities "github.com/Murolando/m_ai_provider/internal/entities"
	"github.com/Murolando/m_ai_provider/options"
	"github.com/shopspring/decimal"
)

// Provider представляет интерфейс для работы с AI провайдерами.
// Провайдер - проводник до модели, будь то владелец модели или другой ai-hub.
//
// Поддерживаемые провайдеры:
//   - openrouter - https://openrouter.ai/ - txt
//   - hydraai - https://hydraai.app/ - txt, mcp
type Provider interface {
	// SendMessage отправляет сообщения в AI модель через провайдера.
	// ctx - контекст для управления временем жизни запроса
	// messages - массив сообщений для отправки в модель (история чата)
	// modelName - название модели для использования
	// options - дополнительные опции (MCP tools, температура и т.д.)
	// Возвращает ответ от модели с текстом, количеством токенов, стоимостью и возможными tool calls
	SendMessage(ctx context.Context, messages []*entities.Message, modelName entities.ModelName, options ...options.SendMessageOption) (*entities.ProviderMessageResponseDTO, error)

	// GetModelInfo получает информацию о конкретной модели.
	// modelName - название модели для получения информации
	// Возвращает структуру с названием, алиасом и ценой модели в рублях
	GetModelInfo(modelName entities.ModelName) (*entities.ModelInfo, error)

	// getModels загружает и кэширует список доступных моделей от провайдера.
	// Приватный метод для внутреннего использования провайдером.
	// Возвращает ошибку если не удалось получить список моделей
	getModels() error

	// calculatePrice рассчитывает цену на основе переданных параметров.
	// params - параметры для расчета цены (поддерживает разные типы через интерфейс PricingParams)
	// Возвращает цену в рублях и ошибку если расчет невозможен
	calculatePrice(params int_entities.PricingParams) (decimal.Decimal, error)
}
