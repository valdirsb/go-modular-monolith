package middleware

import (
    "net/http"
    "time"
    "log"
)

// LoggerMiddleware logs the details of each HTTP request
func LoggerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        
        next.ServeHTTP(w, r)
        
        log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
    })
}

// AuthMiddleware checks for authentication tokens in the request
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Here you would add logic to validate the token
        
        next.ServeHTTP(w, r)
    })
}