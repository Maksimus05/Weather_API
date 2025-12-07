package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
)

func main() {
    // 1. Проверяем существование файлов
    log.Println("🔍 Проверяем файлы...")
    
    if _, err := os.Stat("web/templates/index.html"); err != nil {
        log.Printf("❌ index.html не найден: %v", err)
    } else {
        log.Println("✅ index.html найден")
    }
    
    if _, err := os.Stat("web/static/style.css"); err != nil {
        log.Printf("❌ style.css не найден: %v", err)
    } else {
        log.Println("✅ style.css найден")
    }
    
    // 2. Настраиваем статические файлы
    http.Handle("/static/", 
        http.StripPrefix("/static/", 
            http.FileServer(http.Dir("web/static"))))
    
    // 3. Главная страница - ПРОСТО ОТДАЕМ ФАЙЛ
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        http.ServeFile(w, r, "web/templates/index.html")
    })
    
    // 4. API endpoints
    http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        city := r.URL.Query().Get("city")
        if city == "" {
            city = "Moscow"
        }
        
        json.NewEncoder(w).Encode(map[string]interface{}{
            "city":        city,
            "temp":        22.5,
            "description": "ясно",
            "success":     true,
        })
    })
    
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "status": "healthy",
        })
    })
    
    // 5. Запуск
    port := "8080"
    log.Println("\n" + strings.Repeat("=", 50))
    log.Println("✅ Weather API запущен!")
    log.Println("📍 http://localhost:" + port)
    log.Println("🎨 CSS: http://localhost:" + port + "/static/style.css")
    log.Println("📡 API: http://localhost:" + port + "/weather?city=Moscow")
    log.Println(strings.Repeat("=", 50))
    
    log.Fatal(http.ListenAndServe(":"+port, nil))
}