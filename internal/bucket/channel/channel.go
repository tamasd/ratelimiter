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

package channel

import "github.com/tamasd/ratelimiter/internal/bucket"

type inputMessage struct {
	reply chan bool
}

func newInputMessage() inputMessage {
	return inputMessage{
		reply: make(chan bool),
	}
}

// Bucket decorates a bucket with channels.
//
// This allows multiple goroutines to use the decorated bucket.
type Bucket struct {
	inputch chan inputMessage
	leakch  chan struct{}
	quitch  chan struct{}
	bucket  bucket.Bucket
}

// New creates a new channel bucket.
func New(bucket bucket.Bucket) *Bucket {
	return &Bucket{
		inputch: make(chan inputMessage),
		leakch:  make(chan struct{}),
		quitch:  make(chan struct{}),
		bucket:  bucket,
	}
}

// Input calls the decorated bucket's Input().
func (cb *Bucket) Input() bool {
	message := newInputMessage()
	cb.inputch <- message
	return <-message.reply
}

// Leak calls the decorated bucket's Leak().
func (cb *Bucket) Leak() {
	cb.leakch <- struct{}{}
}

// Start starts the service.
func (cb *Bucket) Start() {
	defer func() {
		cb.quitch = make(chan struct{})
	}()
	for {
		select {
		case input := <-cb.inputch:
			input.reply <- cb.bucket.Input()
		case <-cb.leakch:
			cb.bucket.Leak()
		case <-cb.quitch:
			return
		}
	}
}

// Stop stops the service.
func (cb *Bucket) Stop() {
	close(cb.quitch)
}
