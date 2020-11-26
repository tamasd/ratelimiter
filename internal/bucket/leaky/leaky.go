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

package leaky

// Bucket is the implementation of the leaky bucket algorithm.
//
// This implementation is not thread-safe.
type Bucket struct {
	limit   uint
	counter uint
}

// New creates a new instance of the leaky bucket.
//
// The limit parameter sets the maximum capacity of the bucket.
func New(limit uint) *Bucket {
	return &Bucket{
		limit: limit,
	}
}

// Input tries to increment the internal counter.
//
// If the internal counter is at the limit, this will return false, and don't
// increment the counter further. Otherwise the counter will be incremented and
// true is returned.
func (lb *Bucket) Input() bool {
	if lb.counter < lb.limit {
		lb.counter++
		return true
	}

	return false
}

// Leak decrements the internal counter.
func (lb *Bucket) Leak() {
	if lb.counter > 0 {
		lb.counter--
	}
}
