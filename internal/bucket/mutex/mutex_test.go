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

package mutex_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tamasd/ratelimiter/internal/bucket/mutex"
)

type mockBucket struct {
	mock.Mock
}

func (mb *mockBucket) Input() bool {
	args := mb.Called()
	return args.Bool(0)
}

func (mb *mockBucket) Leak() {
	mb.Called()
}

func TestBucket_Leak(t *testing.T) {
	mb := new(mockBucket)
	mb.On("Leak").Return()
	mutexBucket := mutex.New(mb)

	mutexBucket.Leak()

	mb.AssertCalled(t, "Leak")
}

func TestBucket_Input_True(t *testing.T) {
	mb := new(mockBucket)
	mb.On("Input").Return(true)
	mutexBucket := mutex.New(mb)

	require.True(t, mutexBucket.Input())
	mb.AssertCalled(t, "Input")
	mb.AssertExpectations(t)
}

func TestBucket_Input_False(t *testing.T) {
	mb := new(mockBucket)
	mb.On("Input").Return(false)
	mutexBucket := mutex.New(mb)

	require.False(t, mutexBucket.Input())
	mb.AssertCalled(t, "Input")
	mb.AssertExpectations(t)
}
