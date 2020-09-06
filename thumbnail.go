package thumbnail

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/image/draw"
)

// An Image is an image and it's information.
type Image struct {
	Path        string
	ContentType string
	Data        []byte
	Size        int
	Current     dimensions
	Future      dimensions
}

type dimensions struct {
	Width, Height, X, Y int
}

// A Configuration sets all the configurable options for
// thumbnail generation.
type Configuration struct {
	Path              string
	ContentType       string
	DestinationPrefix string
}

// NewGenerator creates a new thumbnail generator and its configuration.
func NewGenerator() *Generator {
	return &Generator{
		Width:             300,
		Height:            300,
		DestinationPath:   "",
		DestinationPrefix: "thumb_",
	}
}

// NewImage returns a new Image.
func (gen *Generator) NewImage(path string) (*Image, error) {
	imageBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	contentType := detectContentType(imageBytes)
	return &Image{
		Path:        path,
		ContentType: contentType,
		Data:        imageBytes,
		Size:        len(imageBytes),
		Current: dimensions{
			Width:  0,
			Height: 0,
		},
		Future: dimensions{
			Width:  gen.Width,
			Height: gen.Height,
		},
	}, nil
}

// Generator registers a geneator configuration to be used when creating
// thumbnails.
type Generator struct {
	// Width is the destination thumbnail width.
	Width int

	// Height is the destination thumbnail height.
	Height int

	// DestinationPath is the dentination thumbnail path.
	DestinationPath string

	// DestinationPrefix is the prefix for the destination thumbnail filename.
	DestinationPrefix string

	// Scaler is the scaler to be used when generating thumbnails.
	Scaler string
}

// Create generates ta thumbnail.
func (gen *Generator) Create(i *Image) ([]byte, error) {

	dst := gen.createRect(i)
	var buffer bytes.Buffer
	var err error
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

func (gen *Generator) createRect(i *Image) *image.RGBA {
	img, _, err := image.Decode(bytes.NewReader(i.Data))
	if err != nil {
		log.Print(err)
	}
	var (
		x = gen.Width
		y = gen.Height
	)
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
		log.Print("Error: No scaler selected.")
	}
	scaler.Scale(dst, rect, img, img.Bounds(), draw.Over, nil)
	return dst

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
