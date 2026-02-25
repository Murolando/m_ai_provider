// Package entities содержит основные структуры данных для работы с AI провайдерами.
package entities

import (
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/shopspring/decimal"
)

// ModelName представляет название модели AI.
type ModelName string

const (
	// AuthorTypeUser представляет тип автора сообщения - пользователь.
	AuthorTypeUser = "user"
	// AuthorTypeRobot представляет тип автора сообщения - AI ассистент.
	AuthorTypeRobot = "terminator"

	// MessageText представляет тип сообщения - текст.
	MessageText = "message_text"
	// MessageImage представляет тип сообщения - изображение.
	MessageImage = "message_image"

	// Константы для причин завершения генерации
	// FinishReasonStop естественное завершение генерации
	FinishReasonStop = "stop"
	// FinishReasonLength достигнут лимит максимального количества токенов
	FinishReasonLength = "length"
	// FinishReasonToolCalls модель вызвала инструмент
	FinishReasonToolCalls = "tool_calls"
	// FinishReasonContentFilter контент отфильтрован
	FinishReasonContentFilter = "content_filter"
)

// ModelInfo содержит информацию о модели AI.
type ModelInfo struct {
	Name          string          `json:"name"`            // Человекочитаемое название модели
	Alias         ModelName       `json:"alias"`           // Алиас модели для использования в API
	PriceInRubles decimal.Decimal `json:"price_in_rubles"` // Цена модели в рублях
}

// Message представляет сообщение в чате.
type Message struct {
	ChatID      string `json:"chat_id" db:"chat_id"`           // Идентификатор чата
	MessageText string `json:"message_text" db:"message_text"` // Текст сообщения
	AuthorType  string `json:"author_type" db:"author"`        // Тип автора сообщения
	MessageType string `json:"message_type" db:"message_type"` // Тип сообщения
}

// ProviderMessageResponseDTO содержит ответ от AI провайдера.
type ProviderMessageResponseDTO struct {
	MessageText   string          `json:"message_text"`    // Текст ответа от AI модели
	TotalTokens   int64           `json:"total_tokens"`    // Общее количество использованных токенов
	PriceInRubles decimal.Decimal `json:"price_in_rubles"` // Стоимость запроса в рублях

	// Новые поля для поддержки MCP tool calls
	ToolCalls    []mcpgo.CallToolRequest `json:"tool_calls,omitempty"`    // Вызовы инструментов в MCP формате
	FinishReason *string                 `json:"finish_reason,omitempty"` // Причина завершения (stop, tool_calls, length, etc.)
}
