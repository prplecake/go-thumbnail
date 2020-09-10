package thumbnail

import (
	"bytes"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	testDataPath      = "./test_data/"
	testJpegImagePath = testDataPath + "test_image.jpg"
	testPngImagePath  = testDataPath + "test_image.png"
)

var thumbTests = []struct {
	mimeType string
}{
	{"image/jpeg"},
	{"image/png"},
}

// TestThumbTests tests thumbTests
func TestThumbTests(t *testing.T) {
	config := Generator{
		DestinationPath:   "",
		DestinationPrefix: "thumb_",
		Scaler:            "CatmullRom",
	}
	var testImagePath string
	for _, tt := range thumbTests {
		t.Run(tt.mimeType, func(t *testing.T) {
			switch mimeType := tt.mimeType; mimeType {
			case "image/jpeg":
				testImagePath = testJpegImagePath
			case "image/png":
				testImagePath = testPngImagePath
			}
			gen := NewGenerator(config)

			i, err := gen.NewImage(testImagePath)
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

func setupTestCase(t *testing.T) func(t *testing.T, path string) {
	t.Log("Setting up test case.")
	return func(t *testing.T, path string) {
		t.Log("Tearing down test case.")
		err := os.Remove(path)
		if err != nil {
			t.Errorf("Error tearing down test case: %q", err)
		}
	}
}

func checkFileExists(t *testing.T, path string) {
	_, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
			return
		}
		t.Errorf("Error running os.Lstat(%q): %q", path, err)
		return
	}
	t.Log("File exists.")
	return
}

func checkImageDimensions(t *testing.T, path string) (int, int, error) {
	imageBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, 0, err
	}
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return 0, 0, err
	}
	var (
		width  = img.Bounds().Max.X
		height = img.Bounds().Max.Y
	)
	return width, height, nil
}

func Example() {
	var config = Generator{
		DestinationPath:   "",
		DestinationPrefix: "thumb_",
		Scaler:            "CatmullRom",
	}

	imagePath := "path/to/image.jpg"
	dest := "path/to/thumb_image.jpg"
	gen := NewGenerator(config)

	i, err := gen.NewImage(imagePath)
	if err != nil {
		panic(err)
	}

	thumbBytes, err := gen.Create(i)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(dest, thumbBytes, 0644)
	if err != nil {
		panic(err)
	}
}
