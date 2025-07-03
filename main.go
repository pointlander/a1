// Copyright 2025 The A1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
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

// linearRegression calculates the slope (m) and intercept (b) of a linear regression line
// for a given set of data points using the least squares method.
func linearRegression(data plotter.XYs) (m, b float64) {
	n := float64(len(data))
	if n == 0 {
		return 0, 0 // Handle empty data set
	}

	var sumX, sumY, sumXY, sumXX float64
	for _, p := range data {
		sumX += p.X
		sumY += p.Y
		sumXY += p.X * p.Y
		sumXX += p.X * p.X
	}

	// Calculate slope (m)
	numeratorM := n*sumXY - sumX*sumY
	denominatorM := n*sumXX - sumX*sumX
	if denominatorM == 0 {
		// Handle cases where all X values are the same (vertical line)
		return math.Inf(1), 0 // Slope is infinite
	}
	m = numeratorM / denominatorM

	// Calculate intercept (b)
	b = (sumY - m*sumX) / n

	return m, b
}

var (
	// FlagCorr is the correlation mode
	FlagCorr = flag.Bool("corr", false, "correlation mode")
	// FlagCountBits count the bits
	FlagCountBits = flag.Bool("count", false, "count the number of bits")
)

// Corr is the correlation mode
func Corr() {
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
			fmt.Println(stddev)
			cov := make([][]float64, Bits)
			for i := range cov {
				cov[i] = make([]float64, Bits)
			}
			for _, number := range numbers {
				for i := range cov {
					for ii := range cov[i] {
						a := float64((number>>i)&1) - avg[i]
						b := float64((number>>ii)&1) - avg[ii]
						cov[i][ii] += a * b
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
				if i == 0 || i == 1 || ii == 0 || ii == 1 {
					continue
				}
				count++
				scale := math.Abs(a[i][ii] / b[i][ii])
				if math.IsInf(scale, 0) || math.IsNaN(scale) || scale > 20 {
					fmt.Printf("%f ", scale)
					continue
				}
				sum += scale
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

	m, b := linearRegression(points)
	phi := (1 + math.Sqrt(5)) / 2
	fmt.Println("phi=", phi)
	fmt.Println("1/phi=", 1/phi)
	fmt.Println("m=", m)
	fmt.Println("b=", b)
}

// CountBits count the bits
func CountBits() {
	primes := []uint64{2, 3}
	composites := []uint64{}
Search:
	for i := uint64(4); i < 1<<24; i++ {
		max := uint64(math.Sqrt(float64(i)) + 1)
		for _, prime := range primes {
			if prime > max {
				break
			}
			if i%prime == 0 {
				composites = append(composites, i)
				continue Search
			}
		}
		primes = append(primes, i)
	}

	{
		points := make(plotter.XYs, 0, 8)
		for _, prime := range primes {
			p, count := prime, 0
			for p != 0 {
				if p&1 == 1 {
					count++
				}
				p >>= 1
			}
			points = append(points, plotter.XY{X: float64(prime), Y: float64(count)})
		}

		p := plot.New()

		p.Title.Text = "number vs bits"
		p.X.Label.Text = "number"
		p.Y.Label.Text = "bits"

		scatter, err := plotter.NewScatter(points)
		if err != nil {
			panic(err)
		}
		scatter.GlyphStyle.Radius = vg.Length(1)
		scatter.GlyphStyle.Shape = draw.CircleGlyph{}
		p.Add(scatter)

		err = p.Save(8*vg.Inch, 8*vg.Inch, "primes.png")
		if err != nil {
			panic(err)
		}
	}

	{
		points := make(plotter.XYs, 0, 8)
		for _, composite := range composites {
			c, count := composite, 0
			for c != 0 {
				if c&1 == 1 {
					count++
				}
				c >>= 1
			}
			points = append(points, plotter.XY{X: float64(composite), Y: float64(count)})
		}

		p := plot.New()

		p.Title.Text = "number vs bits"
		p.X.Label.Text = "number"
		p.Y.Label.Text = "bits"

		scatter, err := plotter.NewScatter(points)
		if err != nil {
			panic(err)
		}
		scatter.GlyphStyle.Radius = vg.Length(1)
		scatter.GlyphStyle.Shape = draw.CircleGlyph{}
		p.Add(scatter)

		err = p.Save(8*vg.Inch, 8*vg.Inch, "composites.png")
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	flag.Parse()

	if *FlagCorr {
		Corr()
		return
	}

	if *FlagCountBits {
		CountBits()
		return
	}
}
