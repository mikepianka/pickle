package main

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_alphaIntToFloat(t *testing.T) {
	tolerance := 1e-6

	// test cases
	data := []struct {
		name     string
		alpha    uint8
		expected float64
	}{
		{"0", 0, 0.0},
		{"255", 255, 1.0},
		{"128", 128, 0.501961},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := alphaIntToFloat(d.alpha)

			assert.InDelta(t, d.expected, result, tolerance, "alphaIntToFloat(%d) returned %f, expected %f", d.alpha, result, d.expected)
		})
	}
}

func Test_getPixelValue_ColorJPG(t *testing.T) {
	img, err := loadImage(filepath.Join("testdata", "gophers.jpg"))
	if err != nil {
		t.Fatalf("could not load image: %v", err)
	}

	// test cases
	data := []struct {
		name        string
		x           int
		y           int
		expected    color.NRGBA
		expectedErr error
	}{
		{"nose", 297, 85, color.NRGBA{R: 179, G: 192, B: 200, A: 255}, nil},
		{"body", 299, 130, color.NRGBA{R: 90, G: 218, B: 255, A: 255}, nil},
		{"out of bounds (pos)", 601, 240, color.NRGBA{R: 0, G: 0, B: 0, A: 0}, fmt.Errorf("coordinates are out of image bounds")},
		{"out of bounds (neg)", -1, -1, color.NRGBA{R: 0, G: 0, B: 0, A: 0}, fmt.Errorf("coordinates are out of image bounds")},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result, err := getPixelValue(img, d.x, d.y)

			if err != nil && strings.HasPrefix(d.name, "out of bounds") {
				assert.Equal(t, d.expectedErr, err, "getPixelValue(%d, %d) returned error %v, expected %v", d.x, d.y, err, d.expectedErr)
			} else if err != nil {
				t.Fatalf("getPixelValue(%d, %d) returned error %v, expected %v", d.x, d.y, err, d.expectedErr)
			}

			assert.Equal(t, d.expected, result, "getPixelValue(%d, %d) returned %v, expected %v", d.x, d.y, result, d.expected)
		})
	}
}
func Test_getPixelValue_GrayscalePNG(t *testing.T) {
	img, err := loadImage(filepath.Join("testdata", "gophers_grayscale.png"))
	if err != nil {
		t.Fatalf("could not load image: %v", err)
	}

	// test cases
	data := []struct {
		name        string
		x           int
		y           int
		expected    color.NRGBA
		expectedErr error
	}{
		{"nose", 297, 85, color.NRGBA{R: 189, G: 189, B: 189, A: 255}, nil},
		{"body", 299, 130, color.NRGBA{R: 184, G: 184, B: 184, A: 255}, nil},
		{"out of bounds (pos)", 601, 240, color.NRGBA{R: 0, G: 0, B: 0, A: 0}, fmt.Errorf("coordinates are out of image bounds")},
		{"out of bounds (neg)", -1, -1, color.NRGBA{R: 0, G: 0, B: 0, A: 0}, fmt.Errorf("coordinates are out of image bounds")},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result, err := getPixelValue(img, d.x, d.y)

			if err != nil && strings.HasPrefix(d.name, "out of bounds") {
				assert.Equal(t, d.expectedErr, err, "getPixelValue(%d, %d) returned error %v, expected %v", d.x, d.y, err, d.expectedErr)
			} else if err != nil {
				t.Fatalf("getPixelValue(%d, %d) returned error %v, expected %v", d.x, d.y, err, d.expectedErr)
			}

			assert.Equal(t, d.expected, result, "getPixelValue(%d, %d) returned %v, expected %v", d.x, d.y, result, d.expected)
		})
	}
}

func Test_validatePixelCoords(t *testing.T) {
	img, err := loadImage(filepath.Join("testdata", "gophers.jpg"))
	if err != nil {
		t.Fatalf("could not load image: %v", err)
	}

	// test cases
	data := []struct {
		name        string
		x           int
		y           int
		expectedErr error
	}{
		{"valid", 100, 100, nil},
		{"invalid (neg)", 1, -100, fmt.Errorf("coordinates are out of image bounds")},
		{"invalid (pos)", 1000, 1, fmt.Errorf("coordinates are out of image bounds")},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := validatePixelCoords(img, d.x, d.y)

			assert.Equal(t, d.expectedErr, err, "validatePixelCoords(%d, %d) returned error %v, expected %v", d.x, d.y, err, d.expectedErr)
		})
	}
}

