package mappers

import (
	"encoding/json"
	"testing"

	"github.com/Murolando/m_ai_provider/entities/mcp"
	"github.com/Murolando/m_ai_provider/internal/entities/openai"
)

func TestToolsMapper_OpenAIToolToMCP(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем тестовый OpenAI инструмент
	description := "Get current weather information"
	openaiTool := openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: openai.Function{
			Name:        "get_weather",
			Description: &description,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "City name",
					},
					"unit": map[string]interface{}{
						"type": "string",
						"enum": []interface{}{"celsius", "fahrenheit"},
					},
				},
				"required": []interface{}{"location"},
			},
		},
	}

	// Конвертируем в MCP
	mcpTool, err := mapper.OpenAIToolToMCP(openaiTool)
	if err != nil {
		t.Fatalf("Failed to convert OpenAI tool to MCP: %v", err)
	}

	// Проверяем результат
	if mcpTool.Name != "get_weather" {
		t.Errorf("Expected name 'get_weather', got '%s'", mcpTool.Name)
	}

	if mcpTool.Description == nil || *mcpTool.Description != description {
		t.Errorf("Expected description '%s', got %v", description, mcpTool.Description)
	}

	if mcpTool.InputSchema.Type != mcp.SchemaTypeObject {
		t.Errorf("Expected schema type 'object', got '%s'", mcpTool.InputSchema.Type)
	}

	// Проверяем свойства схемы
	if len(mcpTool.InputSchema.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(mcpTool.InputSchema.Properties))
	}

	locationProp, exists := mcpTool.InputSchema.Properties["location"]
	if !exists {
		t.Error("Expected 'location' property to exist")
	} else if locationProp.Type != mcp.SchemaTypeString {
		t.Errorf("Expected location type 'string', got '%s'", locationProp.Type)
	}

	// Проверяем обязательные поля
	if len(mcpTool.InputSchema.Required) != 1 || mcpTool.InputSchema.Required[0] != "location" {
		t.Errorf("Expected required fields ['location'], got %v", mcpTool.InputSchema.Required)
	}
}

func TestToolsMapper_MCPToolToOpenAI(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем тестовый MCP инструмент
	description := "Search the web for information"
	mcpTool := mcp.Tool{
		Name:        "web_search",
		Description: &description,
		InputSchema: mcp.Schema{
			Type: mcp.SchemaTypeObject,
			Properties: map[string]mcp.SchemaProperty{
				"query": {
					Type:        mcp.SchemaTypeString,
					Description: &[]string{"Search query"}[0],
				},
				"max_results": {
					Type:    mcp.SchemaTypeInteger,
					Default: 10,
				},
			},
			Required: []string{"query"},
		},
	}

	// Конвертируем в OpenAI
	openaiTool, err := mapper.MCPToolToOpenAI(mcpTool)
	if err != nil {
		t.Fatalf("Failed to convert MCP tool to OpenAI: %v", err)
	}

	// Проверяем результат
	if openaiTool.Type != openai.ToolTypeFunction {
		t.Errorf("Expected type 'function', got '%s'", openaiTool.Type)
	}

	if openaiTool.Function.Name != "web_search" {
		t.Errorf("Expected name 'web_search', got '%s'", openaiTool.Function.Name)
	}

	if openaiTool.Function.Description == nil || *openaiTool.Function.Description != description {
		t.Errorf("Expected description '%s', got %v", description, openaiTool.Function.Description)
	}

	// Проверяем параметры
	params, ok := openaiTool.Function.Parameters.(map[string]interface{})
	if !ok {
		t.Fatal("Expected parameters to be map[string]interface{}")
	}

	if params["type"] != "object" {
		t.Errorf("Expected type 'object', got '%v'", params["type"])
	}

	properties, ok := params["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected properties to be map[string]interface{}")
	}

	if len(properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(properties))
	}
}

