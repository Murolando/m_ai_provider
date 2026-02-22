package mappers

import (
	"fmt"

	"github.com/Murolando/m_ai_provider/internal/entities/mcp"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

// MCPLibraryMapper предоставляет методы для конвертации между нашими MCP типами и типами библиотеки mark3labs/mcp-go.
type MCPLibraryMapper interface {
	// Конвертация наших типов в типы библиотеки
	OurToolToLibrary(ourTool mcp.Tool) (mcpgo.Tool, error)
	OurToolCallToLibrary(ourCall mcp.ToolCall) (mcpgo.CallToolRequest, error)
	OurToolResultToLibrary(ourResult mcp.ToolResult) (mcpgo.CallToolResult, error)
	OurContentToLibrary(ourContent mcp.Content) (mcpgo.Content, error)

	// Конвертация типов библиотеки в наши типы
	LibraryToolToOur(libTool mcpgo.Tool) (mcp.Tool, error)
	LibraryToolCallToOur(libCall mcpgo.CallToolRequest) (mcp.ToolCall, error)
	LibraryToolResultToOur(libResult mcpgo.CallToolResult) (mcp.ToolResult, error)
	LibraryContentToOur(libContent mcpgo.Content) (mcp.Content, error)

	// Массовые операции
	OurToolsToLibrary(ourTools []mcp.Tool) ([]mcpgo.Tool, error)
	LibraryToolsToOur(libTools []mcpgo.Tool) ([]mcp.Tool, error)
}

// mcpLibraryMapper реализует интерфейс MCPLibraryMapper.
type mcpLibraryMapper struct{}

// NewMCPLibraryMapper создает новый экземпляр маппера библиотеки MCP.
func NewMCPLibraryMapper() MCPLibraryMapper {
	return &mcpLibraryMapper{}
}

// OurToolToLibrary конвертирует наш Tool в Tool библиотеки.
func (m *mcpLibraryMapper) OurToolToLibrary(ourTool mcp.Tool) (mcpgo.Tool, error) {
	// Конвертируем InputSchema
	inputSchema, err := m.convertOurSchemaToLibrarySchema(ourTool.InputSchema)
	if err != nil {
		return mcpgo.Tool{}, fmt.Errorf("failed to convert input schema: %w", err)
	}

	// Получаем описание
	description := ""
	if ourTool.Description != nil {
		description = *ourTool.Description
	}

	// Создаем Tool библиотеки
	libTool := mcpgo.Tool{
		Name:        ourTool.Name,
		Description: description,
		InputSchema: mcpgo.ToolInputSchema(inputSchema),
	}

	return libTool, nil
}

// LibraryToolToOur конвертирует Tool библиотеки в наш Tool.
func (m *mcpLibraryMapper) LibraryToolToOur(libTool mcpgo.Tool) (mcp.Tool, error) {
	// Конвертируем InputSchema
	inputSchema, err := m.convertLibrarySchemaToOurSchema(mcpgo.ToolArgumentsSchema(libTool.InputSchema))
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("failed to convert input schema: %w", err)
	}

	// Создаем наш Tool
	var description *string
	if libTool.Description != "" {
		description = &libTool.Description
	}

	ourTool := mcp.Tool{
		Name:        libTool.Name,
		Description: description,
		InputSchema: inputSchema,
	}

	return ourTool, nil
}

// OurToolCallToLibrary конвертирует наш ToolCall в CallToolRequest библиотеки.
func (m *mcpLibraryMapper) OurToolCallToLibrary(ourCall mcp.ToolCall) (mcpgo.CallToolRequest, error) {
	// Создаем CallToolParams
	params := mcpgo.CallToolParams{
		Name:      ourCall.Name,
		Arguments: ourCall.Arguments,
	}

	// Создаем CallToolRequest
	libRequest := mcpgo.CallToolRequest{
		Params: params,
	}

	return libRequest, nil
}

