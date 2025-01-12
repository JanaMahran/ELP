package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"log"
)


// pour decoder image 
func DecodeImage(filename string) ([][][4]uint8, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create a 3D matrix [height][width][4] for RGBA
	matrix := make([][][4]uint8, height)
	for y := 0; y < height; y++ {
		matrix[y] = make([][4]uint8, width)
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			matrix[y][x] = [4]uint8{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
				uint8(a >> 8),
			}
		}
	}
	return matrix, nil
}


// pour encoder image
func EncodeImage(matrix [][][4]uint8, filename string) error {
	height := len(matrix)
	width := len(matrix[0])

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgba := matrix[y][x]
			img.SetRGBA(x, y, color.RGBA{
				R: rgba[0],
				G: rgba[1],
				B: rgba[2],
				A: rgba[3],
			})
		}
	}

	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return png.Encode(outFile, img)
}

