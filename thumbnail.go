package thumbnail

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

// An Image is an image and it's information.
type Image struct {
	Filename    string
	ContentType string
	Data        []byte
	Size        int
}

// A Configuration sets all the configurable options for
// thumbnail generation.
type Configuration struct {
	Path              string
	ContentType       string
	DestinationPrefix string
}

// Create generates ta thumbnail.
func Create(fb []byte, cfg Configuration) error {
	i, err := process(cfg.Path, fb, cfg.ContentType)
	if err != nil {
		return err
	}
	thumbPath := cfg.DestinationPrefix + filepath.Base(cfg.Path)

	dst := createRect(i)
	var buffer bytes.Buffer
	switch i.ContentType {
	case "image/jpeg":
		err := jpeg.Encode(&buffer, dst, nil)
		if err != nil {
			return err
		}
		err = writeFile(thumbPath, buffer.Bytes(), 0644)
		if err != nil {
			return err
		}
	case "image/png":
		err := png.Encode(&buffer, dst)
		if err != nil {
			return err
		}
		err = writeFile(thumbPath, buffer.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func process(path string, fb []byte, contentType string) (*Image, error) {
	_, _, err := image.Decode(bytes.NewReader(fb))
	if err != nil {
		return nil, err
	}

	i := &Image{
		Filename:    filepath.Base(path),
		ContentType: contentType,
		Data:        fb,
		Size:        len(fb),
	}

	return i, nil
}

func createRect(i *Image) *image.RGBA {
	img, _, err := image.Decode(bytes.NewReader(i.Data))
	if err != nil {
		log.Print(err)
	}
	var (
		height = img.Bounds().Max.Y
		width  = img.Bounds().Max.X
		y      = 300
		x      = 300 * width / height
	)
	rect := image.Rect(0, 0, x, y)
	dst := image.NewRGBA(rect)
	// scaler can be one of:
	//  - CatmullRom        - best quality, slowest
	//  - ApproxBiLinear    - mid quality, mid-speed
	//  - NearestNeighbor   - low quality, fast
	scaler := draw.CatmullRom
	scaler.Scale(dst, rect, img, img.Bounds(), draw.Over, nil)
	return dst

}

func writeFile(p string, f []byte, fmode os.FileMode) error {
	return ioutil.WriteFile(p, f, fmode)
}
