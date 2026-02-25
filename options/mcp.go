package options

import mcpgo "github.com/mark3labs/mcp-go/mcp"

// MCPToolsOption представляет опцию для добавления MCP инструментов к запросу.
// MCP (Model Context Protocol) - стандартизированный протокол для работы с инструментами AI.
// Документация: https://modelcontextprotocol.io/
type MCPToolsOption struct {
	Tools []mcpgo.Tool
}

// OptionType возвращает тип опции для идентификации провайдером.
func (o MCPToolsOption) OptionType() string {
	return OptionTypeMCPTools
}

// GetTools возвращает MCP инструменты.
func (o MCPToolsOption) GetTools() []mcpgo.Tool {
	return o.Tools
}

// WithMCPTools создает опцию для добавления MCP инструментов.
func WithMCPTools(tools []mcpgo.Tool) SendMessageOption {
	return MCPToolsOption{Tools: tools}
}

// ExtractMCPToolsOption извлекает опцию MCP инструментов из списка опций.
// Возвращает инструменты и флаг найдена ли опция.
func ExtractMCPToolsOption(options []SendMessageOption) ([]mcpgo.Tool, bool) {
	for _, option := range options {
		if mcpOption, ok := option.(MCPToolsOption); ok {
			return mcpOption.GetTools(), true
		}
	}
	return nil, false
}

// MCPToolCallsOption представляет опцию для добавления вызовов MCP инструментов.
type MCPToolCallsOption struct {
	ToolCalls []mcpgo.CallToolRequest
}

// OptionType возвращает тип опции для идентификации провайдером.
func (o MCPToolCallsOption) OptionType() string {
	return "mcp_tool_calls"
}

// GetToolCalls возвращает вызовы MCP инструментов.
func (o MCPToolCallsOption) GetToolCalls() []mcpgo.CallToolRequest {
	return o.ToolCalls
}

// WithMCPToolCalls создает опцию для добавления вызовов MCP инструментов.
func WithMCPToolCalls(toolCalls []mcpgo.CallToolRequest) SendMessageOption {
	return MCPToolCallsOption{ToolCalls: toolCalls}
}

// ExtractMCPToolCallsOption извлекает опцию вызовов MCP инструментов из списка опций.
func ExtractMCPToolCallsOption(options []SendMessageOption) ([]mcpgo.CallToolRequest, bool) {
	for _, option := range options {
		if mcpOption, ok := option.(MCPToolCallsOption); ok {
			return mcpOption.GetToolCalls(), true
		}
	}
	return nil, false
}
