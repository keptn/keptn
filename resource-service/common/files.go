package common

import (
	"encoding/base64"
	"errors"
	"fmt"
	archive "github.com/mholt/archiver/v3"
	"github.com/otiai10/copy"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// IFileSystem is an interface for writing files
//go:generate moq -pkg common_mock -skip-ensure -out ./fake/file_writer_mock.go . IFileSystem
type IFileSystem interface {
	WriteBase64EncodedFile(path string, content string) error
	WriteHelmChart(path string) error
	WriteFile(path string, content []byte) error
	ReadFile(filename string) ([]byte, error)
	DeleteFile(path string) error
	FileExists(path string) bool
	MakeDir(path string) error
	WalkPath(path string, walkFunc filepath.WalkFunc) error
}

type FileSystem struct {
	tmpDirLocation string
}

func NewFileSystem(tmpDirLocation string) *FileSystem {
	return &FileSystem{tmpDirLocation: tmpDirLocation}
}

func (fw FileSystem) WriteBase64EncodedFile(path string, content string) error {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return err
	}
	return fw.WriteFile(path, data)
}

func (fw FileSystem) WriteFile(path string, content []byte) error {
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

func (fw FileSystem) WriteHelmChart(path string) error {
	// remove previous helm/resourceURI folder
	targetFolderPath := strings.TrimSuffix(path, ".tgz")
	if err := os.RemoveAll(targetFolderPath); err != nil {
		return fmt.Errorf("could not delete existing folder %s, %v", targetFolderPath, err)
	}
	if err := fw.untarHelm(path); err != nil {
		return err
	}
	return nil
}

func (fw FileSystem) ReadFile(filename string) ([]byte, error) {
	filename = filepath.Clean(filename)
	return ioutil.ReadFile(filename)
}

func (FileSystem) DeleteFile(path string) error {
	var err = os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

func (FileSystem) FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (fw FileSystem) MakeDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func (fw FileSystem) WalkPath(path string, walkFunc filepath.WalkFunc) error {
	return filepath.Walk(path, walkFunc)
}

func (fw FileSystem) untarHelm(filePath string) error {
	tmpDir, err := ioutil.TempDir(fw.tmpDirLocation, "*")
	if err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			logger.WithError(err).Errorf("Could not remove directory %s", tmpDir)
		}
	}()

	tarGz := archive.NewTarGz()
	tarGz.OverwriteExisting = true
	if err := tarGz.Unarchive(filePath, tmpDir); err != nil {
		return fmt.Errorf("could not unarchive Helm chart: %w", err)
	}

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("could not read unpacked files: %w", err)
	}

	if len(files) != 1 {
		return errors.New("unexpected amount of unpacked files")
	}

	folderName := filepath.Join(tmpDir, filePath[strings.LastIndex(filePath, "/")+1:len(filePath)-4])
	oldPath := filepath.Join(tmpDir, files[0].Name())
	if oldPath != folderName {
		if err := os.Rename(oldPath, folderName); err != nil {
			return fmt.Errorf("could not rename unpacked folder: %w", err)
		}
	}

	dir, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		return fmt.Errorf("patch of helm chart is invalid: %w", err)
	}

	if err := copy.Copy(tmpDir, dir); err != nil {
		return fmt.Errorf("could not copy folder: %w", err)
	}

	// remove Helm chart .tgz file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete helm chart package: %w", err)
	}
	return nil
}

func IsHelmChartPath(resourcePath string) bool {
	resourcePathSlice := strings.Split(resourcePath, "/")
	if sliceLen := len(resourcePathSlice); sliceLen >= 2 {
		// return true if the resource path ends with "helm/<resourceName>.tgz
		return resourcePathSlice[sliceLen-2] == "helm" && strings.HasSuffix(resourcePathSlice[sliceLen-1], ".tgz")
	}
	return false
}
