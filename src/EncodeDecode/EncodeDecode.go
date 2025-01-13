package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"log"
)


// pour decoder image 

func decodeImage(filepath string) ([][][4]uint8, image.Image, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to open image: %v", err)
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to decode image: %v", err)
    }

    bounds := img.Bounds()
    width, height := bounds.Dx(), bounds.Dy()

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

    return matrix, img, nil
}



// pour encoder image
func encodeImage(filepath string, matrix [][][4]uint8) error {
    height := len(matrix)
    if height == 0 {
        return fmt.Errorf("matrix is empty")
    }
    width := len(matrix[0])

    img := image.NewRGBA(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            pixel := matrix[y][x]
            img.Set(x, y, color.RGBA{
                R: pixel[0],
                G: pixel[1],
                B: pixel[2],
                A: pixel[3],
            })
        }
    }

    file, err := os.Create(filepath)
    if err != nil {
        return fmt.Errorf("failed to create image file: %v", err)
    }
    defer file.Close()

    if err := png.Encode(file, img); err != nil {
        return fmt.Errorf("failed to encode image to PNG: %v", err)
    }
    return nil
}


