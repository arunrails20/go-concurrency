package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to send directory path of images")
	}

	start := time.Now()

	err := walkFiles(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Time taken: %s\n", time.Since(start))
}

func walkFiles(root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// check if it is a file

		if !info.Mode().IsRegular() {
			return nil
		}

		// check if it is image/jpg

		contentType, _ := getFileContentType(path)
		if contentType != "image/jpeg" {
			return nil
		}
		// Process the image

		thumbnailImage, err := processImage(path)

		if err != nil {
			return err
		}

		// Save the thumbnail image to disk

		err = saveThumbnail(path, thumbnailImage)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil

}

// processImage - take image file as input
// return pointer to thumbnail image to memory
func processImage(path string) (*image.NRGBA, error) {
	srcImage, err := imaging.Open(path)

	// load the image from file
	if err != nil {
		return nil, err
	}

	// scale the image to 100px * 100px
	thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

	return thumbnailImage, nil
}

// saveThumbnail - save the thumnail image to folder
func saveThumbnail(srcImagePath string, thumbnailImage *image.NRGBA) error {
	filename := filepath.Base(srcImagePath)
	dstImagePath := "thumbnail/" + filename

	// save the image in the thumbnail folder.
	err := imaging.Save(thumbnailImage, dstImagePath)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", srcImagePath, dstImagePath)
	return nil

}

func getFileContentType(file string) (string, error) {

	out, err := os.Open(file)

	if err != nil {
		return "", err
	}

	defer out.Close()

	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