func Test_validateFiletype(t *testing.T) {
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	allowedExtensionsCap := []string{".JPG", ".JPEG", ".PNG"}

	// test cases
	data := []struct {
		path              string
		allowedExtensions []string
		expectedErr       string
	}{
		{"a.jpg", allowedExtensions, ""},
		{"a.JPG", allowedExtensions, ""},
		{"a.jpeg", allowedExtensions, ""},
		{"a.JPEG", allowedExtensions, ""},
		{"a.png", allowedExtensions, ""},
		{"a.PNG", allowedExtensions, ""},
		{"./some/dir/a.jpg", allowedExtensions, ""},
		{"./some/dir/a.JPG", allowedExtensions, ""},
		{"a.jpg", allowedExtensionsCap, ""},
		{"a.JPG", allowedExtensionsCap, ""},
		{"a.jpeg", allowedExtensionsCap, ""},
		{"a.JPEG", allowedExtensionsCap, ""},
		{"a.png", allowedExtensionsCap, ""},
		{"a.PNG", allowedExtensionsCap, ""},
		{"./some/dir/a.jpg", allowedExtensionsCap, ""},
		{"./some/dir/a.JPG", allowedExtensionsCap, ""},
		{"a.txt", allowedExtensions, "not one of the supported types:"},
		{"a.txt", allowedExtensionsCap, "not one of the supported types:"},
		{"a.TXT", allowedExtensions, "not one of the supported types:"},
		{"a.TXT", allowedExtensionsCap, "not one of the supported types:"},
	}

	for _, d := range data {
		t.Run(d.path, func(t *testing.T) {
			err := validateFiletype(d.path, d.allowedExtensions)

			if err != nil {
				assert.Contains(t, err.Error(), d.expectedErr, "validateFiletype(%s) returned error %v, expected %v", d.path, err, d.expectedErr)
			}
		})
	}
}

func Test_validateFilepath(t *testing.T) {
	// create a temporary directory and file
	tmpDir, err := os.MkdirTemp("", "validateFilepath_test")
	if err != nil {
		t.Fatalf("could not create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile, err := os.CreateTemp(tmpDir, "validateFilepath_test_*.txt")
	if err != nil {
		t.Fatalf("could not create temporary file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// test cases
	data := []struct {
		name        string
		path        string
		expectedErr error
	}{
		{"empty path", "", fmt.Errorf("no file path provided")},
		{"dir path", tmpDir, fmt.Errorf("expected file path but directory path was provided")},
		{"nonexistent file", "nonexistent_file_path", os.ErrNotExist},
		{"valid file", tmpFile.Name(), nil},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := validateFilepath(d.path)

			if d.name == "nonexistent file" {
				assert.True(t, errors.Is(err, os.ErrNotExist), "validateFilepath(%q) should have returned an error indicating the file does not exist, but got: %v", d.path, err)
			} else {
				assert.Equal(t, d.expectedErr, err, "validateFilepath(%q) returned error %v, expected %v", d.path, err, d.expectedErr)
			}
		})
	}
}

func Test_loadImage(t *testing.T) {
	// test cases
	data := []struct {
		name          string
		path          string
		expectSuccess bool
	}{
		{"valid file", "testdata/gophers.jpg", true},
		{"valid png", "testdata/gophers_grayscale.png", true},
		{"dir path", "testdata", false},
		{"nonexistent", "testdata/noexist.png", false},
		{"unsupported type", "testdata/not_an_image.txt", false},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			_, err := loadImage(d.path)

			// the error messages will vary so we just check if we expect the case to succeed or not
			if d.expectSuccess {
				assert.Nil(t, err, "loadImage(%s) returned error %v, expected nil", d.path, err)
			} else {
				assert.NotNil(t, err, "loadImage(%s) returned nil error, expected error", d.path)
			}
		})
	}
}
