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
	composite := []uint64{}
Search:
	for i := uint64(4); i < Max; i++ {
		max := uint64(math.Sqrt(float64(i)) + 1)
		for _, prime := range primes {
			if prime > max {
				break
			}
			if i%prime == 0 {
				composite = append(composite, i)
				continue Search
			}
		}
		primes = append(primes, i)
	}
	fmt.Println(primes)

	corr := func(numbers []uint64) {
		avg := make([]float64, 10)
		for _, number := range numbers {
			for i := range avg {
				avg[i] += float64((number >> i) & 1)
			}
		}
		for i := range avg {
			avg[i] /= float64(len(numbers))
		}
		stddev := make([]float64, 10)
		for _, number := range numbers {
			for i := range stddev {
				diff := avg[i] - float64((number>>i)&1)
				stddev[i] += diff * diff
			}
		}
		for i := range stddev {
			stddev[i] = math.Sqrt(stddev[i] / float64(len(numbers)))
		}
		cov := make([][]float64, 10)
		for i := range cov {
			cov[i] = make([]float64, 10)
		}
		for _, number := range numbers {
			for i := range cov {
				for ii := range cov[i] {
					a := float64((number>>i)&1) - avg[i]
					b := float64((number>>ii)&1) - avg[ii]
					cov[i][ii] = a + b
				}
			}
		}
		for i := range cov {
			for ii := range cov[i] {
				cov[i][ii] /= float64(len(numbers))
				cov[i][ii] /= (stddev[i] * stddev[ii])
			}
		}
		for _, row := range cov {
			fmt.Println(row)
		}
	}
	corr(primes)
	corr(composite)
}
