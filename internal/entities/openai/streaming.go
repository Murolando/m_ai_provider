package openai

// ChatCompletionStreamResponse представляет streaming ответ от OpenAI Chat Completions API.
// Используется при stream=true в запросе.
type ChatCompletionStreamResponse struct {
	ID                string                        `json:"id"`                           // Уникальный идентификатор чат-сессии
	Object            string                        `json:"object"`                       // Тип объекта, обычно "chat.completion.chunk"
	Created           int64                         `json:"created"`                      // Unix-время создания ответа
	Model             string                        `json:"model"`                        // Модель, которая обработала запрос
	SystemFingerprint *string                       `json:"system_fingerprint,omitempty"` // Отпечаток системы для отслеживания изменений
	Choices           []ChatCompletionStreamChoice  `json:"choices"`                      // Массив с вариантами ответов
	Usage             *Usage                        `json:"usage,omitempty"`              // Информация об использовании токенов (только в последнем chunk)
}

// ChatCompletionStreamChoice представляет один вариант ответа в streaming режиме.
type ChatCompletionStreamChoice struct {
	Index        int                        `json:"index"`                  // Индекс варианта ответа в массиве choices
	Delta        ChatCompletionStreamDelta  `json:"delta"`                  // Дельта изменений для этого chunk
	Logprobs     *ChoiceLogprobs           `json:"logprobs,omitempty"`     // Логарифмы вероятностей токенов
	FinishReason *string                   `json:"finish_reason,omitempty"` // Причина завершения (только в последнем chunk)
}

// ChatCompletionStreamDelta представляет изменения в streaming ответе.
// Содержит только те поля, которые изменились в данном chunk.
type ChatCompletionStreamDelta struct {
	Role      *string    `json:"role,omitempty"`       // Роль (только в первом chunk)
	Content   *string    `json:"content,omitempty"`    // Часть контента
	ToolCalls []ToolCallDelta `json:"tool_calls,omitempty"` // Изменения в вызовах инструментов
}

// ToolCallDelta представляет изменения в вызове инструмента в streaming режиме.
type ToolCallDelta struct {
	Index    *int                 `json:"index,omitempty"`    // Индекс вызова инструмента
	ID       *string              `json:"id,omitempty"`       // ID вызова (только в первом chunk для этого вызова)
	Type     *string              `json:"type,omitempty"`     // Тип вызова (только в первом chunk для этого вызова)
	Function *FunctionCallDelta   `json:"function,omitempty"` // Изменения в вызове функции
}

// FunctionCallDelta представляет изменения в вызове функции в streaming режиме.
type FunctionCallDelta struct {
	Name      *string `json:"name,omitempty"`      // Имя функции (только в первом chunk)
	Arguments *string `json:"arguments,omitempty"` // Часть аргументов функции
}

// ChoiceLogprobs содержит информацию о логарифмах вероятностей для streaming ответа.
type ChoiceLogprobs struct {
	Content []TokenLogprob `json:"content,omitempty"` // Логарифмы вероятностей для токенов контента
}

// Константы для streaming объектов
const (
	ObjectChatCompletionChunk = "chat.completion.chunk" // Тип объекта для streaming chunk
)

// NewChatCompletionStreamResponse создает новый streaming ответ.
func NewChatCompletionStreamResponse(id, model string, created int64) *ChatCompletionStreamResponse {
	return &ChatCompletionStreamResponse{
		ID:      id,
		Object:  ObjectChatCompletionChunk,
		Created: created,
		Model:   model,
		Choices: make([]ChatCompletionStreamChoice, 0),
	}
}

// NewStreamChoice создает новый выбор для streaming ответа.
func NewStreamChoice(index int, delta ChatCompletionStreamDelta) ChatCompletionStreamChoice {
	return ChatCompletionStreamChoice{
		Index: index,
		Delta: delta,
	}
}

// NewContentDelta создает дельту с текстовым контентом.
func NewContentDelta(content string) ChatCompletionStreamDelta {
	return ChatCompletionStreamDelta{
		Content: &content,
	}
}

// NewRoleDelta создает дельту с ролью (обычно для первого chunk).
func NewRoleDelta(role string) ChatCompletionStreamDelta {
	return ChatCompletionStreamDelta{
		Role: &role,
	}
}

// NewToolCallDelta создает дельту для вызова инструмента.
func NewToolCallDelta(index int, id, toolType string) ToolCallDelta {
	return ToolCallDelta{
		Index: &index,
		ID:    &id,
		Type:  &toolType,
	}
}

// NewFunctionCallDelta создает дельту для вызова функции.
func NewFunctionCallDelta(name, arguments string) *FunctionCallDelta {
	delta := &FunctionCallDelta{}
	if name != "" {
		delta.Name = &name
	}
	if arguments != "" {
		delta.Arguments = &arguments
	}
	return delta
}