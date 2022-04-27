package chizerolog

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

func LoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.Error().
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						//Bytes("debug_stack", debug.Stack()).
						Msg("System Error")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// log end request
				if r.Header.Get("Content-Length") != "" && ww.BytesWritten() >= 0 {
					log.Info().
						Str("type", "access").
						Timestamp().
						Fields(map[string]interface{}{
							"remote_ip": r.RemoteAddr,
							"url":       r.URL.Path,
							"proto":     r.Proto,
							"method":    r.Method,
							//"user_agent": r.Header.Get("User-Agent"),
							"status":     ww.Status(),
							"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
							"bytes_in":   r.Header.Get("Content-Length"),
							"bytes_out":  ww.BytesWritten(),
						}).
						Msg("Request:")
				} else if r.Header.Get("Content-Length") == "" {
					log.Info().
						Str("type", "access").
						Timestamp().
						Fields(map[string]interface{}{
							"remote_ip": r.RemoteAddr,
							"url":       r.URL.Path,
							"proto":     r.Proto,
							"method":    r.Method,
							//"user_agent": r.Header.Get("User-Agent"),
							"status":     ww.Status(),
							"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
							"bytes_out":  ww.BytesWritten(),
						}).
						Msg("Request:")
				} else if ww.BytesWritten() == nil {
					log.Info().
						Str("type", "access").
						Timestamp().
						Fields(map[string]interface{}{
							"remote_ip": r.RemoteAddr,
							"url":       r.URL.Path,
							"proto":     r.Proto,
							"method":    r.Method,
							//"user_agent": r.Header.Get("User-Agent"),
							"status":     ww.Status(),
							"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
							"bytes_in":   r.Header.Get("Content-Length"),
						}).
						Msg("Request:")
				} else {
					log.Info().
						Str("type", "access").
						Timestamp().
						Fields(map[string]interface{}{
							"remote_ip": r.RemoteAddr,
							"url":       r.URL.Path,
							"proto":     r.Proto,
							"method":    r.Method,
							//"user_agent": r.Header.Get("User-Agent"),
							"status":     ww.Status(),
							"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						}).
						Msg("Request:")
				}
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
