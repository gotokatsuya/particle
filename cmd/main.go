package main

import (
	"fmt"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"

	"github.com/gotokatsuya/particle"
)

type point struct {
	x int
	y int
}

func makeField(w, h, mx, my int) (map[point]int, []point) {
	field := map[point]int{}
	targetPoints := []point{}

	length := 10

	maxX := length + mx
	maxY := length + my

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if mx < 0 || maxX > w {
				continue
			}
			if my < 0 || maxY > h {
				continue
			}
			if x >= mx && x <= maxX && y >= my && y <= maxY {
				field[point{x, y}] = 1
				targetPoints = append(targetPoints, point{x, y})
			} else {
				field[point{x, y}] = 0
			}
		}
	}
	return field, targetPoints
}

func drawPlot(title string, ps []point) {
	p, err := plot.New()
	if err != nil {
		return
	}
	p.Title.Text = title
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	pts := make(plotter.XYs, len(ps))
	for i, pp := range ps {
		pts[i].X = float64(pp.x)
		pts[i].Y = float64(pp.y)
	}
	plotutil.AddLinePoints(p, pts)
	width := 4.0 * vg.Inch
	height := 4.0 * vg.Inch
	p.Save(width, height, fmt.Sprintf("%s.png", title))
}

func main() {
	// field
	var (
		width  = 400
		height = 400

		src, currentTargetPoints = makeField(width, height, 0, 0)
	)

	// target
	var (
		mx = 1
		my = 1

		targets = currentTargetPoints
	)

	// particle
	var (
		particlePoints = []point{}

		number    = 800
		dimension = 4

		upper = []int{width, height, 10, 10}
		lower = []int{0, 0, -10, -10}
		noise = []int{30, 30, 10, 10}
	)

	f := particle.NewParticleFilter(number, dimension, upper, lower, noise)

	for i := 0; i < 380; i++ {

		f.Resample()

		f.Predict(func(j int, noises []int) particle.Particle {
			// uniform linear motion
			f.Particles[j].X[0] += f.Particles[j].X[2] + noises[0]
			f.Particles[j].X[1] += f.Particles[j].X[3] + noises[1]
			f.Particles[j].X[2] += noises[2]
			f.Particles[j].X[3] += noises[3]
			return f.Particles[j]
		})

		f.Weight(func(j int) particle.Particle {
			x := f.Particles[j].X[0]
			y := f.Particles[j].X[1]
			if v, ok := src[point{x, y}]; ok {
				if v > 0 {
					f.Particles[j].Weight = 1.0
				} else {
					f.Particles[j].Weight = 0.0001
				}
			}
			return f.Particles[j]
		})

		res := f.Measure()
		particlePoints = append(particlePoints, point{res.X[0], res.X[1]})

		// target moves
		mx++
		my++
		src, currentTargetPoints = makeField(width, height, mx, my)
		targets = append(targets, currentTargetPoints...)
	}

	drawPlot("target", targets)
	drawPlot("particles", particlePoints)

	fmt.Println("Bye")
}
