package middleware

import (
    "log"
    "net/http"
    "time"
)

type ResponseWriter struct {
    http.ResponseWriter
    status      int
    wroteHeader bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
    return &ResponseWriter{ResponseWriter: w}
}

func (rw *ResponseWriter) Status() int {
    return rw.status
}

func (rw *ResponseWriter) WriteHeader(code int) {
    if rw.wroteHeader {
        return
    }

    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
    rw.wroteHeader = true
}

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // Логируем входящий запрос
            logger.Printf(
                "Входящий запрос | Method: %s | Path: %s | RemoteAddr: %s",
                r.Method,
                r.URL.Path,
                r.RemoteAddr,
            )

            // Создаем обертку для ResponseWriter чтобы отслеживать статус ответа
            wrapped := NewResponseWriter(w)

            // Обрабатываем запрос
            next.ServeHTTP(wrapped, r)

            // Логируем результат обработки запроса
            logger.Printf(
                "Исходящий ответ | Status: %d | Duration: %v | Path: %s",
                wrapped.Status(),
                time.Since(start),
                r.URL.Path,
            )
        })
    }
}
