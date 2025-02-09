package filters

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"sync"
	"time"
)

// ApplyFilters permet l'ouverture du fichier image, la détermination de son format, et de son décodage en un objet image.Image
// elle englobe l'application du filtre sélectionné à l'image d'entrée selon son extension et la sauvegarde le résultat, en fournissant les paramètres nécessaires à processImage
func ApplyFilters(filterType int, inputPath, outputPath string) error {
	// Ouverture de l'image
	reader, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture de l'image : %w", err)
	}
	defer reader.Close()

	// Décodage de l'image en fonction de son extension (jpg ou png!)
	switch {
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

// applique filtre sur image et enregistre le résultat
func processImage(filterType int, img image.Image, outputPath string) error {
	processedImg, err := applyFilterToImage(filterType, img)
	if err != nil {
		return fmt.Errorf("erreur lors du traitement de l'image : %w", err)
	}

	// Sauvegarde de l'image traitée
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

// applyFilterToImage applique le filtre en entrée à une image
func applyFilterToImage(filterType int, img image.Image) (image.Image, error) {
	// On convertit d'abord l'image en une matrice de pixels pour pouvoir agir dessus
	matrix := imageToMatrix(img)

	// Définition du kernel qui sera utilisé en fonction du filtre sélectionné
	var kernel [][]float64
	switch filterType {
	case 1:
		kernel = [][]float64{ // Grayscale (dans ce cas pas de kernel, conversion directe)
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
		//on met .0 pour avoir des floats car en go division de deux entiers donne entier
	default:
		return nil, fmt.Errorf("filtre non reconnu : %d", filterType)
	}

	// Application du kernel à la matrice de pixels
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

	// On reconvertit la matrice de pixels en image
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

// applyGrayscale convertit une matrice de pixels en niveaux de gris $
func applyGrayscale(matrix [][][4]uint8) [][][4]uint8 {
	height := len(matrix)
	width := len(matrix[0])

	output := make([][][4]uint8, height)
	for y := 0; y < height; y++ {
		output[y] = make([][4]uint8, width)
		for x := 0; x < width; x++ {
			pixel := matrix[y][x]
			gray := uint8(0.299*float64(pixel[0]) + 0.587*float64(pixel[1]) + 0.114*float64(pixel[2]))
			//on obtient les valeurs de gris en pondérant avec des coefficients connus la valeur des canaux r g et b (choix de ne pas moyenner pour ce filtre)
			output[y][x] = [4]uint8{gray, gray, gray, pixel[3]}
		}
	}
	return output
}

// applyKernelParallel applique un kernel à une matrice de pixels en parallèle, avec plusieurs goroutines qui gèrent chacune un bout de la matruce
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

	// Démarrage des goroutines
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

	// on limite les valeurs des pixels à l'intervalle 0 à 255
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
