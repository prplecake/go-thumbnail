package thumbnail

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	testImageOrientation6Path = testDataPath + "test_image_orientation_6.jpg"
)

// TestEXIFOrientation tests that EXIF orientation is properly applied
func TestEXIFOrientation(t *testing.T) {
	config := Generator{
		DestinationPath:   "",
		DestinationPrefix: "thumb_",
		Scaler:            "CatmullRom",
	}

	gen := NewGenerator(config)

	// Load image with EXIF orientation 6 (90 degrees CW)
	// The original image is 400x600 (portrait), but with orientation 6
	// it should be treated as 600x400 (landscape) after rotation
	i, err := gen.NewImageFromFile(testImageOrientation6Path)
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

	err = os.WriteFile(dest, thumbBytes, 0644)
	if err != nil {
		t.Error(err)
	}

	checkFileExists(t, dest)
	
	// After applying EXIF orientation 6, the 400x600 image becomes 600x400
	// The thumbnail should maintain the aspect ratio of the rotated image
	gotWidth, gotHeight, err := checkImageDimensions(dest)
	if err != nil {
		t.Error(err)
	}

	// The rotated image (600x400) should produce a thumbnail where width > height
	if gotWidth <= gotHeight {
		t.Errorf("Expected landscape thumbnail (width > height) after EXIF rotation, got %dx%d", gotWidth, gotHeight)
	}

	t.Logf("Thumbnail dimensions: %dx%d (correctly rotated)", gotWidth, gotHeight)
}

// TestOrientationFunctions tests the individual orientation transformation functions
func TestOrientationFunctions(t *testing.T) {
	config := Generator{
		Scaler: "CatmullRom",
	}
	gen := NewGenerator(config)

	// Load the test image
	i, err := gen.NewImageFromFile(testImageOrientation6Path)
	if err != nil {
		t.Fatal(err)
	}

	// Test that getImageOrientation correctly reads EXIF data
	orientation := getImageOrientation(i.Data)
	if orientation != 6 {
		t.Errorf("Expected orientation 6, got %d", orientation)
	}
	t.Logf("Correctly read EXIF orientation: %d", orientation)

	// Test PNG images return orientation 1 (no EXIF)
	pngImage, err := gen.NewImageFromFile(testPngImagePath)
	if err != nil {
		t.Fatal(err)
	}
	pngOrientation := getImageOrientation(pngImage.Data)
	if pngOrientation != 1 {
		t.Errorf("Expected PNG to return orientation 1, got %d", pngOrientation)
	}
}
