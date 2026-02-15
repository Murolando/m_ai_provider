package utils

import (
	"strings"

	"github.com/Murolando/m_ai_provider/entities"
)

// MakeRequestMessageString объединяет массив сообщений в одну строку.
// messages - массив сообщений для объединения
// Возвращает строку, где каждое сообщение разделено символом новой строки.
func MakeRequestMessageString(messages []*entities.Message) string {
	var result strings.Builder
	for _, item := range messages {
		result.WriteString(item.MessageText)
		result.WriteString("\n")
	}
	return result.String()
}
