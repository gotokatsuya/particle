package particle

import (
	"math/rand"
	"time"
)

// Particle
// X 状態（len = Dimension）
// Weight 重み
type Particle struct {
	X      []int
	Weight float64
}

func NewParticle(dimension int) *Particle {
	weight := 0.0
	x := make([]int, dimension)
	return &Particle{x, weight}
}

// ParticleFilter
// Number 粒子数
// Dimension 次元数
// Upper 最大値（len = Dimension）
// Lower 最大値（len = Dimension）
// Noise 最大値（len = Dimension）
// Particles 粒子の構造体
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

// InitialParticles
// 粒子を一様に分布する
func (f *ParticleFilter) InitialParticles() {
	// 乱数シード
	rand.Seed(time.Now().UnixNano())

	// 粒子の配列生成
	f.Particles = make([]Particle, f.Number)

	// 分布
	for i := 0; i < f.Number; i++ {
		particle := NewParticle(f.Dimension)
		for j := 0; j < f.Dimension; j++ {
			particle.X[j] = (rand.Int() % (f.Upper[j] - f.Lower[j])) + f.Lower[j]
		}
		particle.Weight = 1.0 / float64(f.Number)
		f.Particles[i] = *particle
	}
}

// Resample
// 前状態における重みに基き、粒子を選び直す。（ルーレット選択）
func (f *ParticleFilter) Resample() {
	// 累積重み
	var ws = make([]float64, f.Number)
	ws[0] = f.Particles[0].Weight
	for i := 1; i < f.Number; i++ {
		ws[i] = ws[i-1] + f.Particles[i].Weight
	}

	// 一時的な変数に前状態の粒子を入れる
	var temp = make([]Particle, f.Number)
	for i := 0; i < f.Number; i++ {
		temp[i] = f.Particles[i]
	}

	// 粒子を選び直す
	for i := 0; i < f.Number; i++ {
		var (
			r = float64((rand.Int() % 10000)) / 10000.0
			m = 0
		)
		for j := 0; j < f.Number; j++ {
			if ws[j] >= r {
				m = j
				break
			}
		}
		for k := 0; k < f.Dimension; k++ {
			f.Particles[i].X[k] = temp[m].X[k]
		}
		f.Particles[i].Weight = 0.0
	}
}

// Predict
// 予測モデル（calculate）に従って粒子の位置を予測する
func (f *ParticleFilter) Predict(calculate func(int, []int) Particle) {
	for i := 0; i < f.Number; i++ {
		// ノイズ
		var noises = make([]int, f.Dimension)
		for j := 0; j < f.Dimension; j++ {
			noises[j] = (rand.Int() % (f.Noise[j] * 2)) - f.Noise[j]
		}

		// 推測
		f.Particles[i] = calculate(i, noises)

		// 定められた最大値と最小値を考慮する
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

// Weight
// 尤度計算式（calculate）に従って重み付けする
func (f *ParticleFilter) Weight(calculate func(int) Particle) {
	// 粒子の尤度計算
	for i := 0; i < f.Number; i++ {
		f.Particles[i] = calculate(i)
	}

	// 正規化
	var sum = 0.0
	for i := 0; i < f.Number; i++ {
		sum += f.Particles[i].Weight
	}
	for i := 0; i < f.Number; i++ {
		f.Particles[i].Weight /= sum
	}
}

// Measure
// 重み付き平均から、現状態の粒子を推定する
func (f *ParticleFilter) Measure() Particle {

	var (
		res = NewParticle(f.Dimension)

		x = make([]float64, f.Dimension)
	)

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
