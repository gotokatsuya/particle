package particle

import (
	"math/rand"
)

type Particle struct {
	X      []int
	Weight float64
}

func NewParticle(dimension int) *Particle {
	weight := 0.0
	x := make([]int, dimension)
	for i := 0; i < dimension; i++ {
		x[i] = 0
	}
	return &Particle{x, weight}
}

type ParticleFilter struct {
	Number    int
	Dimension int
	Upper     []int
	Lower     []int
	Noise     []int

	Particles []Particle
}

func NewParticleFilter(number, dimension int, upper, lower, noise []int) *ParticleFilter {
	f := ParticleFilter{
		Number:    number,
		Dimension: dimension,
		Upper:     upper,
		Lower:     lower,
		Noise:     noise,
	}
	f.InitialParticles()
	return &f
}

func (f *ParticleFilter) InitialParticles() {
	particles := make([]Particle, f.Number)

	for i := 0; i < f.Number; i++ {
		particle := NewParticle(f.Dimension)
		for j := 0; j < f.Dimension; j++ {
			particle.X[j] = (rand.Int() % (f.Upper[j] - f.Lower[j])) + f.Lower[j]
		}
		particle.Weight = 1.0 / float64(f.Number)
		particles[i] = *particle
	}

	f.Particles = particles
}

func (f *ParticleFilter) Resample() {
	w := make([]float64, f.Number)
	w[0] = f.Particles[0].Weight
	for i := 1; i < f.Number; i++ {
		w[i] = w[i-1] + f.Particles[i].Weight
	}

	temp := make([]Particle, f.Number)
	for i := 0; i < f.Number; i++ {
		temp[i] = f.Particles[i]
	}

	for i := 0; i < f.Number; i++ {
		r := float64((rand.Int() % 10000)) / 10000.0
		for j := 0; j < f.Number; j++ {
			if r > w[j] {
				continue
			}
			for k := 0; k < f.Dimension; k++ {
				f.Particles[i].X[k] = temp[j].X[k]
			}
			f.Particles[i].Weight = 0.0
			break
		}
	}
}

func (f *ParticleFilter) Predict(calculate func(int, []int) Particle) {
	for i := 0; i < f.Number; i++ {
		noises := make([]int, f.Dimension)

		for j := 0; j < f.Dimension; j++ {
			noises[j] = (rand.Int() % (f.Noise[j] * 2)) - f.Noise[j]
		}

		f.Particles[i] = calculate(i, noises)

		for j := 0; j < f.Dimension; j++ {
			if f.Particles[i].X[j] < f.Lower[j] {
				f.Particles[i].X[j] = f.Lower[j]
			}

			if f.Particles[i].X[j] >= f.Upper[j] {
				f.Particles[i].X[j] = f.Upper[j] - 1
			}
		}
	}
}

func (f *ParticleFilter) Weight(calculate func(int) Particle) {
	for i := 0; i < f.Number; i++ {
		f.Particles[i] = calculate(i)
	}

	sum := 0.0
	for i := 0; i < f.Number; i++ {
		sum += f.Particles[i].Weight
	}
	for i := 0; i < f.Number; i++ {
		f.Particles[i].Weight /= sum
	}
}

func (f *ParticleFilter) Measure() Particle {

	res := NewParticle(f.Dimension)

	x := make([]float64, f.Dimension)
	for i := 0; i < f.Dimension; i++ {
		x[i] = 0.0
	}

	for i := 0; i < f.Number; i++ {
		for j := 0; j < f.Dimension; j++ {
			x[j] += float64(f.Particles[i].X[j]) * f.Particles[i].Weight
		}
	}

	for i := 0; i < f.Dimension; i++ {
		res.X[i] = int(x[i])
	}
	return *res
}
