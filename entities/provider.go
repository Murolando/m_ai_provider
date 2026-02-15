// Package entities содержит основные структуры данных для работы с AI провайдерами.
package entities

import "github.com/shopspring/decimal"

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
	MessageText   string          // Текст ответа от AI модели
	TotalTokens   int64           // Общее количество использованных токенов
	PriceInRubles decimal.Decimal // Стоимость запроса в рублях
}
