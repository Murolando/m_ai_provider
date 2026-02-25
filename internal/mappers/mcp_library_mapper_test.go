package mappers

import (
	"testing"

	"github.com/Murolando/m_ai_provider/entities/mcp"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

func TestNewMCPLibraryMapper(t *testing.T) {
	mapper := NewMCPLibraryMapper()
	if mapper == nil {
		t.Fatal("NewMCPLibraryMapper() returned nil")
	}
}

func TestOurToolToLibrary(t *testing.T) {
	mapper := NewMCPLibraryMapper()

	// Создаем тестовый инструмент
	description := "Test tool description"
	ourTool := mcp.Tool{
		Name:        "test_tool",
		Description: &description,
		InputSchema: mcp.Schema{
			Type: mcp.SchemaTypeObject,
			Properties: map[string]mcp.SchemaProperty{
				"param1": {
					Type:        mcp.SchemaTypeString,
					Description: stringPtr("Test parameter"),
				},
			},
			Required: []string{"param1"},
		},
	}

	// Конвертируем
	libTool, err := mapper.OurToolToLibrary(ourTool)
	if err != nil {
		t.Fatalf("OurToolToLibrary() failed: %v", err)
	}

	// Проверяем результат
	if libTool.Name != ourTool.Name {
		t.Errorf("Expected name %s, got %s", ourTool.Name, libTool.Name)
	}

	if libTool.Description != *ourTool.Description {
		t.Errorf("Expected description %s, got %s", *ourTool.Description, libTool.Description)
	}

	// Проверяем схему
	inputSchema := mcpgo.ToolArgumentsSchema(libTool.InputSchema)
	if inputSchema.Type != mcp.SchemaTypeObject {
		t.Errorf("Expected schema type %s, got %s", mcp.SchemaTypeObject, inputSchema.Type)
	}

	if len(inputSchema.Required) != 1 || inputSchema.Required[0] != "param1" {
		t.Errorf("Expected required fields [param1], got %v", inputSchema.Required)
	}
}

func TestLibraryToolToOur(t *testing.T) {
	mapper := NewMCPLibraryMapper()

	// Создаем тестовый инструмент библиотеки
	libTool := mcpgo.Tool{
		Name:        "test_tool",
		Description: "Test tool description",
		InputSchema: mcpgo.ToolInputSchema{
			Type: mcp.SchemaTypeObject,
			Properties: map[string]interface{}{
				"param1": map[string]interface{}{
					"type":        mcp.SchemaTypeString,
					"description": "Test parameter",
				},
			},
			Required: []string{"param1"},
		},
	}

	// Конвертируем
	ourTool, err := mapper.LibraryToolToOur(libTool)
	if err != nil {
		t.Fatalf("LibraryToolToOur() failed: %v", err)
	}

	// Проверяем результат
	if ourTool.Name != libTool.Name {
		t.Errorf("Expected name %s, got %s", libTool.Name, ourTool.Name)
	}

	if ourTool.Description == nil || *ourTool.Description != libTool.Description {
		t.Errorf("Expected description %s, got %v", libTool.Description, ourTool.Description)
	}

	// Проверяем схему
	if ourTool.InputSchema.Type != mcp.SchemaTypeObject {
		t.Errorf("Expected schema type %s, got %s", mcp.SchemaTypeObject, ourTool.InputSchema.Type)
	}

	if len(ourTool.InputSchema.Required) != 1 || ourTool.InputSchema.Required[0] != "param1" {
		t.Errorf("Expected required fields [param1], got %v", ourTool.InputSchema.Required)
	}
}

func TestOurToolCallToLibrary(t *testing.T) {
	mapper := NewMCPLibraryMapper()

	// Создаем тестовый вызов инструмента
	ourCall := mcp.ToolCall{
		ID:   "test_id",
		Name: "test_tool",
		Arguments: map[string]interface{}{
			"param1": "value1",
			"param2": 42,
		},
	}

	// Конвертируем
	libRequest, err := mapper.OurToolCallToLibrary(ourCall)
	if err != nil {
		t.Fatalf("OurToolCallToLibrary() failed: %v", err)
	}

	// Проверяем результат
	if libRequest.Params.Name != ourCall.Name {
		t.Errorf("Expected name %s, got %s", ourCall.Name, libRequest.Params.Name)
	}

	// Проверяем аргументы
	if len(libRequest.Params.Arguments.(map[string]interface{})) != len(ourCall.Arguments) {
		t.Errorf("Expected %d arguments, got %d", len(ourCall.Arguments), len(libRequest.Params.Arguments.(map[string]interface{})))
	}
}

func TestLibraryToolCallToOur(t *testing.T) {
	mapper := NewMCPLibraryMapper()

	// Создаем тестовый запрос библиотеки
	libRequest := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Name: "test_tool",
			Arguments: map[string]interface{}{
				"param1": "value1",
				"param2": 42,
			},
		},
	}

	// Конвертируем
	ourCall, err := mapper.LibraryToolCallToOur(libRequest)
	if err != nil {
		t.Fatalf("LibraryToolCallToOur() failed: %v", err)
	}

	// Проверяем результат
	if ourCall.Name != libRequest.Params.Name {
		t.Errorf("Expected name %s, got %s", libRequest.Params.Name, ourCall.Name)
	}

	// Проверяем аргументы
	if len(ourCall.Arguments) != len(libRequest.Params.Arguments.(map[string]interface{})) {
		t.Errorf("Expected %d arguments, got %d", len(libRequest.Params.Arguments.(map[string]interface{})), len(ourCall.Arguments))
	}
}

