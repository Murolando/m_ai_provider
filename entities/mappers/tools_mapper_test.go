package mappers

import (
	"encoding/json"
	"testing"

	"github.com/Murolando/m_ai_provider/internal/entities/openai"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

func TestOpenAIToolToMCP(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем тестовый OpenAI tool
	description := "Get the current weather in a given location"
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
						"description": "The city and state, e.g. San Francisco, CA",
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

	if mcpTool.Description != description {
		t.Errorf("Expected description '%s', got '%s'", description, mcpTool.Description)
	}

	// Проверяем схему
	inputSchema := mcpTool.InputSchema
	if inputSchema.Type != "object" {
		t.Errorf("Expected schema type 'object', got '%s'", inputSchema.Type)
	}

	// Проверяем свойства
	if len(inputSchema.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(inputSchema.Properties))
	}

	// Проверяем свойство location
	if locationProp, ok := inputSchema.Properties["location"].(map[string]interface{}); !ok {
		t.Error("Expected 'location' property to exist")
	} else if locationProp["type"] != "string" {
		t.Errorf("Expected location type 'string', got '%v'", locationProp["type"])
	}

	// Проверяем required
	if len(inputSchema.Required) != 1 || inputSchema.Required[0] != "location" {
		t.Errorf("Expected required field 'location', got %v", inputSchema.Required)
	}
}

func TestMCPToolToOpenAI(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем тестовый MCP tool
	description := "Search the web for information"
	mcpTool := mcpgo.Tool{
		Name:        "web_search",
		Description: description,
		InputSchema: mcpgo.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query",
				},
				"max_results": map[string]interface{}{
					"type":    "integer",
					"default": 10,
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
		t.Errorf("Expected description '%s'", description)
	}

	// Проверяем параметры
	params, ok := openaiTool.Function.Parameters.(map[string]interface{})
	if !ok {
		t.Fatal("Expected parameters to be a map")
	}

	if params["type"] != "object" {
		t.Errorf("Expected parameters type 'object', got '%v'", params["type"])
	}

	properties, ok := params["properties"].(map[string]interface{})
	if !ok || len(properties) != 2 {
		t.Error("Expected 2 properties in parameters")
	}
}

func TestOpenAIToolCallToMCP(t *testing.T) {
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
	if mcpToolCall.Params.Name != "get_weather" {
		t.Errorf("Expected name 'get_weather', got '%s'", mcpToolCall.Params.Name)
	}

	// Проверяем аргументы
	args, ok := mcpToolCall.Params.Arguments.(map[string]interface{})
	if !ok {
		t.Fatal("Expected arguments to be a map")
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(args))
	}

	if args["location"] != "Moscow" {
		t.Errorf("Expected location 'Moscow', got '%v'", args["location"])
	}

	if args["unit"] != "celsius" {
		t.Errorf("Expected unit 'celsius', got '%v'", args["unit"])
	}
}

func TestMCPToolCallToOpenAI(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем тестовый MCP tool call
	mcpToolCall := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Name: "web_search",
			Arguments: map[string]interface{}{
				"query":       "golang testing",
				"max_results": 5,
			},
		},
	}

	// Конвертируем в OpenAI
	openaiToolCall, err := mapper.MCPToolCallToOpenAI(mcpToolCall)
	if err != nil {
		t.Fatalf("Failed to convert MCP tool call to OpenAI: %v", err)
	}

	// Проверяем результат
	if openaiToolCall.Type != openai.ToolTypeFunction {
		t.Errorf("Expected type 'function', got '%s'", openaiToolCall.Type)
	}

	if openaiToolCall.Function.Name != "web_search" {
		t.Errorf("Expected name 'web_search', got '%s'", openaiToolCall.Function.Name)
	}

	// Проверяем ID (должен быть сгенерирован)
	if openaiToolCall.ID == "" {
		t.Error("Expected generated ID, got empty string")
	}

	// Проверяем аргументы
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(openaiToolCall.Function.Arguments), &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args["query"] != "golang testing" {
		t.Errorf("Expected query 'golang testing', got '%v'", args["query"])
	}

	// JSON unmarshal конвертирует числа в float64
	if maxResults, ok := args["max_results"].(float64); !ok || int(maxResults) != 5 {
		t.Errorf("Expected max_results 5, got '%v'", args["max_results"])
	}
}

