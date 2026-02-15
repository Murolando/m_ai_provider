package entities

// ChatCompletionRequest представляет запрос к HydraAI API для генерации чата.
type ChatCompletionRequest struct {
	Model            string        `json:"model"`                       // ID модели для генерации (обязательный)
	Messages         []ChatMessage `json:"messages"`                    // Массив объектов сообщений (обязательный)
	MaxTokens        *int          `json:"max_tokens,omitempty"`        // Максимальное количество токенов в ответе (0 = без ограничений)
	Temperature      *float64      `json:"temperature,omitempty"`       // "Креативность" ответа (от 0.0 до 2.0, по умолчанию 0.7)
	Stream           *bool         `json:"stream,omitempty"`            // true для получения ответа в виде потока
	TopP             *float64      `json:"top_p,omitempty"`             // Ядерная выборка
	TopK             *int          `json:"top_k,omitempty"`             // Модель выбирает из k наиболее вероятных токенов
	FrequencyPenalty *float64      `json:"frequency_penalty,omitempty"` // Штраф за частоту токенов (от -2.0 до 2.0)
	PresencePenalty  *float64      `json:"presence_penalty,omitempty"`  // Штраф за наличие токенов (от -2.0 до 2.0)
	WebSearch        *bool         `json:"web_search,omitempty"`        // true, если нужно выполнить поиск в интернете
}

// ChatMessage представляет одно сообщение в диалоге.
type ChatMessage struct {
	Role    string      `json:"role"`    // Роль автора: user, assistant, system
	Content interface{} `json:"content"` // Содержимое сообщения (строка или массив для мультимодальности)
}

// TextContent представляет текстовое содержимое сообщения.
type TextContent struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"` // Текстовое содержимое
}

// ImageURLContent представляет содержимое с изображением.
type ImageURLContent struct {
	Type     string `json:"type"`      // "image_url"
	ImageURL string `json:"image_url"` // URL изображения (https://... или data:image/jpeg;base64,...)
}

// ChatCompletionResponse представляет ответ от HydraAI API.
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`      // Уникальный идентификатор чат-сессии
	Object  string                 `json:"object"`  // Тип объекта, обычно "chat.completion"
	Created int64                  `json:"created"` // Unix-время создания ответа
	Model   string                 `json:"model"`   // Модель, которая обработала запрос
	Choices []ChatCompletionChoice `json:"choices"` // Массив с вариантами ответов
	Usage   ChatCompletionUsage    `json:"usage"`   // Информация о токенах и стоимости
}

// ChatCompletionChoice представляет один вариант ответа.
type ChatCompletionChoice struct {
	Message      ChatMessage `json:"message"`       // Объект сообщения от ассистента
	FinishReason string      `json:"finish_reason"` // Причина завершения: stop, length, tool_calls
	Index        int         `json:"index"`         // Индекс варианта ответа в массиве choices
}

// ChatCompletionUsage содержит информацию об использовании токенов и стоимости.
type ChatCompletionUsage struct {
	PromptTokens            int                      `json:"prompt_tokens"`                       // Количество токенов в запросе
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`     // Детализация токенов в запросе
	CompletionTokens        int                      `json:"completion_tokens"`                   // Количество токенов в ответе
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"` // Детализация токенов в ответе
	TotalTokens             int                      `json:"total_tokens"`                        // Суммарное количество токенов
	TotalTime               float64                  `json:"total_time"`                          // Общее время генерации в секундах
	CostRequest             float64                  `json:"cost_request"`                        // Итоговая стоимость запроса в рублях
	FreeRequest             *bool                    `json:"free_request,omitempty"`              // Был запрос бесплатным или нет
}

// PromptTokensDetails содержит детализацию токенов в запросе.
type PromptTokensDetails struct {
	// Здесь могут быть дополнительные поля для детализации токенов запроса
	// Пока оставляем пустой, так как в документации не указана точная структура
}

// CompletionTokensDetails содержит детализацию токенов в ответе.
type CompletionTokensDetails struct {
	// Здесь могут быть дополнительные поля для детализации токенов ответа
	// Пока оставляем пустой, так как в документации не указана точная структура
}

// Константы для ролей сообщений
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

// Константы для типов контента
const (
	ContentTypeText     = "text"
	ContentTypeImageURL = "image_url"
)

// Константы для причин завершения
const (
	FinishReasonStop      = "stop"       // Естественное завершение
	FinishReasonLength    = "length"     // Достигнут лимит max_tokens
	FinishReasonToolCalls = "tool_calls" // Модель вызвала инструмент
)

// NewChatCompletionRequest создает новый запрос с базовыми параметрами.
func NewChatCompletionRequest(model string, messages []ChatMessage) *ChatCompletionRequest {
	return &ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}
}

// NewTextMessage создает новое текстовое сообщение.
func NewTextMessage(role, text string) ChatMessage {
	return ChatMessage{
		Role:    role,
		Content: text,
	}
}

// NewMultimodalMessage создает новое мультимодальное сообщение с текстом и изображениями.
// func NewMultimodalMessage(role string, contents []interface{}) ChatMessage {
// 	return ChatMessage{
// 		Role:    role,
// 		Content: contents,
// 	}
// }

// NewTextContent создает новый текстовый контент для мультимодального сообщения.
func NewTextContent(text string) TextContent {
	return TextContent{
		Type: ContentTypeText,
		Text: text,
	}
}

// NewImageURLContent создает новый контент с изображением для мультимодального сообщения.
func NewImageURLContent(url string) ImageURLContent {
	return ImageURLContent{
		Type:     ContentTypeImageURL,
		ImageURL: url,
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
