# M ai provider
Единая точка входа для интеграции с разными ИИ системами с поддержкой MCP tools.

Все модели интеграции к провайдерам соответствуют интерфейсу в ```m_ai_provider/provider/provider.go```

## Список провайдеров:
* [openrouter](https://openrouter.ai/) - active
* [hydraai](https://hydraai.app/) - active ✅ MCP tools support

## 🛠️ Поддержка MCP Tools

Проект поддерживает [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) от Anthropic для работы с инструментами. MCP tools автоматически конвертируются в OpenAI формат для совместимости с различными провайдерами(если это потребуется).

## Примеры использования с Hydra AI

### Отправка сообщений

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/Murolando/m_ai_provider/entities"
    "github.com/Murolando/m_ai_provider/provider"
)

func main() {
    // Получаем API ключ и URL из переменных окружения
    apiKey := os.Getenv("HYDRAAI_TOKEN")
    baseURL := os.Getenv("HYDRAAI_URL")

    // Создаем провайдер Hydra AI
    pr, err := provider.NewHydraAIProvider(apiKey, baseURL)
    if err != nil {
        log.Fatalf("Ошибка создания провайдера: %v", err)
    }

    // Создаем сообщения для отправки
    messages := []*entities.Message{
        {
            ChatID:      "example-chat-123",
            MessageText: "Привет! Как дела?",
            AuthorType:  entities.AuthorTypeUser,
            MessageType: entities.MessageText,
        },
        {
            ChatID:      "example-chat-123",
            MessageText: "Привет! У меня всё отлично, спасибо! Чем могу помочь?",
            AuthorType:  entities.AuthorTypeRobot,
            MessageType: entities.MessageText,
        },
        {
            ChatID:      "example-chat-123",
            MessageText: "Расскажи мне интересный факт о космосе",
            AuthorType:  entities.AuthorTypeUser,
            MessageType: entities.MessageText,
        },
    }

    // Выбираем модель (доступные модели см. в config/models.yaml)
    modelName := entities.ModelName("claude-3-5-haiku")

    // Отправляем сообщения
    ctx := context.Background()
    response, err := pr.SendMessage(ctx, messages, modelName)
    if err != nil {
        log.Fatalf("Ошибка отправки сообщения: %v", err)
    }

    // Выводим результат
    fmt.Printf("Ответ модели: %s\n", response.MessageText)
    fmt.Printf("Использовано токенов: %d\n", response.TotalTokens)
    fmt.Printf("Стоимость в рублях: %s\n", response.PriceInRubles.String())
}
```

### Получение информации о моделях (для сравнения)

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/Murolando/m_ai_provider/entities"
    "github.com/Murolando/m_ai_provider/provider"
)

func main() {
    // Создаем провайдер Hydra AI
    apiKey := os.Getenv("HYDRAAI_TOKEN")
    baseURL := os.Getenv("HYDRAAI_URL")
    
    pr, err := provider.NewHydraAIProvider(apiKey, baseURL)
    if err != nil {
        log.Fatalf("Ошибка создания провайдера: %v", err)
    }

    // Список моделей для сравнения
    modelsToCompare := []entities.ModelName{
        "claude-3-5-haiku",
        "claude-sonnet-4",
        "gpt-4o",
        "deepseek-v3",
        "gemini-2-0-flash",
    }

    fmt.Println("Сравнение моделей Hydra AI:")
    fmt.Println("=" * 50)

    for _, modelName := range modelsToCompare {
        modelInfo, err := pr.GetModelInfo(modelName)
        if err != nil {
            fmt.Printf("❌ %s: %v\n", modelName, err)
            continue
        }

        fmt.Printf("✅ %s\n", modelInfo.Name)
        fmt.Printf("   Алиас: %s\n", modelInfo.Alias)
        fmt.Printf("   Цена: %s руб.\n", modelInfo.PriceInRubles.String())
        fmt.Println()
    }
}
```

### Выбор провайдера для модели

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/Murolando/m_ai_provider/entities"
    "github.com/Murolando/m_ai_provider/provider"
    "github.com/shopspring/decimal"
)

