package openai

// ChatCompletionRequest представляет запрос к OpenAI Chat Completions API.
// Полностью совместим с официальной документацией OpenAI.
// Справочник: https://platform.openai.com/docs/api-reference/chat
type ChatCompletionRequest struct {
	Model            string                 `json:"model"`                       // ID модели для генерации (обязательный)
	Messages         []ChatMessage          `json:"messages"`                    // Массив объектов сообщений (обязательный)
	MaxTokens        *int                   `json:"max_tokens,omitempty"`        // Максимальное количество токенов в ответе
	Temperature      *float64               `json:"temperature,omitempty"`       // "Креативность" ответа (от 0.0 до 2.0, по умолчанию 1.0)
	TopP             *float64               `json:"top_p,omitempty"`             // Ядерная выборка (от 0.0 до 1.0, по умолчанию 1.0)
	N                *int                   `json:"n,omitempty"`                 // Количество вариантов ответа (по умолчанию 1)
	Stream           *bool                  `json:"stream,omitempty"`            // true для получения ответа в виде потока
	Stop             interface{}            `json:"stop,omitempty"`              // Строка или массив строк для остановки генерации
	PresencePenalty  *float64               `json:"presence_penalty,omitempty"`  // Штраф за наличие токенов (от -2.0 до 2.0)
	FrequencyPenalty *float64               `json:"frequency_penalty,omitempty"` // Штраф за частоту токенов (от -2.0 до 2.0)
	LogitBias        map[string]int         `json:"logit_bias,omitempty"`        // Модификация вероятности конкретных токенов
	User             *string                `json:"user,omitempty"`              // Уникальный идентификатор пользователя
	ResponseFormat   *ResponseFormat        `json:"response_format,omitempty"`   // Формат ответа (text или json_object)
	Seed             *int                   `json:"seed,omitempty"`              // Seed для детерминированной генерации
	Logprobs         *bool                  `json:"logprobs,omitempty"`          // Возвращать ли логарифмы вероятностей
	TopLogprobs      *int                   `json:"top_logprobs,omitempty"`      // Количество наиболее вероятных токенов (0-20)
	Tools            []Tool                 `json:"tools,omitempty"`             // Список доступных инструментов
	ToolChoice       interface{}            `json:"tool_choice,omitempty"`       // Управление выбором инструментов (строка или ToolChoiceFunction)
}

// ChatMessage представляет одно сообщение в диалоге.
type ChatMessage struct {
	Role       string      `json:"role"`                 // Роль автора: system, user, assistant, tool
	Content    interface{} `json:"content,omitempty"`    // Содержимое сообщения (строка или массив ContentPart)
	Name       *string     `json:"name,omitempty"`       // Имя автора сообщения (для role=function)
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"` // Вызовы инструментов (только для role=assistant)
	ToolCallID *string     `json:"tool_call_id,omitempty"` // ID вызова инструмента (только для role=tool)
}

// ContentPart представляет часть контента в мультимодальном сообщении.
type ContentPart struct {
	Type     string    `json:"type"`               // Тип контента: "text" или "image_url"
	Text     *string   `json:"text,omitempty"`     // Текстовое содержимое (для type="text")
	ImageURL *ImageURL `json:"image_url,omitempty"` // URL изображения (для type="image_url")
}

// ImageURL представляет изображение в сообщении.
type ImageURL struct {
	URL    string  `json:"url"`              // URL изображения (https://... или data:image/...)
	Detail *string `json:"detail,omitempty"` // Уровень детализации: "low", "high", "auto"
}

// ResponseFormat определяет формат ответа от модели.
type ResponseFormat struct {
	Type string `json:"type"` // "text" или "json_object"
}

// ChatCompletionResponse представляет ответ от OpenAI Chat Completions API.
type ChatCompletionResponse struct {
	ID                string                  `json:"id"`                           // Уникальный идентификатор чат-сессии
	Object            string                  `json:"object"`                       // Тип объекта, обычно "chat.completion"
	Created           int64                   `json:"created"`                      // Unix-время создания ответа
	Model             string                  `json:"model"`                        // Модель, которая обработала запрос
	SystemFingerprint *string                 `json:"system_fingerprint,omitempty"` // Отпечаток системы для отслеживания изменений
	Choices           []ChatCompletionChoice  `json:"choices"`                      // Массив с вариантами ответов
	Usage             *Usage                  `json:"usage,omitempty"`              // Информация об использовании токенов
}