func TestToolsMapper_OpenAIToolCallToMCP(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем тестовый OpenAI tool call
	openaiToolCall := openai.ToolCall{
		ID:   "call_123",
		Type: openai.ToolTypeFunction,
		Function: openai.FunctionCall{
			Name:      "get_weather",
			Arguments: `{"location": "Moscow", "unit": "celsius"}`,
		},
	}

	// Конвертируем в MCP
	mcpToolCall, err := mapper.OpenAIToolCallToMCP(openaiToolCall)
	if err != nil {
		t.Fatalf("Failed to convert OpenAI tool call to MCP: %v", err)
	}

	// Проверяем результат
	if mcpToolCall.Name != "get_weather" {
		t.Errorf("Expected name 'get_weather', got '%s'", mcpToolCall.Name)
	}

	if len(mcpToolCall.Arguments) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(mcpToolCall.Arguments))
	}

	if mcpToolCall.Arguments["location"] != "Moscow" {
		t.Errorf("Expected location 'Moscow', got '%v'", mcpToolCall.Arguments["location"])
	}

	if mcpToolCall.Arguments["unit"] != "celsius" {
		t.Errorf("Expected unit 'celsius', got '%v'", mcpToolCall.Arguments["unit"])
	}
}

func TestToolsMapper_MCPToolCallToOpenAI(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем тестовый MCP tool call
	mcpToolCall := mcp.ToolCall{
		Name: "web_search",
		Arguments: map[string]interface{}{
			"query":       "golang tutorial",
			"max_results": 5,
		},
	}

	// Конвертируем в OpenAI
	openaiToolCall, err := mapper.MCPToolCallToOpenAI(mcpToolCall)
	if err != nil {
		t.Fatalf("Failed to convert MCP tool call to OpenAI: %v", err)
	}

	// Проверяем результат (ID должен быть сгенерирован, так как в mcpToolCall его нет)
	if openaiToolCall.ID == "" {
		t.Error("Expected non-empty ID")
	}

	if openaiToolCall.Type != openai.ToolTypeFunction {
		t.Errorf("Expected type 'function', got '%s'", openaiToolCall.Type)
	}

	if openaiToolCall.Function.Name != "web_search" {
		t.Errorf("Expected name 'web_search', got '%s'", openaiToolCall.Function.Name)
	}

	// Проверяем аргументы
	var args map[string]interface{}
	err = json.Unmarshal([]byte(openaiToolCall.Function.Arguments), &args)
	if err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args["query"] != "golang tutorial" {
		t.Errorf("Expected query 'golang tutorial', got '%v'", args["query"])
	}

	if args["max_results"].(float64) != 5 {
		t.Errorf("Expected max_results 5, got '%v'", args["max_results"])
	}
}

func TestToolsMapper_RoundTrip(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем исходный OpenAI инструмент
	description := "Calculate sum of two numbers"
	originalTool := openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: openai.Function{
			Name:        "calculate_sum",
			Description: &description,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "number",
						"description": "First number",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "Second number",
					},
				},
				"required": []interface{}{"a", "b"},
			},
		},
	}

	// OpenAI -> MCP -> OpenAI
	mcpTool, err := mapper.OpenAIToolToMCP(originalTool)
	if err != nil {
		t.Fatalf("Failed to convert OpenAI to MCP: %v", err)
	}

	convertedTool, err := mapper.MCPToolToOpenAI(mcpTool)
	if err != nil {
		t.Fatalf("Failed to convert MCP to OpenAI: %v", err)
	}

	// Проверяем, что основные поля сохранились
	if convertedTool.Function.Name != originalTool.Function.Name {
		t.Errorf("Name mismatch: expected '%s', got '%s'",
			originalTool.Function.Name, convertedTool.Function.Name)
	}

	if *convertedTool.Function.Description != *originalTool.Function.Description {
		t.Errorf("Description mismatch: expected '%s', got '%s'",
			*originalTool.Function.Description, *convertedTool.Function.Description)
	}
}
