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

package ratelimiter_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tamasd/ratelimiter"
)

func TestMiddleware(t *testing.T) {
	mw := ratelimiter.New(ratelimiter.CreateMiddlewareConfig(1))
	w := httptest.NewRecorder()

	called := false
	mw.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	require.True(t, called)

	called = false
	mw.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	resp := w.Result()
	require.False(t, called)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}

func TestChannelMiddleware(t *testing.T) {
	mw := ratelimiter.NewChannel(ratelimiter.CreateMiddlewareConfig(10))
	go mw.Start()
	t.Cleanup(func() {
		mw.Stop()
	})

	for i := 0; i < 10; i++ {
		require.Equal(t, http.StatusOK, testResponseCode(mw))
	}
	require.Equal(t, http.StatusTooManyRequests, testResponseCode(mw))

	// Make sure that enough time passes for a tick, even if the CPU is busy.
	<-time.After(time.Second / 5)
	require.Equal(t, http.StatusOK, testResponseCode(mw))
}

func testResponseCode(mw *ratelimiter.Middleware) int {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return w.Result().StatusCode
}

func TestDelay(t *testing.T) {
	config := ratelimiter.CreateMiddlewareConfig(0)
	(&config).SetRetryDelay(1000, 10)
	mw := ratelimiter.NewMutex(config)

	require.InDelta(t, 1000, testRetryAfterHeader(mw), 10)
}

func testRetryAfterHeader(mw *ratelimiter.Middleware) int {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	code, err := strconv.Atoi(w.Result().Header.Get("Retry-After"))
	if err != nil {
		panic(err)
	}
	return code
}
