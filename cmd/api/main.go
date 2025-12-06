package main

import (
    "html/template"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "weather/internal/handlers"
)

func main() {
    // Загружаем HTML шаблон
    tmpl, err := loadTemplate()
    if err != nil {
        log.Fatal("Ошибка загрузки шаблона:", err)
    }

    // Настраиваем обработчики
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        tmpl.Execute(w, nil)
    })

    // Статические файлы (CSS, JS)
    staticDir := http.Dir(filepath.Join("web", "static"))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

    // API эндпоинты
    http.HandleFunc("/weather", handlers.WeatherHandler)
    http.HandleFunc("/health", handlers.HealthHandler)

    // Порт
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Информация о запуске
    log.Println("=====================================")
    log.Println("🌤️  Weather API запущен!")
    log.Println("📍 Адрес: http://localhost:" + port)
    log.Println("📁 HTML шаблоны: web/templates/")
    log.Println("🎨 Статические файлы: web/static/")
    log.Println("=====================================")
    log.Println("📡 Доступные эндпоинты:")
    log.Println("   GET /              - Домашняя страница")
    log.Println("   GET /weather?city= - Погода для города")
    log.Println("   GET /health        - Проверка здоровья")
    log.Println("   GET /static/       - CSS/JS файлы")
    log.Println("=====================================")

    // Запускаем сервер
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("Ошибка запуска сервера:", err)
    }
}

func loadTemplate() (*template.Template, error) {
    // Ищем шаблон в разных местах (для гибкости)
    possiblePaths := []string{
        "web/templates/index.html",
        "../web/templates/index.html",
        "../../web/templates/index.html",
    }

    var tmpl *template.Template
    var err error

    for _, path := range possiblePaths {
        tmpl, err = template.ParseFiles(path)
        if err == nil {
            log.Printf("✅ Шаблон загружен: %s", path)
            return tmpl, nil
        }
    }

    // Если файл не найден, создаем простой HTML
    if err != nil {
        log.Printf("⚠️  Шаблон не найден, используется встроенный HTML")
        
        htmlContent := `
        <!DOCTYPE html>
        <html>
        <head><title>Weather API</title></head>
        <body>
            <h1>Weather API работает!</h1>
            <p>Шаблон не найден. Проверь папку web/templates/</p>
        </body>
        </html>`
        
        return template.New("index").Parse(htmlContent)
    }

    return tmpl, nil
}