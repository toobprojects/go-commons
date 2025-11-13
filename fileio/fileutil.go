package fileio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/toobprojects/go-commons/logs"
)

// ReadContent
// Used to read file content, this will also convert the read bytes to String.
func ReadContent(pathArg string) string {
	if isSymbolicLink(pathArg) {
		pathArg, _ = os.Readlink(pathArg)
	}

	content, err := os.ReadFile(pathArg)
	if err != nil {
		logs.Error(fmt.Sprintf("Could not open file : %v", err))
	}

	return string(content)
}

// Exists
// Checks if a given file Exists on the file
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NotExists(path string) (bool, error) {
	exists, err := Exists(path)
	return !exists, err
}

// List
// All the fileio in the given directory
func List(directory string) []os.FileInfo {

	// Try open the directory
	file, err := os.Open(directory)
	if err != nil {
		logs.Error(fmt.Sprintf("Could not list directory %s: %v", directory, err))
		return nil
	}

	defer file.Close()

	// If there's no error when opening then you can now read the directory
	// This will give you a list of FileInfo records that you can use to see info of each file
	list, err := file.Readdir(-1)
	if err != nil {
		logs.Error(fmt.Sprintf("Could not list directory %s: %v", directory, err))
		return nil
	}

	return list
}

func ListByWildcard(directory string, suffix string) []string {
	matches, err := filepath.Glob(filepath.Join(directory, suffix))
	if err != nil {
		logs.Error(fmt.Sprintf("Could not list directory : %v", directory))
	}

	return matches
}

// RemoveAllFromDirectory
// Remove all the fileio in the given directory,
// This takes advantage of the List function to list first and then remove.
func RemoveAllFromDirectory(directory string) {
	files := List(directory)
	for _, file := range files {
		fullLogFilePath := directory + "/" + file.Name()
		err := os.Remove(fullLogFilePath)
		if err != nil {
			logs.Error(fmt.Sprintf("Could not delete file : %v", err))
		}
	}
}

// RemoveAllFromDirectoryRecursively
// Remove all the fileio in the given directory,
// This takes advantage of the List function to list first and then remove.
// We also remove recursively, meaning when we have sub folders that are not empty
func RemoveAllFromDirectoryRecursively(directory string) error {
	// Walk through the directory tree
	return filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", path, err)
		}
		if path == directory {
			return nil // Skip the root directory itself
		}

		// Remove the file or directory
		err = os.RemoveAll(path)
		if err != nil {
			return fmt.Errorf("failed to remove %s: %v", path, err)
		}
		return nil
	})
}

func ExtractFolderNameFromPath(path string) string {
	return filepath.Base(path)
}

// Move
// Takes in two Absolute path of the source and destination fileio.
// It uses these to Copy the fileio over from one directory to anoother.
func Move(source string, destination string) {
	err := os.Rename(source, destination)
	if err != nil {
		logs.Error(fmt.Sprintf("Could not move file from %s to %s: %v", source, destination, err))
	}
}

// Copy
// Takes in two Absolute paths for the source and destination fileio.
// These are used to copy from one file to another.
func Copy(source string, destination string) {

	// Open the file now and get its contents.
	sourceFile, err := os.Open(source)
	if err != nil {
		logs.Error(fmt.Sprintf("Could not open source file %s: %v", source, err))
	}

	// Defer the CLOSE until we are done, and we can do this as the last line in the function before the return
	defer sourceFile.Close()

	// Create a new file in the destination path and get it ready for receiving contents.
	destinationFile, err := os.Create(destination)
	if err != nil {
		logs.Error(fmt.Sprintf("Could not create destination file %s: %v", destination, err))
	}

	// Defer the CLOSE until we are done, and we can do this as the last line in the function before the return
	defer destinationFile.Close()

	// Let's start the copying from source to destination
	_, err = io.Copy(destinationFile, sourceFile)

	if err != nil {
		logs.Error(fmt.Sprintf("Error copying from %s to %s: %v", source, destination, err))
	}
}

func Pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		logs.Error(fmt.Sprintf("Error getting current working directory: %v", err))
	}
	return pwd
}

// HomePath
// Used to get the Home Path on a Machine.
func HomePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		logs.Error(fmt.Sprintf("Error getting home directory: %v", err))
	}

	return home
}

// isSymbolicLink
// Checks if File is actual file of Symbolic Link to a File Somewhere else.
func isSymbolicLink(path string) bool {
	// Get file information
	fileInfo, err := os.Lstat(path)
	if err != nil {
		logs.Error(fmt.Sprintf("Error getting file information for path %s: %v", path, err))
	}

	// Check if the file is a symbolic link
	return fileInfo.Mode()&os.ModeSymlink != 0
}
