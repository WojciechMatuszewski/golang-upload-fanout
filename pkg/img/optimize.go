package img

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"mime"
	"path/filepath"

	"github.com/disintegration/imaging"
)

type Optimizer struct{}

func NewOptimizer() Optimizer {
	return Optimizer{}
}

func (o Optimizer) Optimize(img image.Image, fn string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	ext := mime.TypeByExtension(filepath.Ext(fn))
	if ext == "" {
		return nil, errors.New("unknown extension")
	}

	ff, err := imaging.FormatFromFilename(fn)
	if err != nil {
		return nil, err
	}

	var encErr error
	switch ext {
	case "image/png":
		encErr = imaging.Encode(&buf, img, ff, imaging.PNGCompressionLevel(png.CompressionLevel(30)))
	case "image/jpeg":
		encErr = imaging.Encode(&buf, img, ff, imaging.JPEGQuality(30))
	default:
		encErr = errors.New(fmt.Sprintf("could not encode unknown format of %v", ff.String()))
	}

	return &buf, encErr

}
