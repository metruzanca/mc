package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"

	"github.com/spf13/cobra"
)

const padLen = 3

var counter = 0
var prefix string

func padCounter(counter int) string {
	return fmt.Sprintf("%0*d", padLen, counter)
}

func renameFile(oldPath string) {
	dir := filepath.Dir(oldPath)
	ext := filepath.Ext(oldPath)
	base := filepath.Base(oldPath)
	base = strings.TrimSuffix(base, ext)

	if strings.HasPrefix(base, prefix) {
		return
	}

	newBase := prefix + padCounter(counter)
	newPath := filepath.Join(dir, newBase+ext)

	err := os.Rename(oldPath, newPath)
	if err != nil {
		log.Printf("Failed to rename file %s: %s\n", oldPath, err)
	} else {
		log.Printf("Renamed %s to %s\n", oldPath, newPath)
		counter++
	}
}

// prefixCmd represents the prefix command
var prefixCmd = &cobra.Command{
	Use:   "prefix [path] <prefix>",
	Short: "Watches current directory and renames new files to have prefix and counter",
	Long:  `A longer description that spans multiple lines`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("Prefix command usage: [path] <prefix>")
			os.Exit(1)
		}

		var path string
		if len(args) == 1 {
			prefix = args[0]
			currentDir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			path = currentDir
		} else {
			path = args[0]
			prefix = args[1]
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		done := make(chan bool)

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Create == fsnotify.Create {
						// time.Sleep(1 * time.Second) // Ensure file is fully written
						renameFile(event.Name)
					}
				case err := <-watcher.Errors:
					log.Println("error:", err)
				}
			}
		}()

		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Watching directory: %s\n", path)
		<-done
	},
}

func init() {
	rootCmd.AddCommand(prefixCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prefixCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prefixCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
