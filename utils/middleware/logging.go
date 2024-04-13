package middleware

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	logName         string
	ConnectorLogger *log.Logger
)

func SetupLogger() {

	logger := log.New()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if logName == "" {
		logName = fmt.Sprintf("%d.log", time.Now().Unix())
	}

	lumberjackLogger := &lumberjack.Logger{
		// Log file absolute path, os agnostic
		Filename:   filepath.ToSlash(cwd + "/logs/" + logName + ".log"),
		MaxSize:    1,     // MB
		MaxBackups: 10,    // max count of files
		MaxAge:     30,    // days count
		Compress:   false, // disabled by default
	}

	// Fork writing into two outputs
	multiWriter := io.MultiWriter(lumberjackLogger) // io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(log.TextFormatter) // new(log.JSONFormatter)
	// logFormatter := new(log.JSONFormatter)

	logFormatter.TimestampFormat = time.RFC3339 // or RFC1123Z
	logFormatter.FullTimestamp = true

	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel // 4
	}

	logger.SetLevel(logLevel)         // log.SetLevel(log.WarnLevel)
	logger.SetFormatter(logFormatter) // log.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(multiWriter)     // log.SetOutput(os.Stdout)

	ConnectorLogger = logger
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func init() {
	logName = "connector"
	SetupLogger()
}

func LoggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// List of endpoints that doesn't log
		notLoggedRoutes := []string{
			"/connector/api/metrics",
		}

		// Current request path
		requestPath := r.URL.Path

		requestSource := GetIP(r)

		// Serve requests that do not require authentication
		for _, value := range notLoggedRoutes {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		start := time.Now()

		lrw := NewLoggingResponseWriter(w)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(lrw, r)

		fields := log.Fields{
			"method":     r.Method,
			"source":     requestSource,
			"user-agent": r.UserAgent(),
			"uri":        r.RequestURI, // r.URL.Path, // r.RequestURI,
			"payload":    "",
			"status":     lrw.statusCode,
			"latency":    time.Since(start).Milliseconds(),
		}

		// Do stuff here
		ConnectorLogger.WithFields(fields).Info()
	})
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
