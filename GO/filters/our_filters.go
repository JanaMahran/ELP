package filters

import (

	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"sync"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// decode image en une matrice de pixels 
func decodeImage(reader io.Reader) ([][][4]uint8, image.Image, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur decoder image: %v", err)
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

// encode matrice de pixels en une image 
func encodeImage(writer io.Writer, matrix [][][4]uint8) error {
	height := len(matrix)
	if height == 0 {
		return fmt.Errorf("matrice est vide")
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
		return fmt.Errorf("erreur encoder image en jpg: %v", err)
	}
	return nil
}

// Multiplie une matrice et une kernel 
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

	// assurer que les valeurs sont appartiennent au range 0-255
	clamp := func(v int) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v)
	}

	return [4]uint8{clamp(r), clamp(g), clamp(b), matrix[y][x][3]} 
}

// applique la multiplication avec kernel aux matrices avec des goroutines
func applyKernelParallel(matrix [][][4]uint8, kernel [][]int) [][][4]uint8 {
	height := len(matrix)
	width := len(matrix[0])

	// Create the output matrix
	output := make([][][4]uint8, height)
	for i := range output {
		output[i] = make([][4]uint8, width)
	}

	var wg sync.WaitGroup
	numWorkers := 2000
	rowsPerWorker := height / numWorkers

	worker := func(start, end int) {
		defer wg.Done()
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				output[y][x] = applyKernel(matrix, kernel, x, y)
			}
		}
	}

	// pour commencer les goroutines 
	for i := 0; i < numWorkers; i++ {
		start := i * rowsPerWorker
		end := start + rowsPerWorker
		if i == numWorkers-1 {
			end = height
		}
		wg.Add(1)
		go worker(start, end)
	}

	wg.Wait()
	return output
}

// applique filtre a une image et enregistre le output 
func ApplyFilters(kerType int, inputPath string, outputPath string) {
	// ouvrir image
	reader, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Erreur pour ouvrir fichier: %v", err)
	}
	defer reader.Close()

	// decoder image
	matrix, _, err := decodeImage(reader)
	if err != nil {
		log.Fatalf("Erreur pour decoder image: %v", err)
	}

	// definir le kernel selon le choix KerType
	var kernel [][]int
	switch kerType {
	case 1:
		fmt.Println("detection de bordure kernel choisis...")
		kernel = [][]int{
			{-1, -1, -1},
			{-1, 8, -1},
			{-1, -1, -1},
		}
	case 2:
		fmt.Println("sharpen kernel choisis...")
		kernel = [][]int{
			{0, -1, 0},
			{-1, 5, -1},
			{0, -1, 0},
		}
	default:
		log.Fatalf(" Kernel type invalide: %d", kerType)
	}

	// Appliquer kernel
	outputMatrix := applyKernelParallel(matrix, kernel)

	// creer fichier
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Error pour créer output file: %v", err)
	}
	defer outputFile.Close()

	// encoder l'image
	err = encodeImage(outputFile, outputMatrix)
	if err != nil {
		log.Fatalf("Erreur pour encoder image: %v", err)
	}

	fmt.Printf("Image compltè enregistré dans %s\n", outputPath)
}