func main() {
    // Модель, для которой ищем лучшего провайдера
    targetModel := entities.ModelName("qwen-3-0-coder")

    // Создаем провайдеры
    providers := make(map[string]provider.Provider)

    // Hydra AI провайдер
    if hydraToken := os.Getenv("HYDRAAI_TOKEN"); hydraToken != "" {
        if hydraURL := os.Getenv("HYDRAAI_URL"); hydraURL != "" {
            if hydraProvider, err := provider.NewHydraAIProvider(hydraToken, hydraURL); err == nil {
                providers["HydraAI"] = hydraProvider
            }
        }
    }

    // OpenRouter провайдер
    if openrouterToken := os.Getenv("OPENROUTER_TOKEN"); openrouterToken != "" {
        if openrouterProvider, err := provider.NewOpenRouterProvider(openrouterToken); err == nil {
            providers["OpenRouter"] = openrouterProvider
        }
    }

    if len(providers) == 0 {
        log.Fatal("Не удалось создать ни одного провайдера. Проверьте переменные окружения.")
    }

    bestProvider := ""
    bestPrice := decimal.NewFromFloat(999999) // Максимальная цена для сравнения
    var bestProviderInstance provider.Provider

    // Проверяем каждого провайдера
    for providerName, pr := range providers {
        modelInfo, err := pr.GetModelInfo(targetModel)
        if err != nil {
            fmt.Printf("❌ %s: модель недоступна (%v)\n", providerName, err)
            continue
        }
        // Сравниваем цены
        if modelInfo.PriceInRubles.LessThan(bestPrice) {
            bestPrice = modelInfo.PriceInRubles
            bestProvider = providerName
            bestProviderInstance = pr
        }
    }

    if bestProvider == "" {
        log.Fatal("Модель недоступна ни у одного провайдера")
    }

    fmt.Printf("🏆 Лучший провайдер: %s (цена: %s руб.)\n", bestProvider, bestPrice.String())
    fmt.Println()

    // Используем лучшего провайдера для отправки сообщения
    messages := []*entities.Message{
        {
            ChatID:      "test-chat",
            MessageText: "Напиши простую функцию на Go для сложения двух чисел",
            AuthorType:  entities.AuthorTypeUser,
            MessageType: entities.MessageText,
        },
    }

    ctx := context.Background()
    response, err := bestProviderInstance.SendMessage(ctx, messages, targetModel)
    if err != nil {
        log.Fatalf("Ошибка отправки сообщения через %s: %v", bestProvider, err)
    }
}
```

### Настройка переменных окружения

Для работы с провайдерами необходимо установить соответствующие переменные окружения:

```bash
# Для Hydra AI
export HYDRAAI_TOKEN="your-hydra-api-key"
export HYDRAAI_URL="https://api.hydraai.app/v1"

# Для OpenRouter
export OPENROUTER_TOKEN="your-openrouter-api-key"
```

### Доступные модели

Полный список поддерживаемых моделей можно найти в файле [`config/models.yaml`](config/models.yaml). Некоторые популярные модели:

- `claude-3-5-haiku` - быстрая и экономичная модель Claude
- `claude-sonnet-4` - мощная модель Claude для сложных задач
- `gpt-4o` - модель GPT-4 Omni
- `deepseek-v3` - модель DeepSeek v3
- `gemini-2-0-flash` - быстрая модель Gemini
- `qwen-3-32b` - модель Qwen 3 32B

## 🛠️ Использование MCP Tools

### Создание и использование MCP инструмента

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/Murolando/m_ai_provider/entities"
    "github.com/Murolando/m_ai_provider/internal/entities/mcp"
    "github.com/Murolando/m_ai_provider/options"
    "github.com/Murolando/m_ai_provider/provider"
)

func main() {
    // Создаем MCP инструмент для поиска в интернете
    description := "Search the web for information"
    queryDesc := "Search query"
    
    schema := mcp.NewSchema(mcp.SchemaTypeObject)
    schema.AddProperty("query", mcp.NewSchemaProperty(mcp.SchemaTypeString, &queryDesc))
    schema.AddRequired("query")
    
    webSearchTool := mcp.NewTool("web_search", &description, schema)

    // Создаем провайдер
    apiKey := os.Getenv("HYDRAAI_TOKEN")
    baseURL := os.Getenv("HYDRAAI_URL")
    
    provider, err := provider.NewHydraAIProvider(apiKey, baseURL)
    if err != nil {
        log.Fatalf("Ошибка создания провайдера: %v", err)
    }

    // Отправляем сообщение с MCP tools
    messages := []*entities.Message{
        {
            ChatID:      "example-chat",
            MessageText: "Найди информацию о последних новостях в области ИИ",
            AuthorType:  entities.AuthorTypeUser,
            MessageType: entities.MessageText,
        },
    }

    ctx := context.Background()
    response, err := provider.SendMessage(ctx, messages, "claude-3-5-haiku",
        options.WithMCPTools([]mcp.Tool{webSearchTool}))
    if err != nil {
        log.Fatalf("Ошибка отправки сообщения: %v", err)
    }

    // Проверяем, есть ли вызовы инструментов
    if len(response.ToolCalls) > 0 {
        fmt.Println("Модель вызвала инструменты:")
        for _, toolCall := range response.ToolCalls {
            fmt.Printf("- %s (ID: %s): %v\n", toolCall.Name, toolCall.ID, toolCall.Arguments)
            
            // Здесь можно выполнить реальный поиск и отправить результат обратно
            if toolCall.Name == "web_search" {
                query := toolCall.Arguments["query"].(string)
                fmt.Printf("Выполняем поиск: %s\n", query)
                // result := performWebSearch(query)
                // Отправить результат обратно в модель...
            }
        }
    }

    fmt.Printf("Ответ: %s\n", response.MessageText)
    fmt.Printf("Причина завершения: %s\n", *response.FinishReason)
}
```