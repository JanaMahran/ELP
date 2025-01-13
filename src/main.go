package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "src/EncodeDecode"
    "image"
	"image/color"
	"image/png"
	"os"
	"log"
)

func main() {
    // to define paths for images, root, etc
    projectRoot, err := os.Getwd()
    if err != nil {
        log.Fatalf("Error determining working directory: %v", err)
    }   
    inputImagePath := filepath.Join(projectRoot, "annecy.jpg")
    outputImagePath := filepath.Join(projectRoot, "annecy_processed.png")

    // decode
    fmt.Println("Decoding image...")
    matrix, _, err := decodeImage(inputImagePath)
    if err != nil {
        log.Fatalf("Error decoding image: %v", err)
    }
    fmt.Println("Image decoded successfully!")

    // encode
    fmt.Println("Encoding image...")
    err = encodeImage(outputImagePath, matrix)
    if err != nil {
        log.Fatalf("Error encoding image: %v", err)
    }
    fmt.Println("Image encoded successfully!")

    fmt.Printf("Processed image saved to: %s\n", outputImagePath)
}

