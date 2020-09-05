package thumbnail

import (
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

var thumbCfg = Configuration{
	DestinationPrefix: testDataPath + "thumb_",
}

// TestJpegThumbnail tests JPEG thumbnail generation.
func TestJpegThumbnail(t *testing.T) {
	thumbCfg.Path = testJpegImagePath
	thumbCfg.ContentType = "image/jpeg"
	teardownTestCase := setupTestCase(t)
	dest := thumbCfg.DestinationPrefix + filepath.Base(thumbCfg.Path)
	defer teardownTestCase(t, dest)

	testImage, err := ioutil.ReadFile(thumbCfg.Path)
	if err != nil {
		t.Error(err)
	}

	err = Create(testImage, thumbCfg)
	if err != nil {
		t.Error(err)
	}

	checkFileExists(t, dest)
}

// TestPngThumbnail tests PNG thumbnail generation.
func TestPngThumbnail(t *testing.T) {
	thumbCfg.Path = testPngImagePath
	thumbCfg.ContentType = "image/png"
	teardownTestCase := setupTestCase(t)
	dest := thumbCfg.DestinationPrefix + filepath.Base(thumbCfg.Path)
	defer teardownTestCase(t, dest)

	testImage, err := ioutil.ReadFile(thumbCfg.Path)
	if err != nil {
		t.Error(err)
	}

	err = Create(testImage, thumbCfg)
	if err != nil {
		t.Error(err)
	}

	checkFileExists(t, dest)
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
