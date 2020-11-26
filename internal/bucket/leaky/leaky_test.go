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

package leaky_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tamasd/ratelimiter/internal/bucket/leaky"
)

func TestBucket_Input_Success(t *testing.T) {
	bucket := leaky.New(1)

	require.True(t, bucket.Input())
}

func TestBucket_Input_Fail(t *testing.T) {
	bucket := leaky.New(1)
	bucket.Input()

	require.False(t, bucket.Input())
}

func TestBucket_Leak(t *testing.T) {
	bucket := leaky.New(1)

	bucket.Input()
	bucket.Leak()

	require.True(t, bucket.Input())
}

func TestBucket_Leak_Multiple(t *testing.T) {
	bucket := leaky.New(1)

	bucket.Input()
	bucket.Leak()
	bucket.Leak()

	require.True(t, bucket.Input())
}

func TestBucket_Leak_OnlyOne(t *testing.T) {
	bucket := leaky.New(2)

	bucket.Input()
	bucket.Input()
	bucket.Leak()

	require.True(t, bucket.Input())
	require.False(t, bucket.Input())
}
