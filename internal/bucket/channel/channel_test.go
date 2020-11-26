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

package channel_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tamasd/ratelimiter/internal/bucket/channel"
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
	channelBucket := channel.New(mb)
	go channelBucket.Start()
	defer channelBucket.Stop()

	channelBucket.Leak()
	// Make sure that the other goroutine has a chance to run.
	<-time.After(time.Millisecond * 10)

	mb.AssertCalled(t, "Leak")
}

func TestBucket_Input_True(t *testing.T) {
	mb := new(mockBucket)
	mb.On("Input").Return(true)
	channelBucket := channel.New(mb)
	go channelBucket.Start()
	defer channelBucket.Stop()

	require.True(t, channelBucket.Input())
	mb.AssertCalled(t, "Input")
	mb.AssertExpectations(t)
}

func TestBucket_Input_False(t *testing.T) {
	mb := new(mockBucket)
	mb.On("Input").Return(false)
	channelBucket := channel.New(mb)
	go channelBucket.Start()
	defer channelBucket.Stop()

	require.False(t, channelBucket.Input())
	mb.AssertCalled(t, "Input")
	mb.AssertExpectations(t)
}
