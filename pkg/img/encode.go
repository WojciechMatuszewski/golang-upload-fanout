package img

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
)

type Encoder struct{}

// EncodeToBuffer encodes image to a buffer
func (e *Encoder) EncodeToBuffer(img image.Image, ct string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	switch ct {
	case "image/png":
		fmt.Println("resizing png")
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, err
		}
	case "image/jpeg":
		fmt.Println("resizing jpeg")
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("")
	}

	return &buf, nil
}

// NewEncoder returns encoder which encodes images to io.Reader
func NewEncoder() *Encoder {
	return &Encoder{}
}
