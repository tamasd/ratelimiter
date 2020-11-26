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

/*
	Package ratelimiter

	This package provides a rate limiter middleware using the "Leaky bucket"
	algorithm's "meter" variant.

	The idea is that the bucket "fills" every time a request is received. If the
	bucket would overflow, then the middleware blocks the request by sending a
	a 429 Too Many Requests status with a Retry-After header. The delay here
	contains a random component to make sure that all clients don't come back
	at the same time in case of a peaky load.

	A background process periodically "leaks" the bucket, so new requests can
	come through.
*/
package ratelimiter
