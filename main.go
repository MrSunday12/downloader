package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	images := []string{"nginx:latest", "squidfunk/mkdocs-material:8.1.2"}

	log.Printf("Images are: %v", images)

	// Creating folder to save images to
	folderPath := os.ExpandEnv("$HOME/images")
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	existingImages, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		log.Println(strings.Repeat("-", 50))
		pulled := false

		log.Printf("Checking if image is already pulled: %q", image)
		for _, existingImage := range existingImages {
			if ok := contains(existingImage.RepoTags, image); ok {
				pulled = true
				break
			}
		}

		if !pulled {
			log.Printf("Image not pulled, pulling it: %q", image)

			out, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
			if err != nil {
				panic(err)
			}

			if _, err := io.Copy(io.Discard, out); err != nil {
				panic(err)
			}

			out.Close()
			log.Printf("Successfully pulled image: %q", image)
		} else {
			log.Printf("Image already pulled: %q", image)
		}

		// Updating the existingImages list after pulling the image
		existingImages, err := cli.ImageList(context.Background(), types.ImageListOptions{})
		if err != nil {
			panic(err)
		}

		imageId := ""
		for _, existingImage := range existingImages {
			if ok := contains(existingImage.RepoTags, image); ok {
				imageId = existingImage.ID
				break
			}
		}
		log.Printf("Image ID is: %q", imageId)

		imagePath := os.ExpandEnv(fmt.Sprintf("%s/%s.docker", folderPath, strings.ReplaceAll(strings.ReplaceAll(image, ":", "."), "/", ".")))
		log.Printf("Saving image: %q in: %q", image, imagePath)

		os.Remove(imagePath)

		out, err := cli.ImageSave(context.Background(), []string{imageId})
		if err != nil {
			panic(err)
		}

		outfile, err := os.Create(imagePath)
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(outfile, out); err != nil {
			panic(err)
		}

		out.Close()
		outfile.Close()
	}
}

// Check whether a slice contains the given string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
