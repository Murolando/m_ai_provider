package provider

import (
	"context"

	"github.com/Murolando/m_ai_provider/entities"
	internalEnt "github.com/Murolando/m_ai_provider/internal/entities"
	"github.com/Murolando/m_ai_provider/internal/utils"
	"github.com/Murolando/m_ai_provider/options"
	"github.com/shopspring/decimal"
)

var _ Provider = (*DefaultProvider)(nil)

// DefaultProvider представляет провайдера-заглушку для тестирования и разработки.
type DefaultProvider struct{}

// NewDefaultProvider создает новый экземпляр DefaultProvider.
func NewDefaultProvider() *DefaultProvider {
	return &DefaultProvider{}
}

// SendMessage отправляет сообщения через DefaultProvider (возвращает тестовый ответ).
func (p *DefaultProvider) SendMessage(ctx context.Context, messages []*entities.Message, modelName entities.ModelName, options ...options.SendMessageOption) (*entities.ProviderMessageResponseDTO, error) {
	message := utils.MakeRequestMessageString(messages)
	return &entities.ProviderMessageResponseDTO{
		MessageText: "DEFAULT ANSWER FOR " + message,
	}, nil
}

// GetModelInfo получает информацию о модели (всегда возвращает nil для DefaultProvider).
func (p *DefaultProvider) GetModelInfo(modelName entities.ModelName) (*entities.ModelInfo, error) {
	return nil, nil
}

// calculatePrice рассчитывает цену (DefaultProvider всегда возвращает нулевую цену).
func (p *DefaultProvider) calculatePrice(params internalEnt.PricingParams) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

// getModels загружает модели (DefaultProvider не загружает модели).
func (p *DefaultProvider) getModels() error {
	return nil
}
