package middleware

import (
	"bytes"
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func generateTraceID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}

func ObservabilityMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 1. Generar Trace ID para correlación
			traceID := r.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = generateTraceID()
			}

			// 2. Inyectar Trace ID en el contexto para que lo use el Repository/Mongo
			ctx := context.WithValue(r.Context(), "trace_id", traceID)
			r = r.WithContext(ctx)

			// 3. Preparar el capturador de respuesta
			rec := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           bytes.NewBuffer(nil),
			}

			// 4. Ejecutar el flujo de la API
			next.ServeHTTP(rec, r)

			// 5. Logging Asíncrono con Zap
			// Zap es extremadamente rápido, pero el procesamiento de strings (ofuscación)
			// es mejor hacerlo fuera del hilo principal de la petición.
			go func(duration time.Duration, status int, body []byte, path string, tid string) {
				logger.Info("HTTP Request Completed",
					zap.String("trace_id", tid),
					zap.String("path", path),
					zap.Int("status", status),
					zap.Duration("latency", duration),
					zap.String("response", maskSensitiveData(body)),
				)
			}(time.Since(start), rec.statusCode, rec.body.Bytes(), r.URL.Path, traceID)
		})
	}
}

func maskSensitiveData(data []byte) string {
	// Aquí podrías usar json.Unmarshal o strings.Replace
	// para ocultar campos como "password", "cvv", etc.
	if len(data) > 1000 { // Evitar logs gigantescos
		return "Body too large to log"
	}
	return string(data)
}
