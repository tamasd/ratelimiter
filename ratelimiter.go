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

package ratelimiter

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/tamasd/ratelimiter/internal/bucket"
	"github.com/tamasd/ratelimiter/internal/bucket/channel"
	"github.com/tamasd/ratelimiter/internal/bucket/leaky"
	"github.com/tamasd/ratelimiter/internal/bucket/mutex"
)

type startStop interface {
	Start()
	Stop()
}

// MiddlewareConfig holds the configuration for the Middleware type.
type MiddlewareConfig struct {
	requestPerSecond uint
	retryDelay       uint
	random           uint
}

// CreateMiddlewareConfig creates the configration for the middleware.
//
// The only parameter is requestPerSecond, which tells the middleware how many
// requests per second it can let through.
func CreateMiddlewareConfig(requestPerSecond uint) MiddlewareConfig {
	return MiddlewareConfig{
		requestPerSecond: requestPerSecond,
		retryDelay:       1,
		random:           5,
	}
}

// SetRetryDelay sets the dynamic delay values of the middleware.
//
// When the middleware rejects a request, it sends a Retry-After header to ask
// the client to retry the request a few seconds later.
//
// The retryDelay parameter is the minimum retry delay in seconds. The second,
// random, parameter adds an extra random value between 0 and random to the
// delay. The point of this is to smooth out the load when clients coming back
// in case of a very spikey load.
func (mc *MiddlewareConfig) SetRetryDelay(retryDelay, random uint) {
	mc.retryDelay = retryDelay
	mc.random = random
}

func (mc MiddlewareConfig) delay() uint {
	return mc.retryDelay + uint(rand.Intn(int(mc.random)))
}

// Middleware is the rate limiter middleware.
//
// When using this middleware, make sure that you call Start() before starting
// the http server.
type Middleware struct {
	config MiddlewareConfig
	bucket bucket.Bucket

	leakch chan struct{}
	quitch chan struct{}
}

// New creates a rate limiter middleware with the default method.
//
// This is an alias to NewMutex() function at the moment.
func New(config MiddlewareConfig) *Middleware {
	return NewMutex(config)
}

// NewMutex creates a rate limiter middleware using mutexes.
//
// While not as elegant, mutexes outperform channels (see the README.md file
// for more details).
func NewMutex(config MiddlewareConfig) *Middleware {
	return newMiddleware(config, func(bucket bucket.Bucket) bucket.Bucket {
		return mutex.New(bucket)
	})
}

// NewChannel creates a rate limiter middleware using channels.
func NewChannel(config MiddlewareConfig) *Middleware {
	return newMiddleware(config, func(bucket bucket.Bucket) bucket.Bucket {
		return channel.New(bucket)
	})
}

func newMiddleware(config MiddlewareConfig, bucketFactory func(bucket bucket.Bucket) bucket.Bucket) *Middleware {
	return &Middleware{
		config: config,
		bucket: bucketFactory(leaky.New(config.requestPerSecond)),
		leakch: make(chan struct{}),
		quitch: make(chan struct{}),
	}
}

// ServeHTTP implements negroni.Handler interface.
//
// If the rate limiter blocks the request, 429 Too Many Requests will be
// returned with a Retry-After header, asking the client to retry the request
// later.
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if m.bucket.Input() {
		next.ServeHTTP(w, r)
	} else {
		w.Header().Set("Retry-After", strconv.Itoa(int(m.config.delay())))
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
	}
}

// Start starts the middleware's "leak" logic.
//
// This will "drain" the bucket at the configured rate. Make sure you call this
// before you start the http server.
func (m *Middleware) Start() {
	if ss, ok := m.bucket.(startStop); ok {
		go func() { ss.Start() }()
	}

	for {
		select {
		case <-time.After(time.Second / time.Duration(m.config.requestPerSecond)):
			m.bucket.Leak()
		case <-m.quitch:
			return
		}
	}
}

// Stop stops the middleware's internal loop.
//
// After calling this function the middleware is not usable anymore. Make sure
// you call this after the http server is stopped.
func (m *Middleware) Stop() {
	close(m.quitch)
	if ss, ok := m.bucket.(startStop); ok {
		ss.Stop()
	}
}
