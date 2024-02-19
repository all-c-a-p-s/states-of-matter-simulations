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

var START_CENTRES []Coordinate = startCentres()

type Coordinate struct {
	x float64
	y float64
}

func distance(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

func startCentres() []Coordinate {
	var centres []Coordinate
	for i := 200; i <= 800; i += 50 {
		for j := 200; j <= 800; j += 50 {
			c := Coordinate{
				x: float64(i),
				y: float64(j),
			}
			centres = append(centres, c)
		}
	}
	return centres
}

func UpdateCentres(centres []Coordinate, temp float64) []Coordinate {
	k := temp / 100.0
	newCentres := centres
	for i := range newCentres {
		newCentres[i].x += rand.Float64()*k - k/2
		newCentres[i].y += rand.Float64()*k - k/2
		resetCount := 0
		for distance(newCentres[i].x-START_CENTRES[i].x, newCentres[i].y-START_CENTRES[i].y) > 5.0 {
			resetCount++
			newCentres[i].x = START_CENTRES[i].x
			newCentres[i].y = START_CENTRES[i].y
			if resetCount >= 3 {
				break
			}
			newCentres[i].x += rand.Float64()*k - k/2
			newCentres[i].y += rand.Float64()*k - k/2

		}
	}
	return newCentres
}

func main() {
	iterations := 10_000
	centres := startCentres()

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

	for i := 0; i < iterations; i++ {
		dc := gg.NewContext(250.0, 250.0)
		dc.SetRGBA(1, 1, 1, 0)
		dc.Clear()
		for _, c := range centres {
			dc.DrawCircle(c.x/4.0, c.y/4.0, 5.0)
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

		centres = UpdateCentres(centres, 1000.0)
	}

	file, err := os.OpenFile("../../images/solid.gif", os.O_WRONLY|os.O_CREATE, 0600)
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
