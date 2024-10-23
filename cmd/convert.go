package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert <file...> <type>",
	Short: "Convert images from one format to another",
	Args:  cobra.MinimumNArgs(2), // Accept at least two arguments
	Run: func(cmd *cobra.Command, args []string) {
		targetFormat := strings.ToLower(args[len(args)-1]) // The last argument is the target format
		inputFiles := args[:len(args)-1]                   // All but the last argument are input files

		for _, file := range inputFiles {
			convertImage(file, targetFormat)
		}
	},
}

func convertImage(file string, targetFormat string) {
	buffer, err := bimg.Read(file)
	if err != nil {
		log.Printf("Failed to read image %s: %v", file, err)
		return
	}

	// Determine the image type from the buffer
	imageType := bimg.DetermineImageType(buffer)
	if imageType == bimg.ImageType(0) {
		log.Printf("Failed to determine image type for %s: unsupported format", file)
		return
	}

	// Set the target type based on the provided format
	options := bimg.Options{Type: imageType}

	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		log.Printf("Failed to process image %s: %v", file, err)
		return
	}

	outputFile := strings.TrimSuffix(file, filepath.Ext(file)) + "." + targetFormat
	if err := bimg.Write(outputFile, newImage); err != nil {
		log.Printf("Failed to write image %s: %v", outputFile, err)
		return
	}

	fmt.Printf("Converted %s to %s\n", file, outputFile)
}

func init() {
	imageCmd.AddCommand(convertCmd)
}
