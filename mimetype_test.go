package thumbnail

import (
	"errors"
	"testing"
)

var mimetypeTests = []struct {
	mimetype string
}{
	{"application/octet-stream"},
}

// TestMimeType tests different mimetypes
func TestMimeType(t *testing.T) {
	config := Generator{
		DestinationPath:   "",
		DestinationPrefix: "thumb_",
		Scaler:            "NearestNeighbor",
	}
	gen := NewGenerator(config)
	for _, tt := range mimetypeTests {
		t.Run(tt.mimetype, func(t *testing.T) {
			t.Log(tt.mimetype)
			// Can't use NewImage to create an image since we need to
			// bypass detectContentType
			image := &Image{
				ContentType: tt.mimetype,
			}
			errWant := ErrInvalidMimeType
			_, err := gen.CreateThumbnail(image)
			if err != nil {
				if !errors.Is(err, errWant) {
					t.Errorf("Got unexpected error. Expected %s, got %s", errWant, err)
				}
			}
		})
	}
}
