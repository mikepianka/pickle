package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"image/color"
	_ "image/jpeg"
	_ "image/png"
)

// config holds the arguments parsed from the command line.
type config struct {
	imagePath string
	x         int
	y         int
	aFloat    bool
}

func main() {
	conf := cli()

	img, err := loadImage(conf.imagePath)
	if err != nil {
		log.Fatal(err)
	}

	rgba, err := getPixelValue(img, conf.x, conf.y)
	if err != nil {
		log.Fatal(err)
	}

	if conf.aFloat {
		a := alphaIntToFloat(rgba.A)
		fmt.Printf("Pixel at (%d, %d) has RGBA values of R: %d, G: %d, B: %d, A: %f\n", conf.x, conf.y, rgba.R, rgba.G, rgba.B, a)
	} else {
		fmt.Printf("Pixel at (%d, %d) has RGBA values of R: %d, G: %d, B: %d, A: %d\n", conf.x, conf.y, rgba.R, rgba.G, rgba.B, rgba.A)
	}
}

// cli parses the command line arguments and returns a config struct.
func cli() config {
	imagePath := flag.String("path", "", "Path to JPG or PNG file")
	x := flag.Int("x", 0, "X coordinate of the pixel")
	y := flag.Int("y", 0, "Y coordinate of the pixel")
	aFloat := flag.Bool("af", false, "Optional: Return alpha value as float instead of int (default false)")

	flag.Parse()

	return config{
		imagePath: *imagePath,
		x:         *x,
		y:         *y,
		aFloat:    *aFloat,
	}
}

// validateFilepath checks if the provided path is a valid file path and returns an error if it is not.
func validateFilepath(path string) error {
	if path == "" {
		return fmt.Errorf("no file path provided")
	}

	info, err := os.Stat(path)
	if err != nil {
		// path doesn't exist, permissions issues, etc.
		return err
	}

	if info.IsDir() {
		// passing a directory path should return an error
		return fmt.Errorf("expected file path but directory path was provided")
	}

	return nil
}

// validateFiletype checks if the provided file path matches one of the allowed file extensions and returns an error if it does not.
func validateFiletype(path string, allowedExtensions []string) error {
	ext := filepath.Ext(path)

	for _, e := range allowedExtensions {
		if strings.EqualFold(ext, e) {
			return nil
		}
	}

	return fmt.Errorf("file type %s is not one of the supported types: %v", ext, allowedExtensions)
}

// validatePixelCoords checks if the provided pixel coordinates are within the bounds of the image and returns an error if they are not.
func validatePixelCoords(img image.Image, x, y int) error {
	bounds := img.Bounds()

	if x < bounds.Min.X || x > bounds.Max.X || y < bounds.Min.Y || y > bounds.Max.Y {
		return fmt.Errorf("coordinates are out of image bounds")
	}

	return nil
}

// loadImage loads and returns the image at the provided path.
// An error is returned if the file path is invalid, the file type is not supported, or the file cannot be decoded.
func loadImage(path string) (image.Image, error) {
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}

	if err := validateFiletype(path, allowedExtensions); err != nil {
		return nil, err
	}

	if err := validateFilepath(path); err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// getPixelValue returns the NRGBA values of the pixel at the provided coordinates.
// An error is returned if the coordinates are out of bounds of the provided image.
func getPixelValue(img image.Image, x, y int) (color.NRGBA, error) {
	if err := validatePixelCoords(img, x, y); err != nil {
		return color.NRGBA{}, err
	}

	pixel := img.At(x, y)
	rgba := color.NRGBAModel.Convert(pixel).(color.NRGBA)
	return rgba, nil
}

// alphaIntToFloat converts an alpha value from uint8 to float64.
func alphaIntToFloat(alpha uint8) float64 {
	return float64(alpha) / 255
}
