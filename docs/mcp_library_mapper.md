# MCP Library Mapper

Маппер для конвертации между нашими внутренними MCP типами и типами библиотеки `github.com/mark3labs/mcp-go`.

## Обзор

`MCPLibraryMapper` предоставляет двунаправленную конвертацию между:
- Нашими типами (`github.com/Murolando/m_ai_provider/internal/entities/mcp`)
- Типами библиотеки (`github.com/mark3labs/mcp-go/mcp`)

## Установка

Добавьте зависимость в ваш `go.mod`:

```go
require github.com/mark3labs/mcp-go v0.44.0
```

## Использование

### Создание маппера

```go
import (
    "github.com/Murolando/m_ai_provider/internal/mappers"
)

mapper := mappers.NewMCPLibraryMapper()
```

### Конвертация инструментов

#### Наш Tool → Library Tool

```go
import (
    "github.com/Murolando/m_ai_provider/internal/entities/mcp"
    mcpgo "github.com/mark3labs/mcp-go/mcp"
)

// Создаем наш инструмент
description := "Получает текущую погоду для указанного города"
ourTool := mcp.Tool{
    Name:        "get_weather",
    Description: &description,
    InputSchema: mcp.Schema{
        Type: mcp.SchemaTypeObject,
        Properties: map[string]mcp.SchemaProperty{
            "city": {
                Type:        mcp.SchemaTypeString,
                Description: stringPtr("Название города"),
            },
            "units": {
                Type:        mcp.SchemaTypeString,
                Description: stringPtr("Единицы измерения температуры"),
                Enum:        []interface{}{"celsius", "fahrenheit"},
                Default:     "celsius",
            },
        },
        Required: []string{"city"},
    },
}

// Конвертируем в тип библиотеки
libTool, err := mapper.OurToolToLibrary(ourTool)
if err != nil {
    log.Fatalf("Ошибка конвертации: %v", err)
}

fmt.Printf("Library tool: %+v\n", libTool)
```

#### Library Tool → Наш Tool

```go
// Создаем инструмент библиотеки
libTool := mcpgo.Tool{
    Name:        "calculate_sum",
    Description: "Вычисляет сумму двух чисел",
    InputSchema: mcpgo.ToolInputSchema{
        Type: "object",
        Properties: map[string]interface{}{
            "a": map[string]interface{}{
                "type":        "number",
                "description": "Первое число",
            },
            "b": map[string]interface{}{
                "type":        "number", 
                "description": "Второе число",
            },
        },
        Required: []string{"a", "b"},
    },
}

// Конвертируем в наш тип
ourTool, err := mapper.LibraryToolToOur(libTool)
if err != nil {
    log.Fatalf("Ошибка конвертации: %v", err)
}

fmt.Printf("Our tool: %+v\n", ourTool)
```

### Конвертация вызовов инструментов

#### Наш ToolCall → Library CallToolRequest

```go
// Создаем наш вызов инструмента
ourCall := mcp.ToolCall{
    ID:   "call_123",
    Name: "get_weather",
    Arguments: map[string]interface{}{
        "city":  "Москва",
        "units": "celsius",
    },
}

// Конвертируем в запрос библиотеки
libRequest, err := mapper.OurToolCallToLibrary(ourCall)
if err != nil {
    log.Fatalf("Ошибка конвертации: %v", err)
}

fmt.Printf("Library request: %+v\n", libRequest)
```

#### Library CallToolRequest → Наш ToolCall

```go
// Создаем запрос библиотеки
libRequest := mcpgo.CallToolRequest{
    Params: mcpgo.CallToolParams{
        Name: "calculate_sum",
        Arguments: map[string]interface{}{
            "a": 10.5,
            "b": 20.3,
        },
    },
}

// Конвертируем в наш вызов
ourCall, err := mapper.LibraryToolCallToOur(libRequest)
if err != nil {
    log.Fatalf("Ошибка конвертации: %v", err)
}

fmt.Printf("Our call: %+v\n", ourCall)
```

### Конвертация результатов

#### Наш ToolResult → Library CallToolResult

```go
// Создаем наш результат
ourResult := mcp.ToolResult{
    Content: []mcp.Content{
        {
            Type: mcp.ContentTypeText,
            Text: stringPtr("Температура в Москве: 15°C, облачно"),
        },
    },
    IsError: nil, // успешный результат
}

// Конвертируем в результат библиотеки
libResult, err := mapper.OurToolResultToLibrary(ourResult)
if err != nil {
    log.Fatalf("Ошибка конвертации: %v", err)
}

fmt.Printf("Library result: %+v\n", libResult)
```

#### Library CallToolResult → Наш ToolResult

```go
// Создаем результат библиотеки
libResult := mcpgo.CallToolResult{
    Content: []mcpgo.Content{
        mcpgo.NewTextContent("Сумма: 30.8"),
    },
    IsError: false,
}

// Конвертируем в наш результат
ourResult, err := mapper.LibraryToolResultToOur(libResult)
if err != nil {
    log.Fatalf("Ошибка конвертации: %v", err)
}

fmt.Printf("Our result: %+v\n", ourResult)
```

### Массовые операции

#### Конвертация массива инструментов

```go
// Наши инструменты → Инструменты библиотеки
ourTools := []mcp.Tool{tool1, tool2, tool3}
libTools, err := mapper.OurToolsToLibrary(ourTools)
if err != nil {
    log.Fatalf("Ошибка конвертации: %v", err)
}

// Инструменты библиотеки → Наши инструменты  
libTools := []mcpgo.Tool{libTool1, libTool2, libTool3}
ourTools, err := mapper.LibraryToolsToOur(libTools)
if err != nil {
    log.Fatalf("Ошибка конвертации: %v", err)
}
```

