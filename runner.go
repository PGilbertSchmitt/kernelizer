package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"

	knl "kernelizer/kernelate"
)

func main() {
	imageName, kernelName, err := handleFlags()

	if err != nil {
		log.Fatalln(err.Error())
	}

	img, err := getImageData(imageName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_ = img.(*image.RGBA)

	_, err = getKernalData(kernelName)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// Returns the reader for the image and the image for the kernel json in that order
func handleFlags() (string, string, error) {
	args := os.Args[1:]

	if len(args) < 2 {
		return "", "", errors.New("You've been a bad boy. Where are the file names?")
	}

	imageName := args[0]
	kernelName := args[1]

	return imageName, kernelName, nil
}

func getImageData(imageName string) (image.Image, error) {
	imageFile, err := os.Open(imageName)
	if err != nil {
		log.Fatalln(errors.New("Could not open image file"))
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return img, err
}

func getKernalData(kernelName string) (*knl.Kernel, error) {
	var k knl.Kernel

	kernelFile, err := os.Open(kernelName)
	if err != nil {
		return &k, errors.New("Could not open kernel json")
	}
	defer kernelFile.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, kernelFile)
	decoder := json.NewDecoder(buf)

	decoder.Decode(&k)
	return &k, nil
}
