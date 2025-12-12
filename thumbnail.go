// Package thumbnail provides a method to create thumbnails from images.
package thumbnail

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

	"github.com/rwcarlsen/goexif/exif"
	"golang.org/x/image/draw"
)

// An Image is an image and information about it.
type Image struct {
	// Path is a path to an image.
	Path string

	// ContentType is the content type of the image.
	ContentType string

	// Data is the image data in a byte-array
	Data []byte

	// Size is the length of Data
	Size int

	// Current stores the existing image's dimensions
	Current Dimensions

	// Future store the new thumbnail dimensions.
	Future Dimensions
}

// Dimensions stores dimensional information for an Image.
type Dimensions struct {
	// Width is the width of an image in pixels.
	Width int

	// Height is the height on an image in pixels.
	Height int

	// X is the right-most X-coordinate.
	X int

	// Y is the top-most Y-coordinate.
	Y int
}

var (
	// ErrInvalidMimeType is returned when a non-image content type is
	// detected.
	ErrInvalidMimeType = errors.New("invalid mimetype")

	// ErrInvalidScaler is returned when an unrecognized scaler is
	// passed to the Generator.
	ErrInvalidScaler = errors.New("invalid scaler")
)

// NewGenerator returns an instance of a thumbnail generator with a
// given configuration.
func NewGenerator(c Generator) *Generator {
	return &Generator{
		Width:             300,
		Height:            300,
		DestinationPath:   c.DestinationPath,
		DestinationPrefix: c.DestinationPrefix,
		Scaler:            c.Scaler,
	}
}

// NewImageFromFile reads in an image file from the file system and
// populates an Image object. That new Image object is returned along
// with any errors that occur during the operation.
func (gen *Generator) NewImageFromFile(path string) (*Image, error) {
	imageBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	contentType := detectContentType(imageBytes)
	return &Image{
		Path:        path,
		ContentType: contentType,
		Data:        imageBytes,
		Size:        len(imageBytes),
		Current: Dimensions{
			Width:  0,
			Height: 0,
		},
		Future: Dimensions{
			Width:  gen.Width,
			Height: gen.Height,
		},
	}, nil
}

// NewImageFromByteArray reads in an image from a byte array and
// populates an Image object. That new Image object is returned along
// with any errors that occur during the operation.
func (gen *Generator) NewImageFromByteArray(imageBytes []byte) (*Image, error) {
	contentType := detectContentType(imageBytes)
	return &Image{
		ContentType: contentType,
		Data:        imageBytes,
		Size:        len(imageBytes),
		Current: Dimensions{
			Width:  0,
			Height: 0,
		},
		Future: Dimensions{
			Width:  gen.Width,
			Height: gen.Height,
		},
	}, nil
}

// Generator registers a generator configuration to be used when
// creating thumbnails.
type Generator struct {
	// Width is the destination thumbnail width.
	Width int

	// Height is the destination thumbnail height.
	Height int

	// DestinationPath is the destination thumbnail path.
	DestinationPath string

	// DestinationPrefix is the prefix for the destination thumbnail
	// filename.
	DestinationPrefix string

	// Scaler is the scaler to be used when generating thumbnails.
	Scaler string
}

