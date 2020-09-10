package thumbnail

import (
	"testing"
)

var mimetypeTests = []struct {
	mimetype string
}{
	{"application/octet-stream"},
}

// TestMimeType tests different mimetypes
func TestMimeType(t *testing.T) {
	gen := NewGenerator()
	gen.Scaler = "CatmullRom"
	for _, tt := range mimetypeTests {
		t.Run(tt.mimetype, func(t *testing.T) {
			t.Log(tt.mimetype)
			// Can't use NewImage to create an image since we need to
			// bypass detectContentType
			image := &Image{
				ContentType: tt.mimetype,
			}
			errWant := ErrInvalidMimeType
			_, err := gen.Create(image)
			if err != nil {
				if err != errWant {
					t.Errorf("Got unexpected error. Expected %s, got %s", errWant, err)
				}
			}
		})
	}
}