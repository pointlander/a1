// Copyright 2025 The A1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const (
	// Bit is the number of bits
	Bits = 20
	// Max is the max number to factor
	Max = 1 << Bits
)

func main() {
	process := func(Bits, Max uint64) float64 {
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

		corr := func(numbers []uint64) [][]float64 {
			avg := make([]float64, Bits)
			for _, number := range numbers {
				for i := range avg {
					avg[i] += float64((number >> i) & 1)
				}
			}
			for i := range avg {
				avg[i] /= float64(len(numbers))
			}
			stddev := make([]float64, Bits)
			for _, number := range numbers {
				for i := range stddev {
					diff := avg[i] - float64((number>>i)&1)
					stddev[i] += diff * diff
				}
			}
			for i := range stddev {
				stddev[i] = math.Sqrt(stddev[i] / float64(len(numbers)))
			}
			cov := make([][]float64, Bits)
			for i := range cov {
				cov[i] = make([]float64, Bits)
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
			return cov
		}
		a := corr(primes)
		b := corr(composite)
		sum, count := 0.0, 0.0
		for i := range a {
			for ii := range a[i] {
				scale := a[i][ii] / b[i][ii]
				sum += scale
				count++
				fmt.Printf("%f ", scale)
			}
			fmt.Println()
		}
		fmt.Println()
		return sum / count
	}

	points := make(plotter.XYs, 0, 8)
	for bits := uint64(8); bits < Bits; bits++ {
		gain := process(bits, 1<<bits)
		points = append(points, plotter.XY{X: float64(bits), Y: float64(gain)})
	}
	p := plot.New()

	p.Title.Text = "bits vs gain"
	p.X.Label.Text = "bits"
	p.Y.Label.Text = "gain"

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Radius = vg.Length(1)
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Add(scatter)

	err = p.Save(8*vg.Inch, 8*vg.Inch, "gain.png")
	if err != nil {
		panic(err)
	}
}
