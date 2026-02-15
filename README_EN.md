# M AI Provider
A unified entry point for integration with different AI systems.

All provider integration models conform to the interface in ```m_ai_provider/provider/provider.go```

## Provider List:
* [openrouter](https://openrouter.ai/) - active
* [hydraai](https://hydraai.app/) - active

## Usage Examples with Hydra AI

### Sending Messages

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
    // Get API key and URL from environment variables
    apiKey := os.Getenv("HYDRAAI_TOKEN")
    baseURL := os.Getenv("HYDRAAI_URL")

    // Create Hydra AI provider
    pr, err := provider.NewHydraAIProvider(apiKey, baseURL)
    if err != nil {
        log.Fatalf("Error creating provider: %v", err)
    }

    // Create messages to send
    messages := []*entities.Message{
        {
            ChatID:      "example-chat-123",
            MessageText: "Hello! How are you?",
            AuthorType:  entities.AuthorTypeUser,
            MessageType: entities.MessageText,
        },
        {
            ChatID:      "example-chat-123",
            MessageText: "Hello! I'm doing great, thank you! How can I help you?",
            AuthorType:  entities.AuthorTypeRobot,
            MessageType: entities.MessageText,
        },
        {
            ChatID:      "example-chat-123",
            MessageText: "Tell me an interesting fact about space",
            AuthorType:  entities.AuthorTypeUser,
            MessageType: entities.MessageText,
        },
    }

    // Select model (available models see in config/models.yaml)
    modelName := entities.ModelName("claude-3-5-haiku")

    // Send messages
    ctx := context.Background()
    response, err := pr.SendMessage(ctx, messages, modelName)
    if err != nil {
        log.Fatalf("Error sending message: %v", err)
    }

    // Output result
    fmt.Printf("Model response: %s\n", response.MessageText)
    fmt.Printf("Tokens used: %d\n", response.TotalTokens)
    fmt.Printf("Cost in rubles: %s\n", response.PriceInRubles.String())
}
```

### Getting Model Information (for comparison)

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
    // Create Hydra AI provider
    apiKey := os.Getenv("HYDRAAI_TOKEN")
    baseURL := os.Getenv("HYDRAAI_URL")
    
    pr, err := provider.NewHydraAIProvider(apiKey, baseURL)
    if err != nil {
        log.Fatalf("Error creating provider: %v", err)
    }

    // List of models for comparison
    modelsToCompare := []entities.ModelName{
        "claude-3-5-haiku",
        "claude-sonnet-4",
        "gpt-4o",
        "deepseek-v3",
        "gemini-2-0-flash",
    }

    fmt.Println("Hydra AI Models Comparison:")
    fmt.Println("==================================================")

    for _, modelName := range modelsToCompare {
        modelInfo, err := pr.GetModelInfo(modelName)
        if err != nil {
            fmt.Printf("❌ %s: %v\n", modelName, err)
            continue
        }

        fmt.Printf("✅ %s\n", modelInfo.Name)
        fmt.Printf("   Alias: %s\n", modelInfo.Alias)
        fmt.Printf("   Price: %s rubles\n", modelInfo.PriceInRubles.String())
        fmt.Println()
    }
}
```

### Choosing Provider for Model

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
    // Model we're looking for the best provider for
    targetModel := entities.ModelName("qwen-3-0-coder")

    // Create providers
    providers := make(map[string]provider.Provider)

    // Hydra AI provider
    if hydraToken := os.Getenv("HYDRAAI_TOKEN"); hydraToken != "" {
        if hydraURL := os.Getenv("HYDRAAI_URL"); hydraURL != "" {
            if hydraProvider, err := provider.NewHydraAIProvider(hydraToken, hydraURL); err == nil {
                providers["HydraAI"] = hydraProvider
            }
        }
    }

    // OpenRouter provider
    if openrouterToken := os.Getenv("OPENROUTER_TOKEN"); openrouterToken != "" {
        if openrouterProvider, err := provider.NewOpenRouterProvider(openrouterToken); err == nil {
            providers["OpenRouter"] = openrouterProvider
        }
    }

    if len(providers) == 0 {
        log.Fatal("Failed to create any provider. Check environment variables.")
    }

    bestProvider := ""
    bestPrice := decimal.NewFromFloat(999999) // Maximum price for comparison
    var bestProviderInstance provider.Provider

    // Check each provider
    for providerName, pr := range providers {
        modelInfo, err := pr.GetModelInfo(targetModel)
        if err != nil {
            fmt.Printf("❌ %s: model unavailable (%v)\n", providerName, err)
            continue
        }
        // Compare prices
        if modelInfo.PriceInRubles.LessThan(bestPrice) {
            bestPrice = modelInfo.PriceInRubles
            bestProvider = providerName
            bestProviderInstance = pr
        }
    }

    if bestProvider == "" {
        log.Fatal("Model unavailable from any provider")
    }

    fmt.Printf("🏆 Best provider: %s (price: %s rubles)\n", bestProvider, bestPrice.String())
    fmt.Println()

    // Use the best provider to send a message
    messages := []*entities.Message{
        {
            ChatID:      "test-chat",
            MessageText: "Write a simple Go function to add two numbers",
            AuthorType:  entities.AuthorTypeUser,
            MessageType: entities.MessageText,
        },
    }

    ctx := context.Background()
    response, err := bestProviderInstance.SendMessage(ctx, messages, targetModel)
    if err != nil {
        log.Fatalf("Error sending message through %s: %v", bestProvider, err)
    }
}
```

### Environment Variables Setup

To work with providers, you need to set the corresponding environment variables:

```bash
# For Hydra AI
export HYDRAAI_TOKEN="your-hydra-api-key"
export HYDRAAI_URL="https://api.hydraai.app/v1"

# For OpenRouter
export OPENROUTER_TOKEN="your-openrouter-api-key"
```

### Available Models

The complete list of supported models can be found in the [`config/models.yaml`](config/models.yaml) file. Some popular models:

- `claude-3-5-haiku` - fast and economical Claude model
- `claude-sonnet-4` - powerful Claude model for complex tasks
- `gpt-4o` - GPT-4 Omni model
- `deepseek-v3` - DeepSeek v3 model
- `gemini-2-0-flash` - fast Gemini model
- `qwen-3-32b` - Qwen 3 32B model