package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Murolando/m_ai_provider/entities"
	"github.com/Murolando/m_ai_provider/internal/config"
	internalEnt "github.com/Murolando/m_ai_provider/internal/entities"
	"github.com/Murolando/m_ai_provider/internal/entities/openai"
	"github.com/Murolando/m_ai_provider/internal/mappers"
	"github.com/Murolando/m_ai_provider/options"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/shopspring/decimal"
)

const (
	hydraAIProviderName = "HydraAI"
)

var (
	fakeCallID = "123321"
)

var _ Provider = (*HydraAIProvider)(nil)

// HydraAIProvider представляет провайдера для работы с HydraAI API.
type HydraAIProvider struct {
	apiKey      string                                     // API ключ для аутентификации
	baseURL     string                                     // Базовый URL для API запросов
	modelMap    map[entities.ModelName]*entities.ModelInfo // Кэш информации о моделях
	toolsMapper *mappers.ToolsMapper                       // Маппер для конвертации инструментов
}

// NewHydraAIProvider создает новый экземпляр HydraAI провайдера.
// apiKey - API ключ для аутентификации в HydraAI
// baseURL - базовый URL для API запросов
// Возвращает настроенный провайдер или ошибку при неудачной инициализации.
func NewHydraAIProvider(apiKey string, baseURL string) (*HydraAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("HYDRAAI_TOKEN is not set")
	}

	if baseURL == "" {
		return nil, fmt.Errorf("HYDRAAI_URL is not set")
	}

	provider := &HydraAIProvider{
		apiKey:      apiKey,
		baseURL:     baseURL,
		modelMap:    make(map[entities.ModelName]*entities.ModelInfo),
		toolsMapper: mappers.NewToolsMapper(),
	}

	if err := provider.getModels(); err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}

	return provider, nil
}

// SendMessage отправляет сообщения в AI модель через HydraAI API.
func (p *HydraAIProvider) SendMessage(ctx context.Context, messages []*entities.Message, modelName entities.ModelName, opts ...options.SendMessageOption) (*entities.ProviderMessageResponseDTO, error) {
	// Получаем модель из маппинга
	hydraModel, exists := config.HydraNamesMap[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s not supported by %s provider", modelName, hydraAIProviderName)
	}

	// Конвертируем сообщения в формат HydraAI
	chatMessages := convertToChatMessages(messages)
	request := internalEnt.NewHydraChatCompletionRequest(hydraModel, chatMessages)

	// Обрабатываем MCP tools опцию если она есть
	if mcpTools, hasMCPTools := options.ExtractMCPToolsOption(opts); hasMCPTools {
		// Конвертируем MCP tools в OpenAI формат
		openaiTools, err := p.toolsMapper.MCPToolsToOpenAI(mcpTools)
		if err != nil {
			return nil, fmt.Errorf("failed to convert MCP tools to OpenAI: %w", err)
		}

		// Добавляем инструменты к запросу
		request.Tools = openaiTools
		request.ToolChoice = openai.ToolChoiceAuto
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := p.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {

		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	// Выполняем запрос
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", response.StatusCode, string(responseBody))
	}

	var chatResponse internalEnt.HydraChatCompletionResponse
	if err := json.Unmarshal(responseBody, &chatResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := chatResponse.Choices[0]
	result := &entities.ProviderMessageResponseDTO{
		TotalTokens:   int64(chatResponse.Usage.TotalTokens),
		PriceInRubles: decimal.NewFromFloat(chatResponse.Usage.CostRequest).Round(3),
		FinishReason:  mapFinishReason(choice.FinishReason),
	}

	// Обрабатываем tool calls если они есть
	if len(choice.Message.ToolCalls) > 0 {
		mcpToolCalls := make([]mcpgo.CallToolRequest, len(choice.Message.ToolCalls))
		for i, toolCall := range choice.Message.ToolCalls {
			mcpToolCall, err := p.toolsMapper.OpenAIToolCallToMCP(toolCall)
			if err != nil {
				return nil, fmt.Errorf("failed to convert tool call %d to MCP: %w", i, err)
			}
			mcpToolCalls[i] = mcpToolCall
		}
		result.ToolCalls = mcpToolCalls
	}

	// Извлекаем текст ответа если есть
	switch content := choice.Message.Content.(type) {
	case string:
		result.MessageText = content
	case []interface{}:
		// Если это массив, ищем текстовые элементы
		for _, item := range content {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if itemType, ok := itemMap["type"].(string); ok && itemType == "text" {
					if text, ok := itemMap["text"].(string); ok {
						result.MessageText += text
					}
				}
			}
		}
	}

	return result, nil
}

// GetModelInfo получает информацию о конкретной модели из кэша.
func (p *HydraAIProvider) GetModelInfo(modelName entities.ModelName) (*entities.ModelInfo, error) {
	if modelInfo, exists := p.modelMap[modelName]; exists {
		return modelInfo, nil
	}
	return nil, fmt.Errorf("model %s not found in %s provider", modelName, hydraAIProviderName)
}

// getModels получает все модели от HydraAI API и заполняет кэш моделей.
func (p *HydraAIProvider) getModels() error {
	url := p.baseURL + "/models"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("API request failed with status %d: %s", response.StatusCode, string(body))
	}

	var modelsResponse internalEnt.ModelsResponse
	if err := json.NewDecoder(response.Body).Decode(&modelsResponse); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Проходим по всем моделям от API
	for _, hydraModel := range modelsResponse.Data {
		// Проверяем, есть ли эта модель в нашем маппинге
		for ourModelName, hydraModelID := range config.HydraNamesMap {
			if hydraModel.ID == hydraModelID && hydraModel.Active {
				// Рассчитываем цену
				pricingParams := internalEnt.HydraPricingParams{Pricing: hydraModel.Pricing}
				price, err := p.calculatePrice(pricingParams)
				if err != nil {
					// Если ошибка расчета, используем нулевую цену
					price = decimal.Zero
				}

				// Сохраняем в кэш
				p.modelMap[ourModelName] = &entities.ModelInfo{
					Name:          hydraModel.Name,
					Alias:         ourModelName,
					PriceInRubles: price,
				}
				break
			}
		}
	}

	return nil
}

