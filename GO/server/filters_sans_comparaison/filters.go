package filters

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"sync"
	"time"
)

// ApplyFilters applique le filtre sélectionné à l'image d'entrée et sauvegarde le résultat
func ApplyFilters(filterType int, inputPath, outputPath string) error {
	// Ouvrir l'image
	reader, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture de l'image : %w", err)
	}
	defer reader.Close()

	// Décoder l'image en fonction de son extension
	switch {
	case strings.HasSuffix(strings.ToLower(inputPath), ".gif"):
		return processGIF(filterType, reader, outputPath)
	case strings.HasSuffix(strings.ToLower(inputPath), ".jpg"), strings.HasSuffix(strings.ToLower(inputPath), ".jpeg"):
		img, err := jpeg.Decode(reader)
		if err != nil {
			return fmt.Errorf("erreur lors du décodage de l'image : %w", err)
		}
		return processImage(filterType, img, outputPath)
	case strings.HasSuffix(strings.ToLower(inputPath), ".png"):
		img, err := png.Decode(reader)
		if err != nil {
			return fmt.Errorf("erreur lors du décodage de l'image : %w", err)
		}
		return processImage(filterType, img, outputPath)
	default:
		return fmt.Errorf("format d'image non supporté : %s", inputPath)
	}
}

// processGIF traite un fichier GIF
func processGIF(filterType int, reader *os.File, outputPath string) error {
	// Décoder le GIF
	gifImg, err := gif.DecodeAll(reader)
	if err != nil {
		return fmt.Errorf("erreur lors du décodage du GIF : %w", err)
	}

	// Traiter chaque frame du GIF
	processedFrames := make([]*image.Paletted, len(gifImg.Image))
	for i, frame := range gifImg.Image {
		// Convertir la frame en image.Image
		img := image.NewRGBA(frame.Bounds())
		draw.Draw(img, img.Bounds(), frame, frame.Bounds().Min, draw.Src)

		// Appliquer le filtre
		processedImg, err := applyFilterToImage(filterType, img)
		if err != nil {
			return fmt.Errorf("erreur lors du traitement de la frame %d : %w", i, err)
		}

		// Convertir l'image traitée en image.Paletted
		palettedImg := image.NewPaletted(processedImg.Bounds(), frame.Palette)
		draw.Draw(palettedImg, palettedImg.Bounds(), processedImg, processedImg.Bounds().Min, draw.Src)
		processedFrames[i] = palettedImg
	}

	// Sauvegarder le GIF traité
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier de sortie : %w", err)
	}
	defer outputFile.Close()

	err = gif.EncodeAll(outputFile, &gif.GIF{
		Image:     processedFrames,
		Delay:     gifImg.Delay,
		LoopCount: gifImg.LoopCount,
	})
	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage du GIF : %w", err)
	}

	return nil
}

// processImage traite une image statique (JPEG, PNG)
func processImage(filterType int, img image.Image, outputPath string) error {
	// Appliquer le filtre
	processedImg, err := applyFilterToImage(filterType, img)
	if err != nil {
		return fmt.Errorf("erreur lors du traitement de l'image : %w", err)
	}

	// Sauvegarder l'image traitée
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier de sortie : %w", err)
	}
	defer outputFile.Close()

	switch outputPath[len(outputPath)-3:] {
	case "jpg", "jpeg":
		err = jpeg.Encode(outputFile, processedImg, nil)
	case "png":
		err = png.Encode(outputFile, processedImg)
	default:
		return fmt.Errorf("format de sortie non supporté : %s", outputPath)
	}
	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage de l'image : %w", err)
	}

	return nil
}

