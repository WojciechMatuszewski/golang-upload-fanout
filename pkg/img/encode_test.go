package img_test

import (
	"image"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"testing-stuff/pkg/img"
)

func TestEncoder_EncodeToReader(t *testing.T) {
	t.Parallel()

	encoder := img.NewEncoder()

	t.Run("returns error on unknown extension", func(t *testing.T) {
		mockImage := image.NewNRGBA(image.Rect(1, 1, 2, 2))
		_, err := encoder.EncodeToBuffer(mockImage, ".FOO")

		assert.Error(t, err)
	})

	t.Run("success for .png", func(t *testing.T) {
		mockImage := image.NewNRGBA(image.Rect(1, 1, 2, 2))
		b, err := encoder.EncodeToBuffer(mockImage, "image/png")

		assert.NoError(t, err)

		_, err = png.Decode(b)

		assert.NoError(t, err)
	})

	t.Run("success for .jpeg", func(t *testing.T) {
		mockImage := image.NewNRGBA(image.Rect(1, 1, 2, 2))
		b, err := encoder.EncodeToBuffer(mockImage, "image/jpeg")

		assert.NoError(t, err)

		_, err = jpeg.Decode(b)

		assert.NoError(t, err)
	})

	t.Run("fails on invalid image proportions", func(t *testing.T) {
		mockImage := image.NewNRGBA(image.Rect(1, 1, 1, 1))
		_, err := encoder.EncodeToBuffer(mockImage, "image/png")

		assert.Error(t, err)
	})
}
