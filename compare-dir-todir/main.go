package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Parse command line arguments
	args := os.Args
	if len(args) != 4 {
		fmt.Println("Usage: dir_compare <dir1> <dir2> <output_dir>")
		return
	}
	dir1 := args[1]
	dir2 := args[2]
	outputDir := args[3]
	//make sure outputDir does not in dir1 or dir2
	existsInDir1 := dirExistsInDir(dir1, outputDir)
	existsInDir2 := dirExistsInDir(dir2, outputDir)
	if existsInDir1 || existsInDir2 {
		fmt.Println("Output Dir can be in compared Dir!")
		return
	}

	// Get list of files in dir1 and dir2
	files1, err := getFileList(dir1)
	if err != nil {
		fmt.Println("Failed to get file list for", dir1)
		return
	}
	files2, err := getFileList(dir2)
	if err != nil {
		fmt.Println("Failed to get file list for", dir2)
		return
	}

	// Compare MD5 checksums of all files
	uniqueFiles := make(map[string]string)
	for _, file1 := range files1 {
		file2, exists := files2[filepath.Base(file1)]
		if exists {
			if md5Checksum(file1) != md5Checksum(file2) {
				fmt.Println("MD5 checksums do not match for file", file1)
				uniqueFiles[file1] = ""
				uniqueFiles[file2] = ""
			}
		} else {
			fmt.Println("File", filepath.Base(file1), "does not exist in", dir2)
			uniqueFiles[file1] = ""
		}
	}
	for file2 := range files2 {
		if _, exists := files1[filepath.Base(file2)]; !exists {
			fmt.Println("File", filepath.Base(file2), "does not exist in", dir1)
			uniqueFiles[file2] = ""
		}
	}

	// Copy all unique files to output directory
	for file := range uniqueFiles {
		destFile := filepath.Join(outputDir, filepath.Base(file))
		err := copyFile(file, destFile)
		if err != nil {
			fmt.Println("Failed to copy file", file)
		}
	}
}

// Returns a map of file names to file paths in the given directory
func getFileList(dir string) (map[string]string, error) {
	fileList := make(map[string]string)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			fileList[file.Name()] = filepath.Join(dir, file.Name())
		}
	}
	return fileList, nil
}

// Calculates the MD5 checksum of the given file
func md5Checksum(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// Copies a file from src to dest
func copyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return err
	}
	if err = destFile.Sync(); err != nil {
		return err
	}
	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dest, srcFileInfo.Mode())
}

// dirExistsInDir
//
//	@Description: check if a dir in parent dir
//	@param parentDir
//	@param childDir
//	@return bool
func dirExistsInDir(parentDir string, childDir string) bool {
	fileInfo, err := os.Stat(parentDir)
	if err != nil {
		return false
	}

	childPath := filepath.Join(parentDir, childDir)
	if _, err := os.Stat(childPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return fileInfo.IsDir() && strings.HasPrefix(childPath, parentDir)
}