// LibraryToolCallToOur конвертирует CallToolRequest библиотеки в наш ToolCall.
func (m *mcpLibraryMapper) LibraryToolCallToOur(libCall mcpgo.CallToolRequest) (mcp.ToolCall, error) {
	// Конвертируем Arguments с type assertion
	var arguments map[string]interface{}
	if libCall.Params.Arguments != nil {
		if args, ok := libCall.Params.Arguments.(map[string]interface{}); ok {
			arguments = args
		} else {
			return mcp.ToolCall{}, fmt.Errorf("failed to convert arguments: expected map[string]interface{}, got %T", libCall.Params.Arguments)
		}
	} else {
		arguments = make(map[string]interface{})
	}

	ourCall := mcp.ToolCall{
		Name:      libCall.Params.Name,
		Arguments: arguments,
	}

	return ourCall, nil
}

// OurToolResultToLibrary конвертирует наш ToolResult в CallToolResult библиотеки.
func (m *mcpLibraryMapper) OurToolResultToLibrary(ourResult mcp.ToolResult) (mcpgo.CallToolResult, error) {
	// Конвертируем Content
	var libContent []mcpgo.Content
	for _, content := range ourResult.Content {
		libContentItem, err := m.OurContentToLibrary(content)
		if err != nil {
			return mcpgo.CallToolResult{}, fmt.Errorf("failed to convert content: %w", err)
		}
		libContent = append(libContent, libContentItem)
	}

	// Создаем CallToolResult
	libResult := mcpgo.CallToolResult{
		Content: libContent,
		IsError: ourResult.IsError != nil && *ourResult.IsError,
	}

	return libResult, nil
}

// LibraryToolResultToOur конвертирует CallToolResult библиотеки в наш ToolResult.
func (m *mcpLibraryMapper) LibraryToolResultToOur(libResult mcpgo.CallToolResult) (mcp.ToolResult, error) {
	// Конвертируем Content
	var ourContent []mcp.Content
	for _, content := range libResult.Content {
		ourContentItem, err := m.LibraryContentToOur(content)
		if err != nil {
			return mcp.ToolResult{}, fmt.Errorf("failed to convert content: %w", err)
		}
		ourContent = append(ourContent, ourContentItem)
	}

	// Создаем наш ToolResult
	var isError *bool
	if libResult.IsError {
		isError = &libResult.IsError
	}

	ourResult := mcp.ToolResult{
		Content: ourContent,
		IsError: isError,
	}

	return ourResult, nil
}

// OurContentToLibrary конвертирует наш Content в Content библиотеки.
func (m *mcpLibraryMapper) OurContentToLibrary(ourContent mcp.Content) (mcpgo.Content, error) {
	switch ourContent.Type {
	case mcp.ContentTypeText:
		if ourContent.Text == nil {
			return nil, fmt.Errorf("text content must have text field")
		}
		return mcpgo.NewTextContent(*ourContent.Text), nil
	case mcp.ContentTypeImage:
		if ourContent.Data == nil {
			return nil, fmt.Errorf("image content must have data field")
		}
		// Для изображений используем AudioContent как базу, так как нет прямого ImageContent
		// В реальной реализации нужно будет создать правильный тип
		return mcpgo.NewAudioContent(*ourContent.Data, "image/png"), nil
	case mcp.ContentTypeResource:
		if ourContent.URI == nil {
			return nil, fmt.Errorf("resource content must have uri field")
		}
		// Для ресурсов создаем текстовый контент с URI
		return mcpgo.NewTextContent(*ourContent.URI), nil
	default:
		return nil, fmt.Errorf("unsupported content type: %s", ourContent.Type)
	}
}

// LibraryContentToOur конвертирует Content библиотеки в наш Content.
func (m *mcpLibraryMapper) LibraryContentToOur(libContent mcpgo.Content) (mcp.Content, error) {
	// Пытаемся привести к TextContent
	if textContent, ok := mcpgo.AsTextContent(libContent); ok {
		return mcp.Content{
			Type: mcp.ContentTypeText,
			Text: &textContent.Text,
		}, nil
	}

	// Пытаемся привести к AudioContent (может содержать изображения)
	if audioContent, ok := mcpgo.AsAudioContent(libContent); ok {
		return mcp.Content{
			Type: mcp.ContentTypeImage,
			Data: &audioContent.Data,
		}, nil
	}

	// Если не удалось определить тип, возвращаем как текст
	text := mcpgo.GetTextFromContent(libContent)
	return mcp.Content{
		Type: mcp.ContentTypeText,
		Text: &text,
	}, nil
}

