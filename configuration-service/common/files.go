package common

import (
	"encoding/base64"
	"os"
	"strings"
)

// WriteBase64EncodedFile writes a file using the base64 encoded input
func WriteBase64EncodedFile(path string, content string) error {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return err
	}
	return WriteFile(path, data)
}

// WriteFile writes to a file in the filesystem if it exists, it is overwritten
func WriteFile(path string, content []byte) error {
	pathArr := strings.Split(path, "/")
	directory := ""
	for _, pathItem := range pathArr[0 : len(pathArr)-1] {
		directory += pathItem + "/"
	}

	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}
	// detect if file exists
	_, err = os.Stat(path)

	// delete the file and re-create it, if it existed previously
	if !os.IsNotExist(err) {
		DeleteFile(path)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(string(content))
	if err != nil {
		return err
	}

	// save changes
	err = file.Sync()
	return nil
}

// DeleteFile deletes a file
func DeleteFile(path string) error {
	var err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

// FileExists checks wether a file is available or not
func FileExists(path string) bool {
	_, err := os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		return false
	}
	return true
}
