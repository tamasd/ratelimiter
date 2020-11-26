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

package bucket

// Bucket interface representing the basic operations of the leaky bucket algo.
type Bucket interface {

	// Input "fills" the bucket.
	//
	// This function will return true, if the bucket can accept the request,
	// and it will return false, if the bucket is overflowing.
	Input() bool

	// Leak lowers the "water level" inside the bucket.
	Leak()
}
