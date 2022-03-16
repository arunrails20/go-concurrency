package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

// Pipeline
// walkfile ------> processImage -----> SaveImage
//          (paths)              (results)

type result struct {
	srcImagePath   string
	thumbnailImage *image.NRGBA
	err            error
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to send directory path of images")
	}

	start := time.Now()

	err := setupPipeLine(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Time taken: %s\n", time.Since(start))
}

func setupPipeLine(root string) error {
	done := make(chan struct{})
	// terminate all GR in pipeline
	defer close(done)
	// first stage of the pipe
	paths, errc := walkFiles(done, root)
	// second stage of the pipe
	results := processImage(done, paths)

	// third stage of the pipe
	for r := range results {
		if r.err != nil {
			return r.err
		}
		saveThumbnail(r.srcImagePath, r.thumbnailImage)
	}
	// if there was an any error. will return the error
	// than the defer close func will call and all pipeline will terminate

	if err := <-errc; err != nil {
		return err
	}
	return nil
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		// close the paths chan once it sent
		defer close(paths)

		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

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

			// sending path of the img file to path chan
			// we use select stmt to preemption
			select {
			case paths <- path:
			case <-done:
				return fmt.Errorf("walk was canceled")
			}
			return nil
		})

	}()
	return paths, errc

}

func processImage(done <-chan struct{}, paths <-chan string) <-chan *result {

	results := make(chan *result)

	thumbnailer := func() {
		for path := range paths {
			//load the image from file
			srcImage, err := imaging.Open(path)
			if err != nil {
				// Handle the error while open the image file
				select {
				case results <- &result{path, nil, err}:
				case <-done:
					return
				}
			}
			// scale the image to 100px * 100px
			thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

			select {
			case results <- &result{path, thumbnailImage, nil}:
			case <-done:
				return
			}

		}
	}

	const numThumbnailer = 5
	// Use wait groups to sync the goroutines
	var wg sync.WaitGroup

	// we are going to spin 5 goroutines to process the image in parrel
	wg.Add(numThumbnailer)
	for i := 0; i < numThumbnailer; i++ {
		go func() {
			thumbnailer()
			wg.Done()
		}()
	}
	// spin a seperate goroutines to terminates the thumbnailre GR
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
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
