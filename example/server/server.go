// A simple rate limiter middleware.
// Copyright (c) 2020. Tamás Demeter-Haludka
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tamasd/ratelimiter"
	"github.com/urfave/negroni"
)

var (
	loggerOut      = os.Stdout
	loggerExitFunc = os.Exit

	logLevel          = flag.String("loglevel", "info", "")
	requestsPerSecond = flag.Uint("req", 100, "requests per second")
)

const (
	loggerContextKey = "logger"
)

// GetLogger returns the logger from the request context.
func GetLogger(r *http.Request) logrus.FieldLogger {
	return r.Context().Value(loggerContextKey).(logrus.FieldLogger)
}

// ProfilerMiddleware creates a profiling middleware using a logger.
//
// This middleware also puts the logger into the request context.
func ProfilerMiddleware(logger logrus.FieldLogger) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()

		l := logger.WithFields(logrus.Fields{
			"reqid":  randomHexString(16),
			"method": r.Method,
			"path":   r.URL.Path,
			"host":   r.Host,
		})

		w.Header().Set("Server", "Unknown")
		r = r.WithContext(context.WithValue(r.Context(), loggerContextKey, l))

		next(w, r)

		status := w.(negroni.ResponseWriter).Status()
		l.WithFields(logrus.Fields{
			"status-code": status,
			"status":      http.StatusText(status),
			"latency":     time.Since(start),
		}).Infoln("completed handling request")
	}
}

// CreateHandler creates a http.Handler using negroni.
func CreateHandler(logger logrus.FieldLogger, requestsPerSecond uint) http.Handler {
	middleware := negroni.New()

	recovery := negroni.NewRecovery()
	recovery.Logger = logger
	middleware.Use(recovery)
	middleware.UseFunc(ProfilerMiddleware(logger))

	ratelimiterMiddleware := ratelimiter.New(ratelimiter.CreateMiddlewareConfig(requestsPerSecond))
	middleware.Use(ratelimiterMiddleware)
	go ratelimiterMiddleware.Start()

	middleware.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetLogger(r).Infoln("working on request")
		w.WriteHeader(http.StatusOK)
	})

	return middleware
}

// CreateLogger creates and configures a logger.
func CreateLogger() logrus.FieldLogger {
	logger := logrus.New()

	// This makes the logger testable.
	logger.Out = loggerOut
	logger.ExitFunc = loggerExitFunc

	lvl, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logger.WithError(err).Fatalln("failed to parse log level")
		return nil
	}
	logger.SetLevel(lvl)

	hostname, _ := os.Hostname()
	return logger.WithField("hostname", hostname)
}

// CreateServer creates and configures a http.Server.
func CreateServer(handler http.Handler) *http.Server {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:    host + ":" + port,
		Handler: handler,
	}
}

// SetupServer creates an http.Server.
//
// This function also creates all the dependencies from the global flags.
func SetupServer() *http.Server {
	logger := CreateLogger()
	handler := CreateHandler(logger, *requestsPerSecond)
	server := CreateServer(handler)

	return server
}

// randomHexString creates a random string with a given length.
//
// This string will only contain "hex" characters: numbers and letters a-f.
func randomHexString(length int) string {
	buf := make([]byte, length/2+length%2)
	_, _ = io.ReadFull(rand.Reader, buf)
	return hex.EncodeToString(buf)[:length]
}

func main() {
	flag.Parse()
	if err := SetupServer().ListenAndServe(); err != nil {
		panic(err)
	}
}
