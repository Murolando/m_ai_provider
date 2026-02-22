package options

// SendMessageOption представляет интерфейс для опций отправки сообщений.
// Каждый провайдер может по-своему интерпретировать опции.
type SendMessageOption interface {
	// OptionType возвращает тип опции для идентификации провайдером
	OptionType() string
}

// Константы типов опций
const (
	// OptionTypeMCPTools тип опции для MCP инструментов
	OptionTypeMCPTools = "mcp_tools"
)