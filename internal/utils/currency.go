// Package utils содержит вспомогательные функции для работы с валютами и сообщениями.
package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// ValCurs представляет корневой элемент XML ответа от ЦБ РФ с курсами валют.
type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"` // Корневой XML элемент
	Date    string   `xml:"Date,attr"` // Дата курсов валют
	Valutes []Valute `xml:"Valute"`   // Массив валют
}

// Valute представляет информацию об одной валюте из API ЦБ РФ.
type Valute struct {
	ID       string `xml:"ID,attr"`   // Уникальный идентификатор валюты
	NumCode  string `xml:"NumCode"`   // Числовой код валюты
	CharCode string `xml:"CharCode"`  // Символьный код валюты (например, USD)
	Nominal  int    `xml:"Nominal"`   // Номинал валюты
	Name     string `xml:"Name"`      // Название валюты
	Value    string `xml:"Value"`     // Курс валюты к рублю
}

// GetUSDToRUBRate получает текущий курс доллара США к рублю от ЦБ РФ.
// Возвращает курс USD/RUB или ошибку при неудачном запросе.
func GetUSDToRUBRate() (float64, error) {
	req, err := http.NewRequest("GET", "https://www.cbr.ru/scripts/XML_daily.asp", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/xml, text/xml, */*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch currency rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	// Конвертируем из windows-1251 в UTF-8
	decoder := charmap.Windows1251.NewDecoder()
	utf8Body, err := io.ReadAll(transform.NewReader(bytes.NewReader(body), decoder))
	if err != nil {
		return 0, fmt.Errorf("failed to convert encoding: %w", err)
	}

	// Заменяем декларацию кодировки в XML на UTF-8
	utf8BodyStr := string(utf8Body)
	utf8BodyStr = strings.Replace(utf8BodyStr, `encoding="windows-1251"`, `encoding="UTF-8"`, 1)
	utf8Body = []byte(utf8BodyStr)

	var valCurs ValCurs
	if err := xml.Unmarshal(utf8Body, &valCurs); err != nil {
		return 0, fmt.Errorf("failed to parse XML: %w", err)
	}

	for _, valute := range valCurs.Valutes {
		if valute.CharCode == "USD" {
			valueStr := strings.Replace(valute.Value, ",", ".", -1)
			rate, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse USD rate: %w", err)
			}
			return rate, nil
		}
	}

	return 0, fmt.Errorf("USD rate not found")
}