## Поддерживаемые типы содержимого

### Текстовое содержимое

```go
// Наш тип
textContent := mcp.Content{
    Type: mcp.ContentTypeText,
    Text: stringPtr("Привет, мир!"),
}

// Конвертация
libContent, err := mapper.OurContentToLibrary(textContent)
```

### Содержимое изображения

```go
// Наш тип
imageContent := mcp.Content{
    Type: mcp.ContentTypeImage,
    Data: stringPtr("base64encodedimagedata"),
}

// Конвертация
libContent, err := mapper.OurContentToLibrary(imageContent)
```

### Ресурсное содержимое

```go
// Наш тип
resourceContent := mcp.Content{
    Type: mcp.ContentTypeResource,
    URI:  stringPtr("https://example.com/resource"),
}

// Конвертация
libContent, err := mapper.OurContentToLibrary(resourceContent)
```

## Обработка ошибок

Все методы маппера возвращают ошибку в случае проблем с конвертацией:

```go
libTool, err := mapper.OurToolToLibrary(ourTool)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "failed to convert input schema"):
        log.Printf("Ошибка конвертации схемы: %v", err)
    case strings.Contains(err.Error(), "failed to convert"):
        log.Printf("Общая ошибка конвертации: %v", err)
    default:
        log.Printf("Неизвестная ошибка: %v", err)
    }
    return
}
```

## Полный пример

```go
package main

import (
    "fmt"
    "log"

    "github.com/Murolando/m_ai_provider/internal/entities/mcp"
    "github.com/Murolando/m_ai_provider/internal/mappers"
    mcpgo "github.com/mark3labs/mcp-go/mcp"
)

func main() {
    // Создаем маппер
    mapper := mappers.NewMCPLibraryMapper()

    // Создаем наш инструмент
    description := "Конвертирует температуру между шкалами"
    ourTool := mcp.Tool{
        Name:        "convert_temperature",
        Description: &description,
        InputSchema: mcp.Schema{
            Type: mcp.SchemaTypeObject,
            Properties: map[string]mcp.SchemaProperty{
                "value": {
                    Type:        mcp.SchemaTypeNumber,
                    Description: stringPtr("Значение температуры"),
                },
                "from": {
                    Type:        mcp.SchemaTypeString,
                    Description: stringPtr("Исходная шкала"),
                    Enum:        []interface{}{"celsius", "fahrenheit", "kelvin"},
                },
                "to": {
                    Type:        mcp.SchemaTypeString,
                    Description: stringPtr("Целевая шкала"),
                    Enum:        []interface{}{"celsius", "fahrenheit", "kelvin"},
                },
            },
            Required: []string{"value", "from", "to"},
        },
    }

    // Конвертируем в тип библиотеки
    libTool, err := mapper.OurToolToLibrary(ourTool)
    if err != nil {
        log.Fatalf("Ошибка конвертации инструмента: %v", err)
    }

    fmt.Printf("Конвертированный инструмент библиотеки:\n")
    fmt.Printf("  Название: %s\n", libTool.Name)
    fmt.Printf("  Описание: %s\n", libTool.Description)

    // Создаем вызов инструмента
    ourCall := mcp.ToolCall{
        ID:   "temp_call_1",
        Name: "convert_temperature",
        Arguments: map[string]interface{}{
            "value": 25.0,
            "from":  "celsius",
            "to":    "fahrenheit",
        },
    }

    // Конвертируем вызов
    libRequest, err := mapper.OurToolCallToLibrary(ourCall)
    if err != nil {
        log.Fatalf("Ошибка конвертации вызова: %v", err)
    }

    fmt.Printf("\nКонвертированный вызов библиотеки:\n")
    fmt.Printf("  Инструмент: %s\n", libRequest.Params.Name)
    fmt.Printf("  Аргументы: %+v\n", libRequest.Params.Arguments)

    // Создаем результат
    ourResult := mcp.ToolResult{
        Content: []mcp.Content{
            {
                Type: mcp.ContentTypeText,
                Text: stringPtr("25°C = 77°F"),
            },
        },
    }

    // Конвертируем результат
    libResult, err := mapper.OurToolResultToLibrary(ourResult)
    if err != nil {
        log.Fatalf("Ошибка конвертации результата: %v", err)
    }

    fmt.Printf("\nКонвертированный результат библиотеки:\n")
    fmt.Printf("  Содержимое: %d элементов\n", len(libResult.Content))
    fmt.Printf("  Ошибка: %t\n", libResult.IsError)

    // Тестируем round-trip конвертацию
    convertedBackTool, err := mapper.LibraryToolToOur(libTool)
    if err != nil {
        log.Fatalf("Ошибка обратной конвертации: %v", err)
    }

    fmt.Printf("\nПроверка round-trip конвертации:\n")
    fmt.Printf("  Исходное название: %s\n", ourTool.Name)
    fmt.Printf("  Конвертированное название: %s\n", convertedBackTool.Name)
    fmt.Printf("  Совпадают: %t\n", ourTool.Name == convertedBackTool.Name)
}

// Вспомогательная функция
func stringPtr(s string) *string {
    return &s
}
```

## Ограничения

1. **Типы содержимого**: Библиотека может не поддерживать все типы содержимого напрямую
2. **Метаданные**: Некоторые метаданные могут быть потеряны при конвертации
3. **Валидация**: Маппер выполняет базовую валидацию, но не все ограничения схемы

## Тестирование

Запустите тесты для проверки функциональности:

```bash
go test ./internal/mappers/ -v
```

Тесты покрывают:
- Базовую конвертацию всех типов
- Round-trip конвертацию
- Обработку ошибок
- Граничные случаи