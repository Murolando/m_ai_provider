package mcp

// MCP (Model Context Protocol) tools structures
// Based on Anthropic's MCP specification: https://modelcontextprotocol.io/

// Tool представляет инструмент в формате MCP.
type Tool struct {
	Name        string      `json:"name"`                  // Имя инструмента (обязательный)
	Description *string     `json:"description,omitempty"` // Описание инструмента
	InputSchema Schema      `json:"inputSchema"`           // JSON Schema для входных параметров
}

// Schema представляет JSON Schema для параметров инструмента.
type Schema struct {
	Type                 string                    `json:"type"`                           // Тип схемы (обычно "object")
	Properties           map[string]SchemaProperty `json:"properties,omitempty"`           // Свойства объекта
	Required             []string                  `json:"required,omitempty"`             // Обязательные поля
	AdditionalProperties *bool                     `json:"additionalProperties,omitempty"` // Разрешены ли дополнительные свойства
}

// SchemaProperty представляет свойство в JSON Schema.
type SchemaProperty struct {
	Type        string                    `json:"type,omitempty"`        // Тип свойства (string, number, boolean, array, object)
	Description *string                   `json:"description,omitempty"` // Описание свойства
	Enum        []interface{}             `json:"enum,omitempty"`        // Возможные значения для enum
	Items       *SchemaProperty           `json:"items,omitempty"`       // Тип элементов для массива
	Properties  map[string]SchemaProperty `json:"properties,omitempty"`  // Свойства для вложенного объекта
	Required    []string                  `json:"required,omitempty"`    // Обязательные поля для вложенного объекта
	Default     interface{}               `json:"default,omitempty"`     // Значение по умолчанию
	Minimum     *float64                  `json:"minimum,omitempty"`     // Минимальное значение для числа
	Maximum     *float64                  `json:"maximum,omitempty"`     // Максимальное значение для числа
	MinLength   *int                      `json:"minLength,omitempty"`   // Минимальная длина для строки
	MaxLength   *int                      `json:"maxLength,omitempty"`   // Максимальная длина для строки
}

// ToolCall представляет вызов инструмента в формате MCP.
type ToolCall struct {
	ID        string                 `json:"id,omitempty"`  // Уникальный идентификатор вызова (для связи с результатом)
	Name      string                 `json:"name"`          // Имя вызываемого инструмента
	Arguments map[string]interface{} `json:"arguments"`     // Аргументы вызова в виде объекта
}

// ToolResult представляет результат выполнения инструмента в формате MCP.
type ToolResult struct {
	Content []Content `json:"content"`           // Содержимое результата
	IsError *bool     `json:"isError,omitempty"` // Является ли результат ошибкой
}

// Content представляет содержимое в MCP формате.
type Content struct {
	Type string      `json:"type"`           // Тип содержимого ("text", "image", "resource")
	Text *string     `json:"text,omitempty"` // Текстовое содержимое
	Data *string     `json:"data,omitempty"` // Данные в base64 (для изображений)
	URI  *string     `json:"uri,omitempty"`  // URI ресурса
}

// Константы для типов содержимого
const (
	ContentTypeText     = "text"
	ContentTypeImage    = "image"
	ContentTypeResource = "resource"
)

// Константы для типов схем
const (
	SchemaTypeObject  = "object"
	SchemaTypeString  = "string"
	SchemaTypeNumber  = "number"
	SchemaTypeInteger = "integer"
	SchemaTypeBoolean = "boolean"
	SchemaTypeArray   = "array"
)

// NewTool создает новый MCP инструмент.
func NewTool(name string, description *string, inputSchema Schema) Tool {
	return Tool{
		Name:        name,
		Description: description,
		InputSchema: inputSchema,
	}
}

// NewSchema создает новую JSON Schema для MCP инструмента.
func NewSchema(schemaType string) Schema {
	return Schema{
		Type:       schemaType,
		Properties: make(map[string]SchemaProperty),
	}
}

// AddProperty добавляет свойство к схеме.
func (s *Schema) AddProperty(name string, property SchemaProperty) {
	if s.Properties == nil {
		s.Properties = make(map[string]SchemaProperty)
	}
	s.Properties[name] = property
}

// AddRequired добавляет обязательное поле к схеме.
func (s *Schema) AddRequired(fieldName string) {
	s.Required = append(s.Required, fieldName)
}

// NewSchemaProperty создает новое свойство схемы.
func NewSchemaProperty(propType string, description *string) SchemaProperty {
	return SchemaProperty{
		Type:        propType,
		Description: description,
	}
}

// NewToolCall создает новый вызов MCP инструмента.
func NewToolCall(name string, arguments map[string]interface{}) ToolCall {
	return ToolCall{
		Name:      name,
		Arguments: arguments,
	}
}

// NewToolResult создает новый результат выполнения MCP инструмента.
func NewToolResult(content []Content, isError *bool) ToolResult {
	return ToolResult{
		Content: content,
		IsError: isError,
	}
}

// NewTextContent создает новое текстовое содержимое.
func NewTextContent(text string) Content {
	return Content{
		Type: ContentTypeText,
		Text: &text,
	}
}

// NewImageContent создает новое содержимое с изображением.
func NewImageContent(data string) Content {
	return Content{
		Type: ContentTypeImage,
		Data: &data,
	}
}

// NewResourceContent создает новое содержимое с ресурсом.
func NewResourceContent(uri string) Content {
	return Content{
		Type: ContentTypeResource,
		URI:  &uri,
	}
}
