package openai

// Tool представляет инструмент (функцию), доступный для модели.
// Справочник: https://platform.openai.com/docs/guides/function-calling
type Tool struct {
	Type     string   `json:"type"`     // Тип инструмента, всегда "function"
	Function Function `json:"function"` // Описание функции
}

// Function представляет описание функции для вызова.
type Function struct {
	Name        string      `json:"name"`                  // Имя функции (обязательный)
	Description *string     `json:"description,omitempty"` // Описание функции
	Parameters  interface{} `json:"parameters,omitempty"`  // JSON Schema параметров функции
}

// ToolCall представляет вызов инструмента моделью.
type ToolCall struct {
	ID       string       `json:"id"`       // Уникальный идентификатор вызова
	Type     string       `json:"type"`     // Тип вызова, всегда "function"
	Function FunctionCall `json:"function"` // Детали вызова функции
}

// FunctionCall представляет конкретный вызов функции.
type FunctionCall struct {
	Name      string `json:"name"`      // Имя вызываемой функции
	Arguments string `json:"arguments"` // Аргументы функции в формате JSON строки
}

// ToolChoice управляет выбором инструментов моделью.
// Может быть строкой ("none", "auto", "required") или объектом ToolChoiceFunction.
type ToolChoice interface{}

// ToolChoiceFunction принуждает модель вызвать конкретную функцию.
type ToolChoiceFunction struct {
	Type     string                   `json:"type"`     // Всегда "function"
	Function ToolChoiceFunctionDetail `json:"function"` // Детали функции для вызова
}

// ToolChoiceFunctionDetail содержит имя функции для принудительного вызова.
type ToolChoiceFunctionDetail struct {
	Name string `json:"name"` // Имя функции для вызова
}

// Константы для типов инструментов
const (
	ToolTypeFunction = "function" // Тип инструмента - функция
)

// Константы для выбора инструментов
const (
	ToolChoiceNone     = "none"     // Не использовать инструменты
	ToolChoiceAuto     = "auto"     // Автоматический выбор (по умолчанию)
	ToolChoiceRequired = "required" // Обязательно использовать инструмент
)

// NewTool создает новый инструмент с функцией.
func NewTool(name string, description *string, parameters interface{}) Tool {
	return Tool{
		Type: ToolTypeFunction,
		Function: Function{
			Name:        name,
			Description: description,
			Parameters:  parameters,
		},
	}
}

// NewToolCall создает новый вызов инструмента.
func NewToolCall(id, functionName, arguments string) ToolCall {
	return ToolCall{
		ID:   id,
		Type: ToolTypeFunction,
		Function: FunctionCall{
			Name:      functionName,
			Arguments: arguments,
		},
	}
}

// NewToolChoiceFunction создает принудительный выбор конкретной функции.
func NewToolChoiceFunction(functionName string) ToolChoiceFunction {
	return ToolChoiceFunction{
		Type: ToolTypeFunction,
		Function: ToolChoiceFunctionDetail{
			Name: functionName,
		},
	}
}

// NewToolMessage создает сообщение с результатом выполнения инструмента.
func NewToolMessage(toolCallID, content string) ChatMessage {
	return ChatMessage{
		Role:       RoleTool,
		Content:    content,
		ToolCallID: &toolCallID,
	}
}
