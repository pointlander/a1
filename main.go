// Copyright 2025 The A1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
)

// Max is the max number to factor
const Max = 1024

func main() {
	primes := []uint64{2, 3}
Search:
	for i := uint64(4); i < Max; i++ {
		max := uint64(math.Sqrt(float64(i)) + 1)
		for _, prime := range primes {
			if prime > max {
				break
			}
			if i%prime == 0 {
				continue Search
			}
		}
		primes = append(primes, i)
	}
	fmt.Println(primes)
}
