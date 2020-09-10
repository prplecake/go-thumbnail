package thumbnail

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

var scalerTests = []struct {
	scaler string
}{
	{"NearestNeighbor"},
	{"ApproxBiLinear"},
	{"BiLinear"},
	{"CatmullRom"},
}

// TestScalers tests different scalers.
func TestScalers(t *testing.T) {
	config := Generator{
		DestinationPath:   "",
		DestinationPrefix: "thumb_",
	}
	for _, tt := range scalerTests {
		t.Run(tt.scaler, func(t *testing.T) {
			config.Scaler = tt.scaler
			t.Log(config)
			gen := NewGenerator(config)

			i, err := gen.NewImage(testJpegImagePath)
			if err != nil {
				t.Error(err)
			}

			teardownTestCase := setupTestCase(t)
			dest := testDataPath + gen.DestinationPrefix + filepath.Base(i.Path)
			defer teardownTestCase(t, dest)

			thumbBytes, err := gen.Create(i)
			if err != nil {
				t.Error(err)
			}

			err = ioutil.WriteFile(dest, thumbBytes, 0644)
			if err != nil {
				t.Error(err)
			}

			checkFileExists(t, dest)
			var (
				wantWidth  = gen.Width
				wantHeight = gen.Height
			)
			gotWidth, gotHeight, err := checkImageDimensions(t, dest)
			if err != nil {
				t.Error(err)
			}
			if wantWidth != gotWidth {
				t.Errorf("checkImageDimensions() got %d, wants %d", gotWidth, wantWidth)
			}
			if wantHeight != gotHeight {
				t.Errorf("checkImageDimensions() got %d, wants %d", gotHeight, wantHeight)
			}
		})
	}
}

func TestInvalidScaler(t *testing.T) {
	config := Generator{
		DestinationPath:   "",
		DestinationPrefix: "thumb_",
		Scaler:            "",
	}

	gen := NewGenerator(config)
	// Can't use NewImage to create an image since we need to
	// bypass detectContentType

	i, err := gen.NewImage(testJpegImagePath)
	if err != nil {
		t.Error(err)
	}

	errWant := ErrInvalidScaler
	_, err = gen.Create(i)
	if err != nil {
		if err != errWant {
			t.Errorf("Got unexpected error. Expected %s, got %s", errWant, err)
		}
	}
}