// ChatCompletionChoice представляет один вариант ответа.
type ChatCompletionChoice struct {
	Index        int          `json:"index"`         // Индекс варианта ответа в массиве choices
	Message      ChatMessage  `json:"message"`       // Объект сообщения от ассистента
	Logprobs     *Logprobs    `json:"logprobs,omitempty"` // Логарифмы вероятностей токенов
	FinishReason *string      `json:"finish_reason"` // Причина завершения: stop, length, tool_calls, content_filter
}

// Usage содержит информацию об использовании токенов.
type Usage struct {
	PromptTokens            int                      `json:"prompt_tokens"`                       // Количество токенов в запросе
	CompletionTokens        int                      `json:"completion_tokens"`                   // Количество токенов в ответе
	TotalTokens             int                      `json:"total_tokens"`                        // Суммарное количество токенов
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`     // Детализация токенов в запросе
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"` // Детализация токенов в ответе
}

// PromptTokensDetails содержит детализацию токенов в запросе.
type PromptTokensDetails struct {
	CachedTokens *int `json:"cached_tokens,omitempty"` // Количество кэшированных токенов
}

// CompletionTokensDetails содержит детализацию токенов в ответе.
type CompletionTokensDetails struct {
	ReasoningTokens *int `json:"reasoning_tokens,omitempty"` // Количество токенов рассуждения (для o1 моделей)
}

// Logprobs содержит информацию о логарифмах вероятностей токенов.
type Logprobs struct {
	Content []TokenLogprob `json:"content,omitempty"` // Логарифмы вероятностей для токенов контента
}

// TokenLogprob содержит информацию о логарифме вероятности одного токена.
type TokenLogprob struct {
	Token   string             `json:"token"`    // Токен
	Logprob float64            `json:"logprob"`  // Логарифм вероятности токена
	Bytes   []int              `json:"bytes"`    // Байтовое представление токена
	TopLogprobs []TopLogprob    `json:"top_logprobs"` // Топ наиболее вероятных токенов
}

// TopLogprob представляет один из наиболее вероятных токенов.
type TopLogprob struct {
	Token   string  `json:"token"`   // Токен
	Logprob float64 `json:"logprob"` // Логарифм вероятности токена
	Bytes   []int   `json:"bytes"`   // Байтовое представление токена
}

// Константы для ролей сообщений
const (
	RoleSystem    = "system"    // Системное сообщение
	RoleUser      = "user"      // Сообщение пользователя
	RoleAssistant = "assistant" // Сообщение ассистента
	RoleTool      = "tool"      // Результат выполнения инструмента
)

// Константы для типов контента
const (
	ContentTypeText     = "text"      // Текстовый контент
	ContentTypeImageURL = "image_url" // Изображение по URL
)

// Константы для причин завершения
const (
	FinishReasonStop          = "stop"           // Естественное завершение
	FinishReasonLength        = "length"         // Достигнут лимит max_tokens
	FinishReasonToolCalls     = "tool_calls"     // Модель вызвала инструмент
	FinishReasonContentFilter = "content_filter" // Контент отфильтрован
)

// Константы для форматов ответа
const (
	ResponseFormatText       = "text"        // Обычный текстовый ответ
	ResponseFormatJSONObject = "json_object" // Ответ в формате JSON
)

// Константы для уровня детализации изображений
const (
	ImageDetailLow  = "low"  // Низкая детализация (быстрее и дешевле)
	ImageDetailHigh = "high" // Высокая детализация (медленнее и дороже)
	ImageDetailAuto = "auto" // Автоматический выбор
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

// NewMultimodalMessage создает новое мультимодальное сообщение с различными типами контента.
func NewMultimodalMessage(role string, contents []ContentPart) ChatMessage {
	return ChatMessage{
		Role:    role,
		Content: contents,
	}
}

// NewTextContent создает новый текстовый контент для мультимодального сообщения.
func NewTextContent(text string) ContentPart {
	return ContentPart{
		Type: ContentTypeText,
		Text: &text,
	}
}

// NewImageURLContent создает новый контент с изображением для мультимодального сообщения.
func NewImageURLContent(url string, detail *string) ContentPart {
	return ContentPart{
		Type: ContentTypeImageURL,
		ImageURL: &ImageURL{
			URL:    url,
			Detail: detail,
		},
	}
}