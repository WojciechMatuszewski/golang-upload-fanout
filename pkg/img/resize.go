package img

import (
	"image"

	"github.com/disintegration/imaging"
)

type Resizer struct {
}

// ResizeToThumbnail resizes given image to an image of width 128 and height of 128
func (r *Resizer) ResizeToThumbnail(img image.Image) *image.NRGBA {
	return imaging.Resize(img, 128, 128, imaging.NearestNeighbor)
}

// NewResizer creates a Resizer
func NewResizer() *Resizer {
	return &Resizer{}
}
