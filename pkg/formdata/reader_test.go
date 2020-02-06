package formdata_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"testing-stuff/pkg/formdata"
)

const contentType string = "multipart/form-data; boundary=--------------------------556307903488247805978216"

func TestReader_ReadBase64Encoded(t *testing.T) {
	t.Parallel()
	reader := formdata.NewReader()

	t.Run("fails when base64 content is malformed", func(t *testing.T) {
		f, err := reader.ReadBase64Encoded("FOO", contentType)

		assert.Nil(t, f)
		assert.Error(t, err, assert.AnError)
	})

	t.Run("succeeds when base64 content is OK", func(t *testing.T) {
		rb, err := ioutil.ReadFile("./testdata/encoded")
		if err != nil {
			t.Fatal(err)
		}

		f, err := reader.ReadBase64Encoded(string(rb), contentType)

		assert.NoError(t, err)
		assert.Equal(t, f.File["image"][0].Filename, "cat-image.png")
		assert.Equal(t, f.File["image"][0].Size, int64(85263))
	})

	t.Run("fails when base64 content is OK but contentType is wrong", func(t *testing.T) {
		rb, err := ioutil.ReadFile("./testdata/encoded")
		if err != nil {
			t.Fatal(err)
		}

		f, err := reader.ReadBase64Encoded(string(rb), "BAR")

		assert.Error(t, err)
		assert.Nil(t, f)
	})
}
