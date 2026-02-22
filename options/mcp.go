package options

import "github.com/Murolando/m_ai_provider/internal/entities/mcp"

// MCPToolsOption опция для передачи MCP инструментов.
// Каждый провайдер может по-своему обработать эти инструменты:
// - HydraAI: конвертирует в OpenAI формат и добавляет к запросу
// - Другие провайдеры: могут игнорировать или обрабатывать по-своему
type MCPToolsOption struct {
	Tools []mcp.Tool
}

// OptionType возвращает тип опции.
func (o MCPToolsOption) OptionType() string {
	return OptionTypeMCPTools
}

// GetTools возвращает MCP инструменты.
func (o MCPToolsOption) GetTools() []mcp.Tool {
	return o.Tools
}

// WithMCPTools создает опцию для добавления MCP инструментов.
func WithMCPTools(tools []mcp.Tool) SendMessageOption {
	return MCPToolsOption{Tools: tools}
}

// ExtractMCPToolsOption извлекает MCP tools опцию из списка опций.
// Возвращает инструменты и флаг найдена ли опция.
func ExtractMCPToolsOption(options []SendMessageOption) ([]mcp.Tool, bool) {
	for _, option := range options {
		if mcpOption, ok := option.(MCPToolsOption); ok {
			return mcpOption.GetTools(), true
		}
	}
	return nil, false
}