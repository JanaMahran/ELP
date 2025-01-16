package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"image/color"
	"image/png"
	"io"
	"sync"


	// image.Decode peut grâce aux importations suivantes comprendre les images de format JPEG, PNG et GIF
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)


// DecodeImage decodes an image from an io.Reader and returns a pixel matrix.
func decodeImage(reader io.Reader) ([][][4]uint8, image.Image, error) {
	img, _, err := image.Decode(reader)
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

// EncodeImage encodes a pixel matrix to an image and writes it to an io.Writer.
func encodeImage(writer io.Writer, matrix [][][4]uint8) error {
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

	if err := png.Encode(writer, img); err != nil {
		return fmt.Errorf("failed to encode image to PNG: %v", err)
	}
	return nil
}

func printMatrix(matrix [][][4]uint8) {
	for y, row := range matrix {
		fmt.Printf("Row %d: ", y)
		for _, pixel := range row {
			fmt.Printf("[%d %d %d %d] ", pixel[0], pixel[1], pixel[2], pixel[3])
		}
		fmt.Println() // Move to the next row
	}
}

func applyKernel(matrix [][][4]uint8, kernel [][]int, x, y int) [4]uint8 {
	height := len(matrix)
	width := len(matrix[0])
	kernelSize := len(kernel)
	offset := kernelSize / 2

	var r, g, b int
	for ky := 0; ky < kernelSize; ky++ {
		for kx := 0; kx < kernelSize; kx++ {
			px := x + kx - offset
			py := y + ky - offset
			if px >= 0 && px < width && py >= 0 && py < height {
				pixel := matrix[py][px]
				r += int(pixel[0]) * kernel[ky][kx]
				g += int(pixel[1]) * kernel[ky][kx]
				b += int(pixel[2]) * kernel[ky][kx]
			}
		}
	}

	// Clamp the values to the 0-255 range
	clamp := func(v int) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v)
	}

	return [4]uint8{clamp(r), clamp(g), clamp(b), matrix[y][x][3]} // Preserve alpha channel
}

func applyKernelParallel(matrix [][][4]uint8, kernel [][]int) [][][4]uint8 {
	height := len(matrix)
	width := len(matrix[0])
	output := make([][][4]uint8, height)
	for i := range output {
		output[i] = make([][4]uint8, width)
	}

	// worker function to process a range of rows
	worker := func(start, end int, done chan<- bool) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				output[y][x] = applyKernel(matrix, kernel, x, y)
			}
		}
		done <- true
	}

	// number of goroutines
	numWorkers := 4
	rowsPerWorker := height / numWorkers
	done := make(chan bool, numWorkers)

	// start goroutines
	for i := 0; i < numWorkers; i++ {
		start := i * rowsPerWorker
		end := start + rowsPerWorker
		if i == numWorkers-1 {
			end = height // ensure the last worker processes any remaining rows
		}
		go worker(start, end, done)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numWorkers; i++ {
		<-done
	}

	return output
}


func main() {

	kernel := [][]int{
		{-1, -1, -1},
		{-1, 8, -1},
		{-1, -1, -1},
	}


	// Open the image file
	reader, err := os.Open("annecy.jpg") // or "lyon_2.png" to test with a PNG file
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// Decode the image to check its dimensions
	im, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Get image dimensions
	bounds := im.Bounds()
	fmt.Printf("Dimensions de l'image : largeur=%d, hauteur=%d\n", bounds.Dx(), bounds.Dy())

	// Reset the reader by reopening the file
	reader.Seek(0, io.SeekStart)

	// Decode the image into a pixel matrix
	fmt.Println("Decoding the image...")
	matrix, _, err := decodeImage(reader)
	if err != nil {
		log.Fatalf("Error decoding image: %v", err)
	}
	fmt.Println("Image decoded successfully!")

	// for testing purposes
	// Print the pixel matrix
	// printMatrix(matrix)
	

	//apply filters 
	outputMatrix := applyKernelParallel(matrix, kernel)
	fmt.Println("Kernel applied successfully!")
	outputFile, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Encode the pixel matrix back into an image
	fmt.Println("Encoding the image...")
	err = encodeImage(outputFile, outputMatrix)
	if err != nil {
		log.Fatalf("Error encoding image: %v", err)
	}
	fmt.Println("Image encoded and saved successfully to output.png")


}
