package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

// Структуры для OpenWeather
type WeatherResponse struct {
    Main struct {
        Temp      float64 `json:"temp"`
        FeelsLike float64 `json:"feels_like"`
        Humidity  int     `json:"humidity"`
        Pressure  int     `json:"pressure"`
    } `json:"main"`
    Weather []struct {
        Description string `json:"description"`
        Icon        string `json:"icon"`
    } `json:"weather"`
    Wind struct {
        Speed float64 `json:"speed"`
    } `json:"wind"`
    Name string `json:"name"`
    Cod  int    `json:"cod"` // Код ответа от API
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
    // Устанавливаем Content-Type ДО отправки ответа
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    
    // Получаем город
    city := r.URL.Query().Get("city")
    if city == "" {
        city = "Moscow"
    }

    // Получаем API ключ
    apiKey := os.Getenv("OPENWEATHER_API_KEY")
    
    // Если ключа нет - возвращаем JSON с примером
    if apiKey == "" {
        showMockData(w, city)
        return
    }

    // Получаем реальные данные
    weather, err := getRealWeather(city, apiKey)
    if err != nil {
        // ВСЕГДА возвращаем JSON при ошибке
        jsonError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Отправляем успешный ответ
    json.NewEncoder(w).Encode(weather)
}

func getRealWeather(city, apiKey string) (map[string]interface{}, error) {
    // Формируем URL
    url := fmt.Sprintf(
        "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ru",
        city, apiKey,
    )

    // Делаем запрос
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("network error: %v", err)
    }
    defer resp.Body.Close()

    // Читаем ответ
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read error: %v", err)
    }

    // Парсим JSON
    var data WeatherResponse
    if err := json.Unmarshal(body, &data); err != nil {
        return nil, fmt.Errorf("parse error: %v", err)
    }

    // Проверяем код ответа от OpenWeather
    if data.Cod != 200 && data.Cod != 0 {
        return nil, fmt.Errorf("openweather error %d: %s", data.Cod, string(body))
    }

    // Форматируем ответ
    result := map[string]interface{}{
        "city":         data.Name,
        "temperature":  data.Main.Temp,
        "feels_like":   data.Main.FeelsLike,
        "humidity":     data.Main.Humidity,
        "pressure":     data.Main.Pressure,
        "wind_speed":   data.Wind.Speed,
        "description":  getDescription(data.Weather),
        "icon":         getIcon(data.Weather),
        "icon_url":     getIconURL(data.Weather),
        "units":        "metric",
        "source":       "OpenWeatherMap",
        "success":      true,
    }
    
    return result, nil
}

// Вспомогательные функции
func getDescription(weather []struct {
    Description string `json:"description"`
    Icon        string `json:"icon"`
}) string {
    if len(weather) > 0 {
        return weather[0].Description
    }
    return "нет данных"
}

func getIcon(weather []struct {
    Description string `json:"description"`
    Icon        string `json:"icon"`
}) string {
    if len(weather) > 0 {
        return weather[0].Icon
    }
    return "01d"
}

func getIconURL(weather []struct {
    Description string `json:"description"`
    Icon        string `json:"icon"`
}) string {
    icon := getIcon(weather)
    return fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", icon)
}

// Функция для возврата ошибок в JSON формате
func jsonError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error":   message,
        "success": false,
        "status":  statusCode,
    })
}

// Показываем примерные данные
func showMockData(w http.ResponseWriter, city string) {
    response := map[string]interface{}{
        "city":         city,
        "temperature":  22.5,
        "feels_like":   21.8,
        "humidity":     65,
        "pressure":     1013,
        "wind_speed":   3.2,
        "description":  "ясно",
        "icon":         "01d",
        "icon_url":     "https://openweathermap.org/img/wn/01d@2x.png",
        "units":        "metric",
        "source":       "Mock данные",
        "success":      true,
        "note":         "Установи OPENWEATHER_API_KEY для реальных данных",
        "how_to": []string{
            "1. Получи ключ на https://openweathermap.org/api",
            "2. Выполни в терминале:",
            "   export OPENWEATHER_API_KEY=твой_ключ",
            "3. Перезапусти сервер",
        },
    }
    
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    json.NewEncoder(w).Encode(response)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    apiKey := os.Getenv("OPENWEATHER_API_KEY")
    
    response := map[string]interface{}{
        "status":              "healthy",
        "service":             "weather-api",
        "api_key_configured":  apiKey != "",
        "success":             true,
    }
    
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    json.NewEncoder(w).Encode(response)
}