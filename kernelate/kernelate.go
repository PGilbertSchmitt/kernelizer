package kernelate

import (
	"image"
)

// Kernel holds the multidemensional matrix of pixel multipliers
// A kernel is a square matrix, generally with an odd width and height
// This means that it should have a center
type Kernel struct {
	K [][]int
}

// Kernelate passes the center of a kernel over every pixel in an image,
// which causes neighboring pixels to fall under other portions of the
// kernel. Every pixel under the kernel is multiplied by its respective
// kernel multiplier, then averaged with the other pixels by the sum of
// multipliers. Returns a new image.
func Kernelate(img *image.RGBA, k Kernel) (*image.RGBA, error) {
	newImg := image.NewRGBA(img.Bounds())
	return newImg, nil
}
