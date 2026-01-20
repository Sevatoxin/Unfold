package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
)

var (
	input_path           string        // This contains the filepath of the folder that got dragged onto the executable
	main_subfolders      []os.DirEntry // This will contain all subfolders in the input_path to be deleted in the end
	all_subfolders_paths []string      // This contains all paths to all subfolders in the input_path folder
)

// Checks error and panics the program to stop it
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Filtering a folder for all subfolders and drag the files to the main folder from the input_path
func FilterFolder() {
	OpenMainFolder()
	CheckForFolders()
}

// Change the current directory
func OpenMainFolder() {
	err := os.Chdir(input_path)
	check(err)
}

// Read the content of the main folder
func CheckForFolders() {

	GetMainSubfolders() // to delete them in the end
	GetAllSubFolders()  // To be able to grab all sub files and put them into the root
	MoveAllFiles()
	DeleteSubfolders()
}

func DeleteSubfolders() {
	for _, dir := range main_subfolders {
		os.RemoveAll(path.Join(input_path, dir.Name()))
	}
}

func MoveAllFiles() {
	for _, folder := range all_subfolders_paths {
		os.Chdir(path.Join(input_path, folder))
		fmt.Println("Currently in folder: ", folder)

		entries, err := os.ReadDir(".")
		for _, container := range entries {
			fmt.Println("It contains: ", container.Name())
		}
		check(err)

		for _, entry := range entries {
			if !entry.IsDir() {
				fmt.Println("Moving: ", entry.Name())
				MoveFile(entry)
			}
		}
	}

	os.Chdir(input_path)
}

func GetAllSubFolders() {
	err := filepath.WalkDir(".", func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() && path != input_path && path != "." {
			println("Adding following path: ", path)
			all_subfolders_paths = append(all_subfolders_paths, path)
		}

		return nil
	})

	check(err)
}

func GetMainSubfolders() {
	entries, err := os.ReadDir(".")
	check(err)

	// Get all main subfolders
	for _, entry := range entries {
		if entry.IsDir() {
			main_subfolders = append(main_subfolders, entry)
		}
	}
}

func MoveFile(file os.DirEntry) {
	var new_path string
	old_path := path.Join(".", file.Name())

	// Check file already exists in input_path
	entries, err := os.ReadDir(input_path)
	check(err)
	for _, entry := range entries {
		if !entry.IsDir() {
			if file.Name() == entry.Name() {
				new_path = path.Join(input_path, "WARNING_CHECK_FILE_"+file.Name())
			} else {
				new_path = path.Join(input_path, file.Name())
			}
		}
	}

	os.Rename(old_path, new_path)
}

func main() {
	start_time := time.Now()
	args := os.Args[1:]

	// Then start program if there was a folder input
	if len(args) > 0 {
		input_path = args[0]

		FilterFolder()

		fmt.Println("\nThe filtering took: ", time.Since(start_time))
		fmt.Println("\n\nSuccessfully filtered all folders. Press Enter to close this program...")
	} else {
		fmt.Println("To use this program please drag a folder onto the executable named unfold.exe")
	}

	fmt.Scanf("c") // The keep the program open
}
