package mappers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Murolando/m_ai_provider/internal/entities/openai"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

// ToolsMapper предоставляет методы для конвертации между OpenAI и MCP форматами инструментов.
type ToolsMapper struct{}

// NewToolsMapper создает новый экземпляр маппера инструментов.
func NewToolsMapper() *ToolsMapper {
	return &ToolsMapper{}
}

// OpenAIToolToMCP конвертирует OpenAI инструмент в MCP формат.
func (m *ToolsMapper) OpenAIToolToMCP(openaiTool openai.Tool) (mcpgo.Tool, error) {
	// Конвертируем JSON Schema из interface{} в mcpgo.ToolInputSchema
	inputSchema, err := m.convertOpenAIParametersToMCPSchema(openaiTool.Function.Parameters)
	if err != nil {
		return mcpgo.Tool{}, fmt.Errorf("failed to convert parameters to MCP schema: %w", err)
	}

	// Получаем описание
	description := ""
	if openaiTool.Function.Description != nil {
		description = *openaiTool.Function.Description
	}

	return mcpgo.Tool{
		Name:        openaiTool.Function.Name,
		Description: description,
		InputSchema: inputSchema,
	}, nil
}

// MCPToolToOpenAI конвертирует MCP инструмент в OpenAI формат.
func (m *ToolsMapper) MCPToolToOpenAI(mcpTool mcpgo.Tool) (openai.Tool, error) {
	// Конвертируем MCP Schema в interface{} для OpenAI
	parameters, err := m.convertMCPSchemaToOpenAIParameters(mcpTool.InputSchema)
	if err != nil {
		return openai.Tool{}, fmt.Errorf("failed to convert MCP schema to OpenAI parameters: %w", err)
	}

	// Создаем описание как указатель
	var description *string
	if mcpTool.Description != "" {
		description = &mcpTool.Description
	}

	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: openai.Function{
			Name:        mcpTool.Name,
			Description: description,
			Parameters:  parameters,
		},
	}, nil
}

// OpenAIToolCallToMCP конвертирует OpenAI вызов инструмента в MCP формат.
func (m *ToolsMapper) OpenAIToolCallToMCP(openaiToolCall openai.ToolCall) (mcpgo.CallToolRequest, error) {
	// Парсим JSON строку аргументов в map
	var arguments map[string]interface{}
	if openaiToolCall.Function.Arguments != "" {
		if err := json.Unmarshal([]byte(openaiToolCall.Function.Arguments), &arguments); err != nil {
			return mcpgo.CallToolRequest{}, fmt.Errorf("failed to parse OpenAI tool call arguments: %w", err)
		}
	} else {
		arguments = make(map[string]interface{})
	}

	// Создаем CallToolParams
	params := mcpgo.CallToolParams{
		Name:      openaiToolCall.Function.Name,
		Arguments: arguments,
	}

	// Создаем CallToolRequest
	return mcpgo.CallToolRequest{
		Params: params,
	}, nil
}

// MCPToolCallToOpenAI конвертирует MCP вызов инструмента в OpenAI формат.
func (m *ToolsMapper) MCPToolCallToOpenAI(mcpCall mcpgo.CallToolRequest) (openai.ToolCall, error) {
	// Конвертируем Arguments с type assertion
	var arguments map[string]interface{}
	if mcpCall.Params.Arguments != nil {
		if args, ok := mcpCall.Params.Arguments.(map[string]interface{}); ok {
			arguments = args
		} else {
			return openai.ToolCall{}, fmt.Errorf("failed to convert arguments: expected map[string]interface{}, got %T", mcpCall.Params.Arguments)
		}
	} else {
		arguments = make(map[string]interface{})
	}

	// Конвертируем map аргументов в JSON строку
	argumentsJSON, err := json.Marshal(arguments)
	if err != nil {
		return openai.ToolCall{}, fmt.Errorf("failed to marshal MCP tool call arguments: %w", err)
	}

	// Генерируем ID для OpenAI tool call
	id := generateToolCallID()

	return openai.ToolCall{
		ID:   id,
		Type: openai.ToolTypeFunction,
		Function: openai.FunctionCall{
			Name:      mcpCall.Params.Name,
			Arguments: string(argumentsJSON),
		},
	}, nil
}

// OpenAIToolsToMCP конвертирует массив OpenAI инструментов в MCP формат.
func (m *ToolsMapper) OpenAIToolsToMCP(openaiTools []openai.Tool) ([]mcpgo.Tool, error) {
	mcpTools := make([]mcpgo.Tool, len(openaiTools))
	for i, openaiTool := range openaiTools {
		mcpTool, err := m.OpenAIToolToMCP(openaiTool)
		if err != nil {
			return nil, fmt.Errorf("failed to convert OpenAI tool %d to MCP: %w", i, err)
		}
		mcpTools[i] = mcpTool
	}
	return mcpTools, nil
}

