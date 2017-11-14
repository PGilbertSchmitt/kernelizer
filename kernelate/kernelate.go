package kernelate

import (
	"image"
	"image/color"
	"math"
)

// Kernel holds the multidemensional matrix of pixel multipliers
// A kernel is a square matrix, generally with an odd width and height
// This means that it should have a center
type Kernel struct {
	K [][]int
	// MaxVal can be overriden
	MaxVal uint32
}

func (k *Kernel) max() uint32 {
	if k.MaxVal > 0 {
		return k.MaxVal
	}

	var sum uint32
	for _, row := range k.K {
		for _, val := range row {
			sum += uint32(val)
		}
	}

	return sum
}

// apply expects an image with the same dimensions as the kernel
func (k *Kernel) apply(img *image.RGBA) color.RGBA {
	var c color.RGBA

	var bigR, bigG, bigB uint32
	b := img.Bounds()

	for x := 0; x < b.Dx(); x++ {
		for y := 0; y < b.Dy(); y++ {
			pixelX := x + b.Min.X
			pixelY := y + b.Min.Y

			// x and y are the kernel coordinates
			// pixelX and pixelY are the image coordinates
			origPixel := img.At(pixelX, pixelY)
			multiplier := uint32(k.K[x][y])

			r32, g32, b32, _ := origPixel.RGBA()
			r := intSqrt(r32)
			g := intSqrt(g32)
			b := intSqrt(b32)

			bigR += r * multiplier
			bigG += g * multiplier
			bigB += b * multiplier
		}
	}

	divisor := k.max()
	if divisor < 0 {
		return c
	}

	bigR /= divisor
	bigG /= divisor
	bigB /= divisor

	c.R = uint8(bigR)
	c.G = uint8(bigG)
	c.B = uint8(bigB)
	c.A = 255

	return c
}

func (k *Kernel) rect(p image.Point) image.Rectangle {
	dist := (len(k.K) - 1) / 2
	pad := image.Point{1, 1}

	offset := image.Point{X: dist, Y: dist}

	return image.Rectangle{Min: p.Sub(offset), Max: p.Add(offset.Add(pad))}
}

// Kernelate passes the center of a kernel over every pixel in an image,
// which causes neighboring pixels to fall under other portions of the
// kernel. Every pixel under the kernel is multiplied by its respective
// kernel multiplier, then averaged with the other pixels by the sum of
// multipliers. Returns a new image.
func Kernelate(img *image.RGBA, k Kernel) (*image.RGBA, error) {
	newImg := image.NewRGBA(img.Bounds())

	// p := image.Point{X: 90, Y: 50}

	// TODO: Seperate image into chunks to be processed in parallel
	// For every point in the image
	for x := 1; x < img.Rect.Dx()-1; x++ {
		for y := 1; y < img.Rect.Dy()-1; y++ {
			p := image.Point{X: x, Y: y}
			subImg := img.SubImage(k.rect(p))
			color := k.apply(subImg.(*image.RGBA))
			newImg.SetRGBA(p.X, p.Y, color)
		}
	}

	// fmt.Println("There are", runtime.NumCPU(), "CPUs in this runtime")

	return newImg, nil
}

// Just needs to return a 16 bit int which holds a max val of 0xFF
func intSqrt(num uint32) uint32 {
	return uint32(math.Sqrt(float64(num)))
}