// OurToolsToLibrary конвертирует массив наших Tools в массив Tools библиотеки.
func (m *mcpLibraryMapper) OurToolsToLibrary(ourTools []mcp.Tool) ([]mcpgo.Tool, error) {
	libTools := make([]mcpgo.Tool, len(ourTools))
	for i, ourTool := range ourTools {
		libTool, err := m.OurToolToLibrary(ourTool)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tool %d: %w", i, err)
		}
		libTools[i] = libTool
	}
	return libTools, nil
}

// LibraryToolsToOur конвертирует массив Tools библиотеки в массив наших Tools.
func (m *mcpLibraryMapper) LibraryToolsToOur(libTools []mcpgo.Tool) ([]mcp.Tool, error) {
	ourTools := make([]mcp.Tool, len(libTools))
	for i, libTool := range libTools {
		ourTool, err := m.LibraryToolToOur(libTool)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tool %d: %w", i, err)
		}
		ourTools[i] = ourTool
	}
	return ourTools, nil
}

// convertOurSchemaToLibrarySchema конвертирует нашу Schema в ToolArgumentsSchema библиотеки.
func (m *mcpLibraryMapper) convertOurSchemaToLibrarySchema(ourSchema mcp.Schema) (mcpgo.ToolArgumentsSchema, error) {
	// Конвертируем Properties
	properties := make(map[string]interface{})
	for propName, propValue := range ourSchema.Properties {
		properties[propName] = m.convertOurSchemaPropertyToMap(propValue)
	}

	// Создаем ToolArgumentsSchema
	libSchema := mcpgo.ToolArgumentsSchema{
		Type:       ourSchema.Type,
		Properties: properties,
		Required:   ourSchema.Required,
	}

	// Обрабатываем AdditionalProperties
	if ourSchema.AdditionalProperties != nil {
		libSchema.AdditionalProperties = *ourSchema.AdditionalProperties
	}

	return libSchema, nil
}

// convertLibrarySchemaToOurSchema конвертирует ToolArgumentsSchema библиотеки в нашу Schema.
func (m *mcpLibraryMapper) convertLibrarySchemaToOurSchema(libSchema mcpgo.ToolArgumentsSchema) (mcp.Schema, error) {
	// Конвертируем Properties
	properties := make(map[string]mcp.SchemaProperty)
	for propName, propValue := range libSchema.Properties {
		if propMap, ok := propValue.(map[string]interface{}); ok {
			property, err := m.convertMapToOurSchemaProperty(propMap)
			if err != nil {
				return mcp.Schema{}, fmt.Errorf("failed to convert property %s: %w", propName, err)
			}
			properties[propName] = property
		}
	}

	// Создаем нашу Schema
	ourSchema := mcp.Schema{
		Type:       libSchema.Type,
		Properties: properties,
		Required:   libSchema.Required,
	}

	// Обрабатываем AdditionalProperties
	if libSchema.AdditionalProperties != nil {
		if additionalProps, ok := libSchema.AdditionalProperties.(bool); ok {
			ourSchema.AdditionalProperties = &additionalProps
		}
	}

	return ourSchema, nil
}

// convertOurSchemaPropertyToMap конвертирует наш SchemaProperty в map для библиотеки.
func (m *mcpLibraryMapper) convertOurSchemaPropertyToMap(property mcp.SchemaProperty) map[string]interface{} {
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
		result["items"] = m.convertOurSchemaPropertyToMap(*property.Items)
	}

	if len(property.Properties) > 0 {
		properties := make(map[string]interface{})
		for nestedPropName, nestedPropValue := range property.Properties {
			properties[nestedPropName] = m.convertOurSchemaPropertyToMap(nestedPropValue)
		}
		result["properties"] = properties
	}

	if len(property.Required) > 0 {
		result["required"] = property.Required
	}

	return result
}

// convertMapToOurSchemaProperty конвертирует map в наш SchemaProperty.
func (m *mcpLibraryMapper) convertMapToOurSchemaProperty(propMap map[string]interface{}) (mcp.SchemaProperty, error) {
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
		itemProperty, err := m.convertMapToOurSchemaProperty(items)
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
				nestedProperty, err := m.convertMapToOurSchemaProperty(nestedPropMap)
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