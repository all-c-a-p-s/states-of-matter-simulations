package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math"
	"math/rand"
	"os"

	"github.com/fogleman/gg"
)

const (
	GRAVITY                    float64 = 10
	COEFFICIENT_OF_RESTITUTION float64 = 0.5
	P_ESCAPE                   float64 = 0.01
	MAX_K                      float64 = 2.0
	MAX_SPEED                  float64 = 10.0
)

type Coordinate struct {
	x float64
	y float64
}

type Velocity struct {
	x float64
	y float64
}

type Particle struct {
	coords   Coordinate
	velocity Velocity
}

func calculateSpeed(p Particle) float64 {
	return math.Sqrt(p.velocity.x*p.velocity.x + p.velocity.y*p.velocity.y)
}

func genParticles(n int, temp float64) []Particle {
	// where n is number of particles to generate
	k := temp / 10
	particles := []Particle{}
	centres := startCentres()
	for i := 0; i < n; i++ {
		c := centres[i]
		vX := rand.Float64()*k - k/2
		vY := rand.Float64()*k - k/2
		v := Velocity{vX, vY}
		newParticle := Particle{c, v}
		particles = append(particles, newParticle)
	}
	return particles
}

func startCentres() []Coordinate {
	var centres []Coordinate
	for i := 20; i <= 970; i += 50 {
		for j := 960; j >= 660; j -= 50 {
			c := Coordinate{
				x: float64(i),
				y: float64(j),
			}
			centres = append(centres, c)
		}
	}
	return centres
}

func distance(p1, p2 Particle) float64 {
	deltaX := p1.coords.x - p2.coords.x
	deltaY := p1.coords.y - p2.coords.y
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

func colliding(p1, p2 Particle, radius float64) bool {
	return distance(p1, p2) < radius*2
}

func updateParticles(particles []Particle, genNumber int) []Particle {
	var speedK float64 = float64((500 + genNumber)) / 500
	speedK = min(speedK, MAX_K)
	gravityK := GRAVITY * float64(500/(genNumber+500))
	// gravityK = max(1/MAX_K, gravityK)
	newParticles := particles
	for i := range newParticles {

		if genNumber < 50 {
			escape_p := rand.Float64()
			if escape_p > (1 - P_ESCAPE) {
				newParticles[i].velocity.x *= 2 // escape from cluster of particles
				newParticles[i].velocity.y *= 2 // if the particle is not on the top layer its speed propagates upwards by collisions
			}
		} else if genNumber > 200 {

			reset := newParticles[i]

			newParticles[i].velocity.x *= float64(speedK)
			newParticles[i].velocity.y *= float64(speedK)

			if calculateSpeed(newParticles[i]) > MAX_SPEED {
				newParticles[i] = reset
			}

		}
		k1 := 1
		k2 := 1

		// negative or positive change
		r := rand.Float32()
		if r > 0.5 {
			k1 = -1
		}
		R := rand.Float32()
		if R > 0.5 {
			k2 = -1
		}

		r1 := rand.Float32()
		newParticles[i].velocity.x += float64(r1) * float64(k1) * 2
		r2 := rand.Float32()
		newParticles[i].velocity.y += float64(r2) * float64(k2) * 2

		// random small velocity to mitigate strange behaviour on the side

		if newParticles[i].coords.y >= 900 {
			newParticles[i].velocity.y -= 1
		}

		if newParticles[i].coords.y-20.0+newParticles[i].velocity.y < 0 {
			newParticles[i].coords.y = 20
			newParticles[i].velocity.y *= -COEFFICIENT_OF_RESTITUTION
		} else if newParticles[i].coords.y+20.0+newParticles[i].velocity.y > 1000 {
			newParticles[i].coords.y = 980
			newParticles[i].velocity.y *= -COEFFICIENT_OF_RESTITUTION
		} else if newParticles[i].coords.x-20.0+newParticles[i].velocity.x < 0 {
			newParticles[i].coords.x = 20
			newParticles[i].velocity.x *= -COEFFICIENT_OF_RESTITUTION
		} else if newParticles[i].coords.x+20.0+newParticles[i].velocity.x > 1000 {
			newParticles[i].coords.x = 980
			newParticles[i].velocity.x *= -COEFFICIENT_OF_RESTITUTION
		}

		for j := range newParticles {
			if j == i {
				continue
			}
			if colliding(newParticles[i], newParticles[j], 20.0) {
				tempX := newParticles[i].velocity.x
				tempY := newParticles[i].velocity.y

				newParticles[i].velocity.x = newParticles[j].velocity.x * COEFFICIENT_OF_RESTITUTION
				newParticles[j].velocity.x = tempX * COEFFICIENT_OF_RESTITUTION

				newParticles[i].velocity.y = newParticles[j].velocity.y * COEFFICIENT_OF_RESTITUTION
				newParticles[j].velocity.y = tempY * COEFFICIENT_OF_RESTITUTION

				distanceX := newParticles[j].coords.x - newParticles[i].coords.x
				distanceY := newParticles[j].coords.y - newParticles[i].coords.y

				max_attempts := 100
				attempts := 0

				for colliding(newParticles[i], newParticles[j], 20.0) { // update position until no longer collliding
					attempts++
					newParticles[j].coords.x += distanceX / 10
					newParticles[j].coords.y += distanceY / 10
					if attempts > max_attempts {
						break
					}
				}
			}
		}

		if newParticles[i].coords.y >= 960 { // close to the bottom
			// done to mitigate weird behaviour at the bottom
			newParticles[i].velocity.y *= 0.1
			newParticles[i].velocity.x *= 0.1
		}
	}
	for i := range newParticles {
		newParticles[i].velocity.y += gravityK
	}
	for i := range newParticles {
		newParticles[i].coords.x += newParticles[i].velocity.x
		newParticles[i].coords.y += newParticles[i].velocity.y
	}
	return newParticles
}

func main() {
	iterations := 1000

	var images []*image.Paletted
	var delays []int
	var disposals []byte

	palette := color.Palette{
		image.Transparent,
		image.Black,
		image.White,
		color.RGBA{0xEE, 0xEE, 0xEE, 255},
		color.RGBA{0xCC, 0xCC, 0xCC, 255},
		color.RGBA{0x99, 0x99, 0x99, 255},
		color.RGBA{0x66, 0x66, 0x66, 255},
		color.RGBA{0x33, 0x33, 0x33, 255},
	}

	particles := genParticles(133, 100.0)
	for i := 0; i < iterations; i++ {
		dc := gg.NewContext(250.0, 250.0)
		dc.SetRGBA(1, 1, 1, 0)
		dc.Clear()
		for _, c := range particles {
			dc.DrawCircle(c.coords.x/4.0, c.coords.y/4.0, 5.0)
		}

		dc.SetRGBA(0, 0, 0, 1)
		dc.Fill()

		img := dc.Image()
		bounds := img.Bounds()
		dst := image.NewPaletted(bounds, palette)
		draw.Draw(dst, bounds, img, bounds.Min, draw.Src)
		images = append(images, dst)
		delays = append(delays, 1)
		disposals = append(disposals, gif.DisposalBackground)

		particles = updateParticles(particles, i)
	}
	file, err := os.OpenFile("../../images/evaporation.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic("error creating file")
	}
	defer file.Close()
	gif.EncodeAll(file, &gif.GIF{
		Image:    images,
		Delay:    delays,
		Disposal: disposals,
	})
}
