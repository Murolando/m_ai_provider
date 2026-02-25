package mappers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Murolando/m_ai_provider/entities/mcp"
	"github.com/Murolando/m_ai_provider/internal/entities/openai"
)

// ToolsMapper предоставляет методы для конвертации между OpenAI и MCP форматами инструментов.
type ToolsMapper struct{}

// NewToolsMapper создает новый экземпляр маппера инструментов.
func NewToolsMapper() *ToolsMapper {
	return &ToolsMapper{}
}

// OpenAIToolToMCP конвертирует OpenAI инструмент в MCP формат.
func (m *ToolsMapper) OpenAIToolToMCP(openaiTool openai.Tool) (mcp.Tool, error) {
	// Конвертируем JSON Schema из interface{} в mcp.Schema
	inputSchema, err := m.convertOpenAIParametersToMCPSchema(openaiTool.Function.Parameters)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("failed to convert parameters to MCP schema: %w", err)
	}

	return mcp.Tool{
		Name:        openaiTool.Function.Name,
		Description: openaiTool.Function.Description,
		InputSchema: inputSchema,
	}, nil
}

// MCPToolToOpenAI конвертирует MCP инструмент в OpenAI формат.
func (m *ToolsMapper) MCPToolToOpenAI(mcpTool mcp.Tool) (openai.Tool, error) {
	// Конвертируем MCP Schema в interface{} для OpenAI
	parameters, err := m.convertMCPSchemaToOpenAIParameters(mcpTool.InputSchema)
	if err != nil {
		return openai.Tool{}, fmt.Errorf("failed to convert MCP schema to OpenAI parameters: %w", err)
	}

	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: openai.Function{
			Name:        mcpTool.Name,
			Description: mcpTool.Description,
			Parameters:  parameters,
		},
	}, nil
}

// OpenAIToolCallToMCP конвертирует OpenAI вызов инструмента в MCP формат.
func (m *ToolsMapper) OpenAIToolCallToMCP(openaiToolCall openai.ToolCall) (mcp.ToolCall, error) {
	// Парсим JSON строку аргументов в map
	var arguments map[string]interface{}
	if openaiToolCall.Function.Arguments != "" {
		if err := json.Unmarshal([]byte(openaiToolCall.Function.Arguments), &arguments); err != nil {
			return mcp.ToolCall{}, fmt.Errorf("failed to parse OpenAI tool call arguments: %w", err)
		}
	} else {
		arguments = make(map[string]interface{})
	}

	return mcp.ToolCall{
		ID:        openaiToolCall.ID, // Сохраняем ID для связи с результатом
		Name:      openaiToolCall.Function.Name,
		Arguments: arguments,
	}, nil
}

// MCPToolCallToOpenAI конвертирует MCP вызов инструмента в OpenAI формат.
func (m *ToolsMapper) MCPToolCallToOpenAI(mcpToolCall mcp.ToolCall) (openai.ToolCall, error) {
	// Конвертируем map аргументов в JSON строку
	argumentsJSON, err := json.Marshal(mcpToolCall.Arguments)
	if err != nil {
		return openai.ToolCall{}, fmt.Errorf("failed to marshal MCP tool call arguments: %w", err)
	}

	// Используем ID из MCP tool call, если он есть, иначе генерируем
	id := mcpToolCall.ID
	if id == "" {
		id = generateToolCallID()
	}

	return openai.ToolCall{
		ID:   id,
		Type: openai.ToolTypeFunction,
		Function: openai.FunctionCall{
			Name:      mcpToolCall.Name,
			Arguments: string(argumentsJSON),
		},
	}, nil
}

