package thumbnail

import (
	"bytes"
	"image"
	"io/ioutil"
	"net/http"
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

			i, err := gen.NewImageFromFile(testImagePath)
			if err != nil {
				t.Error(err)
			}

			teardownTestCase := setupTestCase(t)
			dest := testDataPath + gen.DestinationPrefix + filepath.Base(i.Path)
			defer teardownTestCase(t, dest)

			thumbBytes, err := gen.CreateThumbnail(i)
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
			gotWidth, gotHeight, err := checkImageDimensions(dest)
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

func TestNewImageFromByteArray(t *testing.T) {
	testImageURL := "https://s3-us-west-2.amazonaws.com/freeradical-system/media_attachments/files/107/742/646/709/715/386/original/8d771cb38a5984ab.jpg"
	if testImageURL == "" {
		t.Fatal("No test image URL.")
	}
	resp, err := http.Get(testImageURL)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	config := Generator{
		DestinationPath:   "",
		DestinationPrefix: "thumb_",
		Scaler:            "CatmullRom",
	}
	gen := NewGenerator(config)
	i, err := gen.NewImageFromByteArray(data)
	if err != nil {
		t.Error(err)
	}

	teardownTestCase := setupTestCase(t)
	dest := testDataPath + gen.DestinationPrefix + filepath.Base(testImageURL)
	defer teardownTestCase(t, dest)

	thumbBytes, err := gen.CreateThumbnail(i)
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
	gotWidth, gotHeight, err := checkImageDimensions(dest)
	if err != nil {
		t.Error(err)
	}
	if wantWidth != gotWidth {
		t.Errorf("checkImageDimensions() got %d, wants %d", gotWidth, wantWidth)
	}
	if wantHeight != gotHeight {
		t.Errorf("checkImageDimensions() got %d, wants %d", gotHeight, wantHeight)
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
}

func checkImageDimensions(path string) (int, int, error) {
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

	i, err := gen.NewImageFromFile(imagePath)
	if err != nil {
		panic(err)
	}

	thumbBytes, err := gen.CreateThumbnail(i)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(dest, thumbBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func (*Generator) ExampleNewImageFromByteArray() {
	testImageURL := "https://example.com/image.jpg"
	if testImageURL == "" {
		panic("No test image URL.")
	}
	resp, err := http.Get(testImageURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	config := Generator{
		DestinationPath:   "",
		DestinationPrefix: "thumb_",
		Scaler:            "CatmullRom",
	}
	gen := NewGenerator(config)
	i, err := gen.NewImageFromByteArray(data)
	if err != nil {
		panic(err)
	}

	dest := testDataPath + gen.DestinationPrefix + filepath.Base(testImageURL)

	thumbBytes, err := gen.CreateThumbnail(i)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(dest, thumbBytes, 0644)
	if err != nil {
		panic(err)
	}
}
