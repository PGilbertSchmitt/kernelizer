package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"os"

	knl "kernelizer/kernelate"
)

func main() {
	imageName, kernelName, outputName, err := handleFlags()

	if err != nil {
		log.Fatalln(err.Error())
	}

	img, err := getImageData(imageName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	foo := img.(*image.RGBA)

	k, err := getKernalData(kernelName)
	if err != nil {
		log.Fatalln(err.Error())
	}

	newImg, err := knl.Kernelate(foo, *k)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if outputName == "" {
		outputName = "spooky"
	}

	f, err := os.OpenFile(outputName, os.O_RDWR, 0777)
	if err != nil {
		log.Println(err)

		f, err = os.Create(outputName)
		if err != nil {
			log.Fatalln("You done fucked up now:", err)
		}
	}
	defer f.Close()

	png.Encode(f, newImg)
}

// Returns the name of the input image, the image for the kernel json, and the output image in that order
// TODO: create return object and clean up using flags
func handleFlags() (string, string, string, error) {
	args := os.Args[1:]

	if len(args) < 2 {
		return "", "", "", errors.New("You've been a bad boy. Where are the file names?")
	}

	inputName := args[0]
	kernelName := args[1]
	outputName := ""
	if len(args) > 2 {
		outputName = args[2]
	}

	return inputName, kernelName, outputName, nil
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

	// Sanity check the kernel
	height := len(k.K)

	if height < 1 {
		return nil, errors.New("Kernel must be at least one unit wide")
	}

	if (height % 2) == 0 {
		return nil, errors.New("Kernel must have odd width and height")
	}

	for _, row := range k.K {
		width := len(row)
		if width != height {
			return nil, errors.New("Kernel must be a square")
		}
	}

	return &k, nil
}
