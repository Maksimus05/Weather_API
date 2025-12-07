package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
)

func main() {
    log.Println("🚀 Запускаем Weather API...")
    
    // Проверяем структуру
    log.Println("📁 Структура проекта:")
    log.Println("  web/templates/index.html - HTML страница")
    log.Println("  web/static/style.css - стили")
    
    // 1. Статические файлы
    http.Handle("/static/", 
        http.StripPrefix("/static/", 
            http.FileServer(http.Dir("web/static"))))
    
    // 2. Главная страница
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        http.ServeFile(w, r, "web/templates/index.html")
    })
    
    // 3. API эндпоинты
    http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        city := r.URL.Query().Get("city")
        if city == "" {
            city = "Moscow"
        }
        
        json.NewEncoder(w).Encode(map[string]interface{}{
            "city":        city,
            "temperature": 22.5,
            "feels_like":  21.8,
            "humidity":    65,
            "pressure":    1013,
            "wind_speed":  3.2,
            "description": "ясно",
            "success":     true,
        })
    })
    
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":  "healthy",
            "service": "weather-api",
            "success": true,
        })
    })
    
    // 4. Запуск сервера
    log.Println("\n" + strings.Repeat("=", 50))
    log.Println("✅ Weather API запущен!")
    log.Println("📍 Адрес: http://localhost:8080")
    log.Println(strings.Repeat("=", 50))
    log.Println("📡 Доступные эндпоинты:")
    log.Println("  GET /              - Главная страница")
    log.Println("  GET /weather?city= - Погода для города")
    log.Println("  GET /health        - Проверка здоровья")
    log.Println("  GET /static/       - CSS/JS файлы")
    log.Println(strings.Repeat("=", 50))
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
