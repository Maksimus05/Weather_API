package main

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "weather/internal/handlers"
)

func main() {
    log.Println("🌤️  Запуск Weather API...")
    
    // Получаем абсолютный путь к проекту
    projectRoot, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("📁 Рабочая директория: %s", projectRoot)
    
    // Проверяем API ключ
    apiKey := os.Getenv("OPENWEATHER_API_KEY")
    if apiKey == "" {
        log.Println("⚠️  OPENWEATHER_API_KEY не установлен")
    } else {
        log.Printf("✅ API ключ найден")
    }
    
    // Пути к файлам
    staticPath := filepath.Join(projectRoot, "..", "web", "static")
    templatePath := filepath.Join(projectRoot, "..", "web", "templates", "index.html")
    
    log.Printf("📁 Путь к статике: %s", staticPath)
    log.Printf("📁 Путь к шаблону: %s", templatePath)
    
    // Проверяем существование файлов
    if _, err := os.Stat(templatePath); os.IsNotExist(err) {
        log.Printf("❌ Файл не найден: %s", templatePath)
    } else {
        log.Printf("✅ HTML файл найден")
    }
    
    // Статические файлы
    http.Handle("/static/", 
        http.StripPrefix("/static/", 
            http.FileServer(http.Dir(staticPath))))
    
    // Главная страница - ПРАВИЛЬНЫЙ ПУТЬ
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("📄 Запрос главной страницы: %s", r.URL.Path)
        
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        
        // Пробуем разные пути к файлу
        possiblePaths := []string{
            templatePath,
            filepath.Join(projectRoot, "web", "templates", "index.html"),
            "../web/templates/index.html",
            "../../web/templates/index.html",
        }
        
        for _, path := range possiblePaths {
            if _, err := os.Stat(path); err == nil {
                log.Printf("✅ Отдаю файл: %s", path)
                http.ServeFile(w, r, path)
                return
            }
        }
        
        // Если файл не найден - выводим простую страницу
        log.Println("❌ HTML файл не найден")
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(`
            <!DOCTYPE html>
            <html>
            <head><title>Weather API</title></head>
            <body>
                <h1>🌤️ Weather API работает!</h1>
                <p>Но index.html не найден.</p>
                <p>Проверьте:</p>
                <ul>
                    <li><a href="/health">/health</a> - работает</li>
                    <li><a href="/weather?city=Moscow">/weather</a> - работает</li>
                </ul>
            </body>
            </html>
        `))
    })
    
    // API endpoints
    http.HandleFunc("/weather", handlers.WeatherHandler)
    http.HandleFunc("/health", handlers.HealthHandler)
    
    log.Println("\n" + strings.Repeat("=", 50))
    log.Println("✅ Сервер запущен: http://localhost:8080")
    log.Println("📡 /weather?city=Москва - получение погоды")
    log.Println("❤️  /health - проверка сервиса")
    log.Println(strings.Repeat("=", 50))
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}