func TestMCPToolResultToContent(t *testing.T) {
	mapper := NewToolsMapper()

	tests := []struct {
		name     string
		result   mcpgo.CallToolResult
		expected string
	}{
		{
			name: "success result",
			result: mcpgo.CallToolResult{
				Content: []mcpgo.Content{
					mcpgo.NewTextContent("Weather in Moscow: 20°C, sunny"),
				},
				IsError: false,
			},
			expected: "Weather in Moscow: 20°C, sunny",
		},
		{
			name: "error result",
			result: mcpgo.CallToolResult{
				Content: []mcpgo.Content{
					mcpgo.NewTextContent("City not found"),
				},
				IsError: true,
			},
			expected: "Error: City not found",
		},
		{
			name: "multiple content",
			result: mcpgo.CallToolResult{
				Content: []mcpgo.Content{
					mcpgo.NewTextContent("Line 1"),
					mcpgo.NewTextContent("Line 2"),
				},
				IsError: false,
			},
			expected: `["Line 1","Line 2"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := mapper.MCPToolResultToContent(tt.result)
			if err != nil {
				t.Fatalf("Failed to convert result to content: %v", err)
			}

			if content != tt.expected {
				t.Errorf("Expected content '%s', got '%s'", tt.expected, content)
			}
		})
	}
}

func TestOpenAIToolsToMCP(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем тестовые OpenAI tools
	description1 := "Tool 1"
	description2 := "Tool 2"
	openaiTools := []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: openai.Function{
				Name:        "tool1",
				Description: &description1,
				Parameters:  map[string]interface{}{"type": "object"},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: openai.Function{
				Name:        "tool2",
				Description: &description2,
				Parameters:  map[string]interface{}{"type": "object"},
			},
		},
	}

	// Конвертируем в MCP
	mcpTools, err := mapper.OpenAIToolsToMCP(openaiTools)
	if err != nil {
		t.Fatalf("Failed to convert OpenAI tools to MCP: %v", err)
	}

	// Проверяем результат
	if len(mcpTools) != 2 {
		t.Fatalf("Expected 2 tools, got %d", len(mcpTools))
	}

	if mcpTools[0].Name != "tool1" {
		t.Errorf("Expected first tool name 'tool1', got '%s'", mcpTools[0].Name)
	}

	if mcpTools[1].Name != "tool2" {
		t.Errorf("Expected second tool name 'tool2', got '%s'", mcpTools[1].Name)
	}
}

func TestMCPToolsToOpenAI(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем тестовые MCP tools
	mcpTools := []mcpgo.Tool{
		{
			Name:        "tool1",
			Description: "Tool 1",
			InputSchema: mcpgo.ToolInputSchema{Type: "object"},
		},
		{
			Name:        "tool2",
			Description: "Tool 2",
			InputSchema: mcpgo.ToolInputSchema{Type: "object"},
		},
	}

	// Конвертируем в OpenAI
	openaiTools, err := mapper.MCPToolsToOpenAI(mcpTools)
	if err != nil {
		t.Fatalf("Failed to convert MCP tools to OpenAI: %v", err)
	}

	// Проверяем результат
	if len(openaiTools) != 2 {
		t.Fatalf("Expected 2 tools, got %d", len(openaiTools))
	}

	if openaiTools[0].Function.Name != "tool1" {
		t.Errorf("Expected first tool name 'tool1', got '%s'", openaiTools[0].Function.Name)
	}

	if openaiTools[1].Function.Name != "tool2" {
		t.Errorf("Expected second tool name 'tool2', got '%s'", openaiTools[1].Function.Name)
	}
}

func TestRoundTripConversion(t *testing.T) {
	mapper := NewToolsMapper()

	// Создаем исходный OpenAI tool
	description := "Round trip test tool"
	originalTool := openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: openai.Function{
			Name:        "round_trip_tool",
			Description: &description,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"param1": map[string]interface{}{
						"type":        "string",
						"description": "String parameter",
					},
					"param2": map[string]interface{}{
						"type":    "integer",
						"minimum": 0,
						"maximum": 100,
					},
				},
				"required": []interface{}{"param1"},
			},
		},
	}

	// Конвертируем OpenAI -> MCP -> OpenAI
	mcpTool, err := mapper.OpenAIToolToMCP(originalTool)
	if err != nil {
		t.Fatalf("Failed to convert to MCP: %v", err)
	}

	resultTool, err := mapper.MCPToolToOpenAI(mcpTool)
	if err != nil {
		t.Fatalf("Failed to convert back to OpenAI: %v", err)
	}

	// Проверяем, что основные поля сохранились
	if resultTool.Function.Name != originalTool.Function.Name {
		t.Errorf("Name mismatch: expected '%s', got '%s'",
			originalTool.Function.Name, resultTool.Function.Name)
	}

	if *resultTool.Function.Description != *originalTool.Function.Description {
		t.Errorf("Description mismatch: expected '%s', got '%s'",
			*originalTool.Function.Description, *resultTool.Function.Description)
	}
}
