package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
  "strings"
)

func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func mergeFiles(destPath, srcPath string) error {
	destFile, err := os.OpenFile(destPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func removeDuplicateLines(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	uniqueLines := make(map[string]bool)
	var output []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !uniqueLines[line] {
			uniqueLines[line] = true
			output = append(output, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	tmpFilePath := filePath + ".tmp"
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	for _, line := range output {
		if _, err := tmpFile.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	if err := os.Rename(tmpFilePath, filePath); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory>")
		return
	}

	dir := os.Args[1]
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	fileHashes := make(map[string]string)
	duplicateFiles := make(map[string][]string)
	
	prefixMap := make(map[string][]string)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		fileHash, err := hashFile(filePath)
		if err != nil {
			fmt.Println("Error hashing file:", err)
			continue
		}		

		if existingFile, exists := fileHashes[fileHash]; exists {
			duplicateFiles[existingFile] = append(duplicateFiles[existingFile], filePath)
		} else {
			fileHashes[fileHash] = filePath
		}

		fileName := file.Name()
		prefix := strings.Split(fileName, "_")[0]
		prefixMap[prefix] = append(prefixMap[prefix], fileName)
	}

	for originalFile, dupes := range duplicateFiles {
		for _, dupe := range dupes {
			fmt.Printf("Merging %s into %s\n", dupe, originalFile)
			if err := mergeFiles(originalFile, dupe); err != nil {
				fmt.Println("Error merging files:", err)
				continue
			}
			fmt.Printf("Removing duplicate file %s\n", dupe)
			if err := os.Remove(dupe); err != nil {
				fmt.Println("Error removing file:", err)
			}
		}
		fmt.Printf("Removing duplicate lines from %s\n", originalFile)
		if err := removeDuplicateLines(originalFile); err != nil {
			fmt.Println("Error removing duplicate lines:", err)
		}
	}

	for prefix, fileList := range prefixMap {
		if len(fileList) > 1 {
			var mergedFile string
			for i, file := range fileList {
				filePath := filepath.Join(dir, file)
				if i == 0 {
					mergedFile = filePath
					continue
				}
                fmt.Println("aaa",prefix)
				fileHash, err := hashFile(filePath)
				if err != nil {
					fmt.Println("Error hashing file:", err)
					continue
				}

				mergedFileHash, err := hashFile(mergedFile)
				if err != nil {
					fmt.Println("Error hashing merged file:", err)
					continue
				}

				if fileHash == mergedFileHash {
					fmt.Printf("Removing duplicate file %s\n", filePath)
					if err := os.Remove(filePath); err != nil {
						fmt.Println("Error removing file:", err)
					}
				} else {
					fmt.Printf("Merging %s into %s\n", filePath, mergedFile)
					if err := mergeFiles(mergedFile, filePath); err != nil {
						fmt.Println("Error merging files:", err)
						continue
					}
					fmt.Printf("Removing file %s after merging\n", filePath)
					if err := os.Remove(filePath); err != nil {
						fmt.Println("Error removing file:", err)
					}
					fmt.Printf("Removing duplicate lines from %s\n", mergedFile)
					if err := removeDuplicateLines(mergedFile); err != nil {
						fmt.Println("Error removing duplicate lines:", err)
					}
				}
			}
		}
    }
	fmt.Println("Duplicate file cleanup complete.")
}