// MCPToolsToOpenAI конвертирует массив MCP инструментов в OpenAI формат.
func (m *ToolsMapper) MCPToolsToOpenAI(mcpTools []mcpgo.Tool) ([]openai.Tool, error) {
	openaiTools := make([]openai.Tool, len(mcpTools))
	for i, mcpTool := range mcpTools {
		openaiTool, err := m.MCPToolToOpenAI(mcpTool)
		if err != nil {
			return nil, fmt.Errorf("failed to convert MCP tool %d to OpenAI: %w", i, err)
		}
		openaiTools[i] = openaiTool
	}
	return openaiTools, nil
}

// MCPToolResultToContent конвертирует результат выполнения MCP инструмента в контент для OpenAI.
func (m *ToolsMapper) MCPToolResultToContent(result mcpgo.CallToolResult) (string, error) {
	// Если это ошибка, форматируем как ошибку
	if result.IsError {
		if len(result.Content) > 0 {
			// Извлекаем текст из первого элемента контента
			text := mcpgo.GetTextFromContent(result.Content[0])
			return fmt.Sprintf("Error: %s", text), nil
		}
		return "Error: Unknown error occurred", nil
	}

	// Собираем весь контент в одну строку
	var contents []string
	for _, content := range result.Content {
		text := mcpgo.GetTextFromContent(content)
		if text != "" {
			contents = append(contents, text)
		}
	}

	if len(contents) == 0 {
		return "", nil
	}

	// Объединяем все части контента
	if len(contents) == 1 {
		return contents[0], nil
	}

	// Для множественного контента создаем JSON
	contentJSON, err := json.Marshal(contents)
	if err != nil {
		return "", fmt.Errorf("failed to marshal content: %w", err)
	}
	return string(contentJSON), nil
}

// convertOpenAIParametersToMCPSchema конвертирует OpenAI parameters в MCP ToolInputSchema.
func (m *ToolsMapper) convertOpenAIParametersToMCPSchema(parameters interface{}) (mcpgo.ToolInputSchema, error) {
	if parameters == nil {
		// Возвращаем пустую схему объекта
		return mcpgo.ToolInputSchema{
			Type:       "object",
			Properties: make(map[string]interface{}),
		}, nil
	}

	// Сначала конвертируем в JSON, затем обратно в map для нормализации
	jsonBytes, err := json.Marshal(parameters)
	if err != nil {
		return mcpgo.ToolInputSchema{}, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	var paramMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &paramMap); err != nil {
		return mcpgo.ToolInputSchema{}, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}

	// Создаем ToolInputSchema из map
	schema := mcpgo.ToolInputSchema{
		Type:       "object",
		Properties: make(map[string]interface{}),
	}

	// Извлекаем основные поля схемы
	if schemaType, ok := paramMap["type"].(string); ok {
		schema.Type = schemaType
	}

	if properties, ok := paramMap["properties"].(map[string]interface{}); ok {
		schema.Properties = properties
	}

	if required, ok := paramMap["required"].([]interface{}); ok {
		// Конвертируем []interface{} в []string
		stringRequired := make([]string, 0, len(required))
		for _, req := range required {
			if reqStr, ok := req.(string); ok {
				stringRequired = append(stringRequired, reqStr)
			}
		}
		schema.Required = stringRequired
	}

	if additionalProps, ok := paramMap["additionalProperties"]; ok {
		schema.AdditionalProperties = additionalProps
	}

	return schema, nil
}

// convertMCPSchemaToOpenAIParameters конвертирует MCP ToolInputSchema в OpenAI parameters.
func (m *ToolsMapper) convertMCPSchemaToOpenAIParameters(schema mcpgo.ToolInputSchema) (interface{}, error) {
	// Создаем map для OpenAI parameters
	result := map[string]interface{}{
		"type": schema.Type,
	}

	if len(schema.Properties) > 0 {
		result["properties"] = schema.Properties
	}

	if schema.Required != nil {
		result["required"] = schema.Required
	}

	if schema.AdditionalProperties != nil {
		result["additionalProperties"] = schema.AdditionalProperties
	}

	return result, nil
}

// generateToolCallID генерирует уникальный ID для tool call.
func generateToolCallID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "call_" + hex.EncodeToString(bytes)
}

// Helper функции для работы с MCP Content

// CreateTextContent создает текстовый контент для MCP.
func CreateTextContent(text string) mcpgo.Content {
	return mcpgo.NewTextContent(text)
}

// CreateErrorContent создает контент с ошибкой для MCP.
func CreateErrorContent(errorText string) mcpgo.CallToolResult {
	return mcpgo.CallToolResult{
		Content: []mcpgo.Content{
			mcpgo.NewTextContent(errorText),
		},
		IsError: true,
	}
}

// CreateSuccessContent создает успешный контент для MCP.
func CreateSuccessContent(text string) mcpgo.CallToolResult {
	return mcpgo.CallToolResult{
		Content: []mcpgo.Content{
			mcpgo.NewTextContent(text),
		},
		IsError: false,
	}
}