// applyFilterToImage applique un filtre à une image
func applyFilterToImage(filterType int, img image.Image) (image.Image, error) {
	// Convertir l'image en matrice de pixels
	matrix := imageToMatrix(img)

	// Définir le kernel en fonction du filtre sélectionné
	var kernel [][]float64
	switch filterType {
	case 1:
		kernel = [][]float64{ // Grayscale (pas de kernel, conversion directe)
			{0.299, 0.587, 0.114},
		}
	case 2:
		kernel = [][]float64{ // Détection de contours
			{-1, -1, -1},
			{-1, 8, -1},
			{-1, -1, -1},
		}
	case 3:
		kernel = [][]float64{ // Netteté
			{0, -1, 0},
			{-1, 5, -1},
			{0, -1, 0},
		}
	case 4:
		kernel = [][]float64{ // Flou gaussien
			{1 / 16.0, 2 / 16.0, 1 / 16.0},
			{2 / 16.0, 4 / 16.0, 2 / 16.0},
			{1 / 16.0, 2 / 16.0, 1 / 16.0},
		}
	default:
		return nil, fmt.Errorf("filtre non reconnu : %d", filterType)
	}

	// Appliquer le kernel à la matrice de pixels
	var outputMatrix [][][4]uint8
	if filterType == 1 {
		// Conversion directe en niveaux de gris
		outputMatrix = applyGrayscale(matrix)
	} else {
		startParallel := time.Now()
		outputMatrix = applyKernelParallel(matrix, kernel)
		elapsedParallel := time.Since(startParallel)
		fmt.Printf("Temps d'exécution avec goroutines : %v\n", elapsedParallel)
	}

	// Convertir la matrice de pixels en image
	return matrixToImage(outputMatrix), nil
}

// imageToMatrix convertit une image en une matrice de pixels
func imageToMatrix(img image.Image) [][][4]uint8 {
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
	return matrix
}

// matrixToImage convertit une matrice de pixels en une image
func matrixToImage(matrix [][][4]uint8) *image.RGBA {
	height := len(matrix)
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
	return img
}

// applyGrayscale convertit une matrice de pixels en niveaux de gris
func applyGrayscale(matrix [][][4]uint8) [][][4]uint8 {
	height := len(matrix)
	width := len(matrix[0])

	output := make([][][4]uint8, height)
	for y := 0; y < height; y++ {
		output[y] = make([][4]uint8, width)
		for x := 0; x < width; x++ {
			pixel := matrix[y][x]
			gray := uint8(0.299*float64(pixel[0]) + 0.587*float64(pixel[1]) + 0.114*float64(pixel[2]))
			output[y][x] = [4]uint8{gray, gray, gray, pixel[3]}
		}
	}
	return output
}

// applyKernelParallel applique un kernel à une matrice de pixels en parallèle
func applyKernelParallel(matrix [][][4]uint8, kernel [][]float64) [][][4]uint8 {
	height := len(matrix)
	width := len(matrix[0])

	output := make([][][4]uint8, height)
	for i := range output {
		output[i] = make([][4]uint8, width)
	}

	var wg sync.WaitGroup
	numWorkers := 4 // Nombre de goroutines
	rowsPerWorker := height / numWorkers

	worker := func(start, end int) {
		defer wg.Done()
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				output[y][x] = applyKernel(matrix, kernel, x, y)
			}
		}
	}

	// Démarrer les goroutines
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

// applyKernel applique un kernel à un pixel spécifique
func applyKernel(matrix [][][4]uint8, kernel [][]float64, x, y int) [4]uint8 {
	height := len(matrix)
	width := len(matrix[0])
	kernelSize := len(kernel)
	offset := kernelSize / 2

	var r, g, b float64
	for ky := 0; ky < kernelSize; ky++ {
		for kx := 0; kx < kernelSize; kx++ {
			px := x + kx - offset
			py := y + ky - offset
			if px >= 0 && px < width && py >= 0 && py < height {
				pixel := matrix[py][px]
				r += float64(pixel[0]) * kernel[ky][kx]
				g += float64(pixel[1]) * kernel[ky][kx]
				b += float64(pixel[2]) * kernel[ky][kx]
			}
		}
	}

	// Limiter les valeurs des pixels à l'intervalle 0 à 255
	clamp := func(v float64) uint8 {
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
