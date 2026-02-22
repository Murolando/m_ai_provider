package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Murolando/m_ai_provider/entities"
	"github.com/Murolando/m_ai_provider/internal/config"
	internalEnt "github.com/Murolando/m_ai_provider/internal/entities"
	"github.com/Murolando/m_ai_provider/internal/utils"
	"github.com/Murolando/m_ai_provider/options"
	"github.com/revrost/go-openrouter"
	"github.com/shopspring/decimal"
)

// Константы для провайдера OpenRouter
const (
	openRouterProviderName     = "OpenRouter"
	defaultUSDToRUBRateOnError = 80.0
)

// Проверяем, что OpenRouterProvider реализует интерфейс Provider
var _ Provider = (*OpenRouterProvider)(nil)

// OpenRouterProvider представляет провайдера для работы с OpenRouter API.
type OpenRouterProvider struct {
	client   *openrouter.Client                         // HTTP клиент для работы с OpenRouter API
	modelMap map[entities.ModelName]*entities.ModelInfo // Кэш информации о моделях
}

// NewOpenRouterProvider создает новый экземпляр OpenRouter провайдера.
// token - API токен для аутентификации в OpenRouter
// Возвращает настроенный провайдер или ошибку при неудачной инициализации.
func NewOpenRouterProvider(token string) (*OpenRouterProvider, error) {
	if token == "" {
		return nil, fmt.Errorf("OPENROUTER_TOKEN is not set")
	}

	client := openrouter.NewClient(token)

	provider := &OpenRouterProvider{
		client:   client,
		modelMap: make(map[entities.ModelName]*entities.ModelInfo),
	}

	if err := provider.getModels(); err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}

	return provider, nil
}

// SendMessage отправляет сообщения в AI модель через OpenRouter API.
func (p *OpenRouterProvider) SendMessage(ctx context.Context, messages []*entities.Message, modelName entities.ModelName, options ...options.SendMessageOption) (*entities.ProviderMessageResponseDTO, error) {
	openRouterModel, exists := config.OpenRouterNamesMap[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s not supported by %s provider", modelName, openRouterProviderName)
	}

	message := utils.MakeRequestMessageString(messages)
	response, err := p.client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
		Model: openRouterModel,
		Messages: []openrouter.ChatCompletionMessage{
			openrouter.UserMessage(message),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	var priceInRubles decimal.Decimal
	if response.Usage != nil {
		costUSD := response.Usage.Cost
		usdToRubRate, err := utils.GetUSDToRUBRate()
		if err != nil {
			usdToRubRate = defaultUSDToRUBRateOnError
		}
		costRUB := costUSD * usdToRubRate
		priceInRubles = decimal.NewFromFloat(costRUB).Round(3)
	}

	return &entities.ProviderMessageResponseDTO{
		MessageText:   response.Choices[0].Message.Content.Text,
		PriceInRubles: priceInRubles,
	}, nil
}

// GetModelInfo получает информацию о конкретной модели из кэша.
func (p *OpenRouterProvider) GetModelInfo(modelName entities.ModelName) (*entities.ModelInfo, error) {
	if modelInfo, exists := p.modelMap[modelName]; exists {
		return modelInfo, nil
	}
	return nil, fmt.Errorf("model %s not found in %s provider", modelName, openRouterProviderName)
}

// calculatePrice рассчитывает цену на основе параметров OpenRouter.
func (p *OpenRouterProvider) calculatePrice(params internalEnt.PricingParams) (decimal.Decimal, error) {
	switch pricingParams := params.(type) {
	case internalEnt.OpenRouterPricingParams:
		var basePriceUSD float64

		if promptPriceFloat, err := strconv.ParseFloat(pricingParams.PromptPrice, 64); err == nil {
			basePriceUSD += promptPriceFloat
		}
		if completionPriceFloat, err := strconv.ParseFloat(pricingParams.CompletionPrice, 64); err == nil {
			basePriceUSD += completionPriceFloat
		}
		if requestPriceFloat, err := strconv.ParseFloat(pricingParams.RequestPrice, 64); err == nil {
			basePriceUSD += requestPriceFloat
		}
		if imagePriceFloat, err := strconv.ParseFloat(pricingParams.ImagePrice, 64); err == nil {
			basePriceUSD += imagePriceFloat
		}

		usdToRubRate, err := utils.GetUSDToRUBRate()
		if err != nil {
			usdToRubRate = 100.0
		}

		basePriceRUB := basePriceUSD * usdToRubRate
		finalPriceDecimal := decimal.NewFromFloat(basePriceRUB)
		return finalPriceDecimal.Ceil(), nil
	default:
		return decimal.Zero, fmt.Errorf("unsupported pricing params type for OpenRouter: %T", params)
	}
}

// getModels получает все модели от OpenRouter API и заполняет кэш моделей.
func (p *OpenRouterProvider) getModels() error {
	ctx := context.Background()
	models, err := p.client.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	for _, model := range models {
		for ourModelName, openrouterModelID := range config.OpenRouterNamesMap {
			if model.ID == openrouterModelID {
				pricingParams := internalEnt.OpenRouterPricingParams{
					PromptPrice:     model.Pricing.Prompt,
					CompletionPrice: model.Pricing.Completion,
					RequestPrice:    model.Pricing.Request,
					ImagePrice:      model.Pricing.Image,
				}
				price, err := p.calculatePrice(pricingParams)
				if err != nil {
					// Если ошибка расчета, используем нулевую цену
					price = decimal.Zero
				}

				p.modelMap[ourModelName] = &entities.ModelInfo{
					Name:          model.Name,
					Alias:         ourModelName,
					PriceInRubles: price,
				}
				break
			}
		}
	}

	return nil
}