// calculatePrice рассчитывает цену на основе переданных параметров
func (p *HydraAIProvider) calculatePrice(params internalEnt.PricingParams) (decimal.Decimal, error) {
	switch pricingParams := params.(type) {
	case internalEnt.HydraPricingParams:
		pricing := pricingParams.Pricing
		var basePrice float64

		switch pricing.Type {
		case "tokens":
			if pricing.InCostPerMillion != nil && pricing.OutCostPerMillion != nil {
				// Если есть два поля, берем сумму
				basePrice = *pricing.InCostPerMillion + *pricing.OutCostPerMillion
			} else if pricing.CostPerMillion != nil {
				// Если одно поле
				basePrice = *pricing.CostPerMillion
			}
		case "request":
			if pricing.CostPerRequest != nil {
				basePrice = *pricing.CostPerRequest
			}
		}

		// Округляем вверх до целого рубля
		finalPriceDecimal := decimal.NewFromFloat(basePrice)
		return finalPriceDecimal.Ceil(), nil
	default:
		return decimal.Zero, fmt.Errorf("unsupported pricing params type for HydraAI: %T", params)
	}
}

// convertToChatMessages конвертирует внутренние сообщения в формат OpenAI/HydraAI
func convertToChatMessages(messages []*entities.Message) []openai.ChatMessage {
	chatMessages := make([]openai.ChatMessage, len(messages))

	for i, msg := range messages {
		var role string
		switch msg.AuthorType {
		case entities.AuthorTypeUser:
			role = openai.RoleUser
		case entities.AuthorTypeRobot:
			role = openai.RoleAssistant
			chatMessages[i].ToolCallID = &fakeCallID
		case entities.AuthorTypeTool:
			role = openai.RoleTool
			chatMessages[i].ToolCallID = &fakeCallID
		default:
			role = openai.RoleUser // дефолтная роль
		}
		chatMessages[i] = openai.NewTextMessage(role, msg.MessageText)

		// Для сообщений с ролью tool нужно установить ToolCallID он не отдается в дефолтном mcpgo
		// поэтому обходить это буду путем создания уникальной пары значений
		if msg.AuthorType == entities.AuthorTypeTool && msg.ToolCallID != nil {
			chatMessages[i].ToolCallID = msg.ToolCallID
		}
	}

	return chatMessages
}

// mapFinishReason маппит OpenAI finish reason в общие константы entities.
func mapFinishReason(openaiReason *string) *string {
	if openaiReason == nil {
		return nil
	}

	var mappedReason string
	switch *openaiReason {
	case openai.FinishReasonStop:
		mappedReason = entities.FinishReasonStop
	case openai.FinishReasonLength:
		mappedReason = entities.FinishReasonLength
	case openai.FinishReasonToolCalls:
		mappedReason = entities.FinishReasonToolCalls
	case openai.FinishReasonContentFilter:
		mappedReason = entities.FinishReasonContentFilter
	default:
		// Если неизвестная причина, возвращаем как есть
		mappedReason = *openaiReason
	}

	return &mappedReason
}