// CreateThumbnail generates a thumbnail.
func (gen *Generator) CreateThumbnail(i *Image) ([]byte, error) {
	if i.ContentType == "application/octet-stream" {
		return nil, ErrInvalidMimeType
	}

	dst, err := gen.createRect(i)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	switch i.ContentType {
	case "image/jpeg":
		err = jpeg.Encode(&buffer, dst, nil)
	case "image/png":
		err = png.Encode(&buffer, dst)
	}
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (gen *Generator) createRect(i *Image) (*image.RGBA, error) {
	img, _, err := image.Decode(bytes.NewReader(i.Data))
	if err != nil {
		return nil, err
	}

	// Apply EXIF orientation transformation
	orientation := getImageOrientation(i.Data)
	img = applyOrientation(img, orientation)

	var (
		width  = img.Bounds().Max.X
		height = img.Bounds().Max.Y
		x      = gen.Width * width / height
		y      = gen.Height
	)
	gen.Width = x
	gen.Height = y
	rect := image.Rect(0, 0, x, y)
	dst := image.NewRGBA(rect)
	var scaler draw.Interpolator
	switch scalerChoice := gen.Scaler; scalerChoice {
	case "NearestNeighbor":
		scaler = draw.NearestNeighbor
	case "ApproxBiLinear":
		scaler = draw.ApproxBiLinear
	case "BiLinear":
		scaler = draw.BiLinear
	case "CatmullRom":
		scaler = draw.CatmullRom
	}
	if scaler == nil {
		return nil, ErrInvalidScaler
	}
	scaler.Scale(dst, rect, img, img.Bounds(), draw.Over, nil)
	return dst, nil

}

// getImageOrientation extracts the EXIF orientation tag from an image.
// Returns 1 (normal orientation) if no orientation tag is found or if there's an error.
func getImageOrientation(data []byte) int {
	// Only attempt to read EXIF data for JPEG images
	contentType := http.DetectContentType(data[:min(512, len(data))])
	if contentType != "image/jpeg" {
		return 1
	}

	reader := bytes.NewReader(data)
	x, err := exif.Decode(reader)
	if err != nil {
		// No EXIF data or unable to decode, return normal orientation
		return 1
	}

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		// No orientation tag, return normal orientation
		return 1
	}

	orientationVal, err := orientation.Int(0)
	if err != nil {
		return 1
	}

	return orientationVal
}

// applyOrientation applies the EXIF orientation transformation to an image.
// The orientation parameter should be the EXIF orientation tag value (1-8).
// Reference: http://jpegclub.org/exif_orientation.html
func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 1:
		// Normal orientation, no transformation needed
		return img
	case 2:
		// Flipped horizontally
		return flipHorizontal(img)
	case 3:
		// Rotated 180 degrees
		return rotate180(img)
	case 4:
		// Flipped vertically
		return flipVertical(img)
	case 5:
		// Transposed (flipped over top-left to bottom-right axis)
		// = Flipped horizontally then rotated 90° CCW
		return rotate90(flipHorizontal(img))
	case 6:
		// Rotated 90° CW
		// Note: rotate270 performs 270° CCW rotation, which equals 90° CW
		return rotate270(img)
	case 7:
		// Transverse (flipped over top-right to bottom-left axis)
		// = Flipped horizontally then rotated 90° CW
		return rotate270(flipHorizontal(img))
	case 8:
		// Rotated 90° CCW
		// Note: rotate90 performs 90° CCW rotation
		return rotate90(img)
	default:
		// Unknown orientation, return original
		return img
	}
}

// rotate90 rotates an image 90 degrees counter-clockwise
func rotate90(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	rotated := image.NewRGBA(image.Rect(0, 0, height, width))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rotated.Set(y-bounds.Min.Y, width-(x-bounds.Min.X)-1, img.At(x, y))
		}
	}
	return rotated
}

// rotate180 rotates an image 180 degrees
func rotate180(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	rotated := image.NewRGBA(image.Rect(0, 0, width, height))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rotated.Set(width-(x-bounds.Min.X)-1, height-(y-bounds.Min.Y)-1, img.At(x, y))
		}
	}
	return rotated
}

// rotate270 rotates an image 270 degrees counter-clockwise (or 90 degrees clockwise)
func rotate270(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	rotated := image.NewRGBA(image.Rect(0, 0, height, width))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rotated.Set(height-(y-bounds.Min.Y)-1, x-bounds.Min.X, img.At(x, y))
		}
	}
	return rotated
}

// flipHorizontal flips an image horizontally
func flipHorizontal(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	flipped := image.NewRGBA(image.Rect(0, 0, width, height))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			flipped.Set(width-(x-bounds.Min.X)-1, y-bounds.Min.Y, img.At(x, y))
		}
	}
	return flipped
}

// flipVertical flips an image vertically
func flipVertical(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	flipped := image.NewRGBA(image.Rect(0, 0, width, height))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			flipped.Set(x-bounds.Min.X, height-(y-bounds.Min.Y)-1, img.At(x, y))
		}
	}
	return flipped
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// detectContentType from
// https://golangcode.com/get-the-content-type-of-file/
func detectContentType(fb []byte) string {
	// Only the first 512 bytes are used to sniff the content type.
	// Use the net/http package's handy DetectContentType function.
	// Always seems to return a valid content-type by returning
	// "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(fb[:512])
}
