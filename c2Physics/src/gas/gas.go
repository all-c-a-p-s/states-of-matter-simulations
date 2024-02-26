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

func genParticles(n int, temp float64) []Particle {
	// where n is number of particles to generate
	k := temp / 10
	particles := []Particle{}
	for i := 0; i < n; i++ {
		cX := rand.Float64() * 1000.0
		cY := rand.Float64() * 1000.0
		c := Coordinate{cX, cY}
		vX := rand.Float64()*k - k/2
		vY := rand.Float64()*k - k/2
		v := Velocity{vX, vY}
		newParticle := Particle{c, v}
		particles = append(particles, newParticle)
	}
	return particles
}

func distance(p1, p2 Particle) float64 {
	deltaX := p1.coords.x - p2.coords.x
	deltaY := p1.coords.y - p2.coords.y
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

func colliding(p1, p2 Particle, radius float64) bool {
	return distance(p1, p2) < radius*2
}

func updateParticles(particles []Particle) []Particle {
	newParticles := particles
	for i := range newParticles {
		if newParticles[i].coords.y-20.0+newParticles[i].velocity.y < 0 {
			newParticles[i].coords.y = 20
			newParticles[i].velocity.y *= -1
		} else if newParticles[i].coords.y+20.0+newParticles[i].velocity.y > 1000 {
			newParticles[i].coords.y = 980
			newParticles[i].velocity.y *= -1
		} else if newParticles[i].coords.x-20.0+newParticles[i].velocity.x < 0 {
			newParticles[i].coords.x = 20
			newParticles[i].velocity.x *= -1
		} else if newParticles[i].coords.x+20.0+newParticles[i].velocity.x > 1000 {
			newParticles[i].coords.x = 980
			newParticles[i].velocity.x *= -1
		}
		newParticles[i].coords.x += newParticles[i].velocity.x
		newParticles[i].coords.y += newParticles[i].velocity.y
		for j := range newParticles {
			if j == i {
				continue
			}
			if colliding(newParticles[i], newParticles[j], 20.0) {
				tempX := newParticles[i].velocity.x
				tempY := newParticles[i].velocity.y

				newParticles[i].velocity.x = newParticles[j].velocity.x
				newParticles[j].velocity.x = tempX

				newParticles[i].velocity.y = newParticles[j].velocity.y
				newParticles[j].velocity.y = tempY

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
	}
	return newParticles
}

func main() {
	iterations := 10_000

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

	particles := genParticles(20, 600.0)
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

		particles = updateParticles(particles)
	}
	file, err := os.OpenFile("../../images/gas.gif", os.O_WRONLY|os.O_CREATE, 0600)
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
