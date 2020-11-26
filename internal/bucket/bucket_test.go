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

package bucket_test

import (
	"testing"

	"github.com/tamasd/ratelimiter/internal/bucket"
	"github.com/tamasd/ratelimiter/internal/bucket/channel"
	"github.com/tamasd/ratelimiter/internal/bucket/leaky"
	"github.com/tamasd/ratelimiter/internal/bucket/mutex"
)

func benchBucket(b *testing.B, bucket bucket.Bucket) {
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bucket.Input()
			bucket.Leak()
		}
	})
}

func BenchmarkMutex(b *testing.B) {
	mutexBucket := mutex.New(leaky.New(1))

	benchBucket(b, mutexBucket)
}

func BenchmarkChannel(b *testing.B) {
	channelBucket := channel.New(leaky.New(1))
	go channelBucket.Start()
	b.Cleanup(func() {
		channelBucket.Stop()
	})

	benchBucket(b, channelBucket)
}
