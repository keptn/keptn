package common

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

// IFileWriter is an interface for writing files
//go:generate moq -pkg common_mock -skip-ensure -out ./fake/file_writer_mock.go . IFileWriter
type IFileWriter interface {
	WriteBase64EncodedFile(path string, content string) error
	WriteFile(path string, content []byte) error
	DeleteFile(path string) error
	FileExists(path string) bool
	MakeDir(path string) error
}

type FileWriter struct{}

func (fw FileWriter) WriteBase64EncodedFile(path string, content string) error {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return err
	}
	return fw.WriteFile(path, data)
}

func (fw FileWriter) WriteFile(path string, content []byte) error {
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
		err = fw.DeleteFile(path)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	path = filepath.Clean(path)
	file, err = os.OpenFile(path, os.O_RDWR, 0600)
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
	return err
}

func (FileWriter) DeleteFile(path string) error {
	var err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func (FileWriter) FileExists(path string) bool {
	_, err := os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (fw FileWriter) MakeDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