// OpenAIToolsToMCP конвертирует массив OpenAI инструментов в MCP формат.
func (m *ToolsMapper) OpenAIToolsToMCP(openaiTools []openai.Tool) ([]mcp.Tool, error) {
	mcpTools := make([]mcp.Tool, len(openaiTools))
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
func (m *ToolsMapper) MCPToolsToOpenAI(mcpTools []mcp.Tool) ([]openai.Tool, error) {
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

// convertOpenAIParametersToMCPSchema конвертирует OpenAI parameters в MCP Schema.
func (m *ToolsMapper) convertOpenAIParametersToMCPSchema(parameters interface{}) (mcp.Schema, error) {
	if parameters == nil {
		return mcp.Schema{Type: mcp.SchemaTypeObject}, nil
	}

	// Сначала конвертируем в JSON, затем обратно в map для нормализации
	jsonBytes, err := json.Marshal(parameters)
	if err != nil {
		return mcp.Schema{}, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	var paramMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &paramMap); err != nil {
		return mcp.Schema{}, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}

	return m.convertMapToMCPSchema(paramMap)
}

// convertMCPSchemaToOpenAIParameters конвертирует MCP Schema в OpenAI parameters.
func (m *ToolsMapper) convertMCPSchemaToOpenAIParameters(schema mcp.Schema) (interface{}, error) {
	return m.convertMCPSchemaToMap(schema), nil
}

// convertMapToMCPSchema конвертирует map в MCP Schema.
func (m *ToolsMapper) convertMapToMCPSchema(paramMap map[string]interface{}) (mcp.Schema, error) {
	schema := mcp.Schema{
		Type:       mcp.SchemaTypeObject,
		Properties: make(map[string]mcp.SchemaProperty),
	}

	// Извлекаем основные поля схемы
	if schemaType, ok := paramMap["type"].(string); ok {
		schema.Type = schemaType
	}

	if properties, ok := paramMap["properties"].(map[string]interface{}); ok {
		for propName, propValue := range properties {
			if propMap, ok := propValue.(map[string]interface{}); ok {
				property, err := m.convertMapToSchemaProperty(propMap)
				if err != nil {
					return schema, fmt.Errorf("failed to convert property %s: %w", propName, err)
				}
				schema.Properties[propName] = property
			}
		}
	}

	if required, ok := paramMap["required"].([]interface{}); ok {
		for _, req := range required {
			if reqStr, ok := req.(string); ok {
				schema.Required = append(schema.Required, reqStr)
			}
		}
	}

	if additionalProps, ok := paramMap["additionalProperties"].(bool); ok {
		schema.AdditionalProperties = &additionalProps
	}

	return schema, nil
}

// convertMapToSchemaProperty конвертирует map в SchemaProperty.
func (m *ToolsMapper) convertMapToSchemaProperty(propMap map[string]interface{}) (mcp.SchemaProperty, error) {
	property := mcp.SchemaProperty{}

	if propType, ok := propMap["type"].(string); ok {
		property.Type = propType
	}

	if description, ok := propMap["description"].(string); ok {
		property.Description = &description
	}

	if enum, ok := propMap["enum"].([]interface{}); ok {
		property.Enum = enum
	}

	if defaultValue, ok := propMap["default"]; ok {
		property.Default = defaultValue
	}

	if minimum, ok := propMap["minimum"].(float64); ok {
		property.Minimum = &minimum
	}

	if maximum, ok := propMap["maximum"].(float64); ok {
		property.Maximum = &maximum
	}

	if minLength, ok := propMap["minLength"].(float64); ok {
		minLengthInt := int(minLength)
		property.MinLength = &minLengthInt
	}

	if maxLength, ok := propMap["maxLength"].(float64); ok {
		maxLengthInt := int(maxLength)
		property.MaxLength = &maxLengthInt
	}

	// Обработка items для массивов
	if items, ok := propMap["items"].(map[string]interface{}); ok {
		itemProperty, err := m.convertMapToSchemaProperty(items)
		if err != nil {
			return property, fmt.Errorf("failed to convert items property: %w", err)
		}
		property.Items = &itemProperty
	}

	// Обработка вложенных properties для объектов
	if properties, ok := propMap["properties"].(map[string]interface{}); ok {
		property.Properties = make(map[string]mcp.SchemaProperty)
		for nestedPropName, nestedPropValue := range properties {
			if nestedPropMap, ok := nestedPropValue.(map[string]interface{}); ok {
				nestedProperty, err := m.convertMapToSchemaProperty(nestedPropMap)
				if err != nil {
					return property, fmt.Errorf("failed to convert nested property %s: %w", nestedPropName, err)
				}
				property.Properties[nestedPropName] = nestedProperty
			}
		}
	}

	// Обработка required для вложенных объектов
	if required, ok := propMap["required"].([]interface{}); ok {
		for _, req := range required {
			if reqStr, ok := req.(string); ok {
				property.Required = append(property.Required, reqStr)
			}
		}
	}

	return property, nil
}

// convertMCPSchemaToMap конвертирует MCP Schema в map для OpenAI.
func (m *ToolsMapper) convertMCPSchemaToMap(schema mcp.Schema) map[string]interface{} {
	result := map[string]interface{}{
		"type": schema.Type,
	}

	if len(schema.Properties) > 0 {
		properties := make(map[string]interface{})
		for propName, propValue := range schema.Properties {
			properties[propName] = m.convertSchemaPropertyToMap(propValue)
		}
		result["properties"] = properties
	}

	if len(schema.Required) > 0 {
		result["required"] = schema.Required
	}

	if schema.AdditionalProperties != nil {
		result["additionalProperties"] = *schema.AdditionalProperties
	}

	return result
}

// generateToolCallID генерирует уникальный ID для tool call.
func generateToolCallID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// convertSchemaPropertyToMap конвертирует SchemaProperty в map.
func (m *ToolsMapper) convertSchemaPropertyToMap(property mcp.SchemaProperty) map[string]interface{} {
	result := make(map[string]interface{})

	if property.Type != "" {
		result["type"] = property.Type
	}

	if property.Description != nil {
		result["description"] = *property.Description
	}

	if len(property.Enum) > 0 {
		result["enum"] = property.Enum
	}

	if property.Default != nil {
		result["default"] = property.Default
	}

	if property.Minimum != nil {
		result["minimum"] = *property.Minimum
	}

	if property.Maximum != nil {
		result["maximum"] = *property.Maximum
	}

	if property.MinLength != nil {
		result["minLength"] = *property.MinLength
	}

	if property.MaxLength != nil {
		result["maxLength"] = *property.MaxLength
	}

	if property.Items != nil {
		result["items"] = m.convertSchemaPropertyToMap(*property.Items)
	}

	if len(property.Properties) > 0 {
		properties := make(map[string]interface{})
		for nestedPropName, nestedPropValue := range property.Properties {
			properties[nestedPropName] = m.convertSchemaPropertyToMap(nestedPropValue)
		}
		result["properties"] = properties
	}

	if len(property.Required) > 0 {
		result["required"] = property.Required
	}

	return result
}