func TestOurContentToLibrary(t *testing.T) {
	mapper := NewMCPLibraryMapper()

	tests := []struct {
		name        string
		ourContent  mcp.Content
		expectError bool
	}{
		{
			name: "text content",
			ourContent: mcp.Content{
				Type: mcp.ContentTypeText,
				Text: stringPtr("Hello, world!"),
			},
			expectError: false,
		},
		{
			name: "image content",
			ourContent: mcp.Content{
				Type: mcp.ContentTypeImage,
				Data: stringPtr("base64imagedata"),
			},
			expectError: false,
		},
		{
			name: "resource content",
			ourContent: mcp.Content{
				Type: mcp.ContentTypeResource,
				URI:  stringPtr("https://example.com/resource"),
			},
			expectError: false,
		},
		{
			name: "invalid text content",
			ourContent: mcp.Content{
				Type: mcp.ContentTypeText,
				Text: nil,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			libContent, err := mapper.OurContentToLibrary(tt.ourContent)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("OurContentToLibrary() failed: %v", err)
			}

			if libContent == nil {
				t.Errorf("Expected non-nil content")
			}
		})
	}
}

func TestLibraryContentToOur(t *testing.T) {
	mapper := NewMCPLibraryMapper()

	// Тестируем текстовый контент
	textContent := mcpgo.NewTextContent("Hello, world!")
	ourContent, err := mapper.LibraryContentToOur(textContent)
	if err != nil {
		t.Fatalf("LibraryContentToOur() failed: %v", err)
	}

	if ourContent.Type != mcp.ContentTypeText {
		t.Errorf("Expected type %s, got %s", mcp.ContentTypeText, ourContent.Type)
	}

	if ourContent.Text == nil || *ourContent.Text != "Hello, world!" {
		t.Errorf("Expected text 'Hello, world!', got %v", ourContent.Text)
	}
}

