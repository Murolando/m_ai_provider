package entities

import (
	"github.com/Murolando/m_ai_provider/internal/entities/openai"
)

// HydraChatCompletionRequest представляет запрос к HydraAI API для генерации чата.
// Использует базовые OpenAI структуры с добавлением Hydra-специфичных полей.
type HydraChatCompletionRequest struct {
	openai.ChatCompletionRequest        // Встраиваем базовую OpenAI структуру
	TopK                         *int   `json:"top_k,omitempty"`       // Модель выбирает из k наиболее вероятных токенов (Hydra-специфично)
	WebSearch                    *bool  `json:"web_search,omitempty"`  // true, если нужно выполнить поиск в интернете (Hydra-специфично)
}

// HydraChatCompletionResponse представляет ответ от HydraAI API.
// Использует базовые OpenAI структуры с переопределением Usage для Hydra-специфичных полей.
type HydraChatCompletionResponse struct {
	ID                string                      `json:"id"`                           // Уникальный идентификатор чат-сессии
	Object            string                      `json:"object"`                       // Тип объекта, обычно "chat.completion"
	Created           int64                       `json:"created"`                      // Unix-время создания ответа
	Model             string                      `json:"model"`                        // Модель, которая обработала запрос
	SystemFingerprint *string                     `json:"system_fingerprint,omitempty"` // Отпечаток системы (совместимость с OpenAI)
	Choices           []openai.ChatCompletionChoice `json:"choices"`                    // Массив с вариантами ответов (используем OpenAI структуру)
	Usage             HydraChatCompletionUsage    `json:"usage"`                        // Информация о токенах и стоимости (Hydra-специфично)
}

// HydraChatCompletionUsage содержит информацию об использовании токенов и стоимости для HydraAI.
// Расширяет базовую OpenAI Usage структуру дополнительными полями.
type HydraChatCompletionUsage struct {
	PromptTokens            int                              `json:"prompt_tokens"`                       // Количество токенов в запросе
	PromptTokensDetails     *HydraPromptTokensDetails       `json:"prompt_tokens_details,omitempty"`     // Детализация токенов в запросе
	CompletionTokens        int                              `json:"completion_tokens"`                   // Количество токенов в ответе
	CompletionTokensDetails *HydraCompletionTokensDetails   `json:"completion_tokens_details,omitempty"` // Детализация токенов в ответе
	TotalTokens             int                              `json:"total_tokens"`                        // Суммарное количество токенов
	TotalTime               float64                          `json:"total_time"`                          // Общее время генерации в секундах (Hydra-специфично)
	CostRequest             float64                          `json:"cost_request"`                        // Итоговая стоимость запроса в рублях (Hydra-специфично)
	FreeRequest             *bool                            `json:"free_request,omitempty"`              // Был запрос бесплатным или нет (Hydra-специфично)
}

// HydraPromptTokensDetails содержит детализацию токенов в запросе для HydraAI.
type HydraPromptTokensDetails struct {
	// Базовые поля совместимые с OpenAI
	CachedTokens *int `json:"cached_tokens,omitempty"` // Количество кэшированных токенов
	// Здесь могут быть дополнительные Hydra-специфичные поля
}

// HydraCompletionTokensDetails содержит детализацию токенов в ответе для HydraAI.
type HydraCompletionTokensDetails struct {
	// Базовые поля совместимые с OpenAI
	ReasoningTokens *int `json:"reasoning_tokens,omitempty"` // Количество токенов рассуждения (для o1-подобных моделей)
	// Здесь могут быть дополнительные Hydra-специфичные поля
}

// NewHydraChatCompletionRequest создает новый Hydra запрос с базовыми параметрами.
func NewHydraChatCompletionRequest(model string, messages []openai.ChatMessage) *HydraChatCompletionRequest {
	return &HydraChatCompletionRequest{
		ChatCompletionRequest: openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	}
}

// NewHydraChatCompletionRequestWithWebSearch создает новый Hydra запрос с поиском в интернете.
func NewHydraChatCompletionRequestWithWebSearch(model string, messages []openai.ChatMessage, webSearch bool) *HydraChatCompletionRequest {
	return &HydraChatCompletionRequest{
		ChatCompletionRequest: openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
		WebSearch: &webSearch,
	}
}

// NewHydraChatCompletionRequestWithTopK создает новый Hydra запрос с параметром TopK.
func NewHydraChatCompletionRequestWithTopK(model string, messages []openai.ChatMessage, topK int) *HydraChatCompletionRequest {
	return &HydraChatCompletionRequest{
		ChatCompletionRequest: openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
		TopK: &topK,
	}
}

// ModelsResponse представляет ответ API /models от HydraAI.
type ModelsResponse struct {
	Data []HydraModel `json:"data"` // Массив доступных моделей
}

// HydraModel представляет модель от HydraAI API.
type HydraModel struct {
	ID                 string       `json:"id"`                   // Уникальный идентификатор модели
	Name               string       `json:"name"`                 // Название модели
	Description        string       `json:"description"`          // Описание модели
	Type               string       `json:"type"`                 // Тип модели
	Context            int          `json:"context"`              // Размер контекстного окна
	Active             bool         `json:"active"`               // Активна ли модель
	WebSearch          *bool        `json:"web_search"`           // Поддерживает ли поиск в интернете
	InputModalities    []string     `json:"input_modalities"`     // Поддерживаемые входные модальности
	OutputModalities   []string     `json:"output_modalities"`    // Поддерживаемые выходные модальности
	SupportedFileTypes []string     `json:"supported_file_types"` // Поддерживаемые типы файлов
	MaxImageCount      *int         `json:"max_image_count"`      // Максимальное количество изображений
	MaxFileCount       *int         `json:"max_file_count"`       // Максимальное количество файлов
	Architecture       *string      `json:"architecture"`         // Архитектура модели
	Quantization       *string      `json:"quantization"`         // Тип квантизации
	Pricing            HydraPricing `json:"pricing"`              // Информация о ценообразовании
	OwnedBy            string       `json:"owned_by"`             // Владелец модели
}

// HydraPricing представляет структуру ценообразования модели.
type HydraPricing struct {
	Type              string   `json:"type"`                        // Тип ценообразования (tokens, request)
	InCostPerMillion  *float64 `json:"in_cost_per_million,omitempty"`  // Стоимость входящих токенов за миллион
	OutCostPerMillion *float64 `json:"out_cost_per_million,omitempty"` // Стоимость исходящих токенов за миллион
	CostPerRequest    *float64 `json:"cost_per_request,omitempty"`     // Стоимость за запрос
	CostPerMillion    *float64 `json:"cost_per_million,omitempty"`     // Общая стоимость за миллион токенов
	FreeRequests      bool     `json:"free_requests"`               // Бесплатные ли запросы
}
