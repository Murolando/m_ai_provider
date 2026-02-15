package entities

// PricingParams интерфейс для различных типов параметров ценообразования.
// Пустой интерфейс для type switch в функциях calculatePrice.
type PricingParams interface{}

// HydraPricingParams содержит параметры для расчета цены в HydraAI провайдере.
type HydraPricingParams struct {
	Pricing HydraPricing `json:"pricing"` // Структура ценообразования от HydraAI API
}

// OpenRouterPricingParams содержит параметры для расчета цены в OpenRouter провайдере.
type OpenRouterPricingParams struct {
	PromptPrice     string `json:"prompt_price"`     // Цена за токены входящего сообщения
	CompletionPrice string `json:"completion_price"` // Цена за токены ответа
	RequestPrice    string `json:"request_price"`    // Цена за запрос
	ImagePrice      string `json:"image_price"`      // Цена за обработку изображений
}