func TestOurToolsToLibrary(t *testing.T) {
	mapper := NewMCPLibraryMapper()

	// Создаем массив тестовых инструментов
	description1 := "Tool 1"
	description2 := "Tool 2"
	ourTools := []mcp.Tool{
		{
			Name:        "tool1",
			Description: &description1,
			InputSchema: mcp.Schema{Type: mcp.SchemaTypeObject},
		},
		{
			Name:        "tool2",
			Description: &description2,
			InputSchema: mcp.Schema{Type: mcp.SchemaTypeObject},
		},
	}

	// Конвертируем
	libTools, err := mapper.OurToolsToLibrary(ourTools)
	if err != nil {
		t.Fatalf("OurToolsToLibrary() failed: %v", err)
	}

	// Проверяем результат
	if len(libTools) != len(ourTools) {
		t.Errorf("Expected %d tools, got %d", len(ourTools), len(libTools))
	}

	for i, libTool := range libTools {
		if libTool.Name != ourTools[i].Name {
			t.Errorf("Tool %d: expected name %s, got %s", i, ourTools[i].Name, libTool.Name)
		}
	}
}

func TestLibraryToolsToOur(t *testing.T) {
	mapper := NewMCPLibraryMapper()

	// Создаем массив тестовых инструментов библиотеки
	libTools := []mcpgo.Tool{
		{
			Name:        "tool1",
			Description: "Tool 1",
			InputSchema: mcpgo.ToolInputSchema{Type: mcp.SchemaTypeObject},
		},
		{
			Name:        "tool2",
			Description: "Tool 2",
			InputSchema: mcpgo.ToolInputSchema{Type: mcp.SchemaTypeObject},
		},
	}

	// Конвертируем
	ourTools, err := mapper.LibraryToolsToOur(libTools)
	if err != nil {
		t.Fatalf("LibraryToolsToOur() failed: %v", err)
	}

	// Проверяем результат
	if len(ourTools) != len(libTools) {
		t.Errorf("Expected %d tools, got %d", len(libTools), len(ourTools))
	}

	for i, ourTool := range ourTools {
		if ourTool.Name != libTools[i].Name {
			t.Errorf("Tool %d: expected name %s, got %s", i, libTools[i].Name, ourTool.Name)
		}
	}
}

func TestRoundTripConversion(t *testing.T) {
	mapper := NewMCPLibraryMapper()

	// Создаем оригинальный инструмент
	description := "Round trip test tool"
	originalTool := mcp.Tool{
		Name:        "round_trip_tool",
		Description: &description,
		InputSchema: mcp.Schema{
			Type: mcp.SchemaTypeObject,
			Properties: map[string]mcp.SchemaProperty{
				"param1": {
					Type:        mcp.SchemaTypeString,
					Description: stringPtr("String parameter"),
				},
				"param2": {
					Type:    mcp.SchemaTypeInteger,
					Minimum: float64Ptr(0),
					Maximum: float64Ptr(100),
				},
			},
			Required: []string{"param1"},
		},
	}

	// Конвертируем наш -> библиотека -> наш
	libTool, err := mapper.OurToolToLibrary(originalTool)
	if err != nil {
		t.Fatalf("OurToolToLibrary() failed: %v", err)
	}

	convertedTool, err := mapper.LibraryToolToOur(libTool)
	if err != nil {
		t.Fatalf("LibraryToolToOur() failed: %v", err)
	}

	// Проверяем, что основные поля сохранились
	if convertedTool.Name != originalTool.Name {
		t.Errorf("Name mismatch: expected %s, got %s", originalTool.Name, convertedTool.Name)
	}

	if convertedTool.Description == nil || *convertedTool.Description != *originalTool.Description {
		t.Errorf("Description mismatch: expected %s, got %v", *originalTool.Description, convertedTool.Description)
	}

	if convertedTool.InputSchema.Type != originalTool.InputSchema.Type {
		t.Errorf("Schema type mismatch: expected %s, got %s", originalTool.InputSchema.Type, convertedTool.InputSchema.Type)
	}

	if len(convertedTool.InputSchema.Required) != len(originalTool.InputSchema.Required) {
		t.Errorf("Required fields mismatch: expected %v, got %v", originalTool.InputSchema.Required, convertedTool.InputSchema.Required)
	}
}

// Вспомогательные функции для создания указателей
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
