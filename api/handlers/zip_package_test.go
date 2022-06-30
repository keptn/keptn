package handlers

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractZipFileHappyPath(t *testing.T) {

	sourceImportPackage := "../test/data/import/sample-package"

	tempDir, err := ioutil.TempDir("", "test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tempZipFile, err := ioutil.TempFile(tempDir, "test-archive*"+defaultImportArchiveExtension)
	require.NoError(t, err)

	err = writeZip(tempZipFile, sourceImportPackage)
	require.NoError(t, err)

	err = tempZipFile.Close()
	require.NoError(t, err)

	p, err := NewZippedPackage(tempZipFile.Name(), testArchiveSize20MB)
	require.NoError(t, err)
	require.NotNil(t, p)

	// assert that creating a package from a zipped file extracted the files in a subdir with the same name as the
	// file minus the extension
	expectedExtractedPath := strings.TrimSuffix(tempZipFile.Name(), defaultImportArchiveExtension)
	assert.DirExists(t, expectedExtractedPath)
	assertDirEqual(t, sourceImportPackage, expectedExtractedPath)

	// assert that closing the package cleans up the extracted files
	err = p.Close()
	assert.NoError(t, err)
	assert.NoDirExists(t, expectedExtractedPath)
}

func TestExtractErrorZipFilePackageTooBig(t *testing.T) {

	sourceImportPackage := "../test/data/import/sample-package"

	tempDir, err := ioutil.TempDir("", "test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tempZipFile, err := ioutil.TempFile(tempDir, "test-archive*"+defaultImportArchiveExtension)
	require.NoError(t, err)

	err = writeZip(tempZipFile, sourceImportPackage)
	require.NoError(t, err)

	err = tempZipFile.Close()
	require.NoError(t, err)

	p, err := NewZippedPackage(tempZipFile.Name(), 10)
	assert.ErrorIs(t, err, ErrorUncompressedSizeTooBig)
	assert.Nil(t, p)

	// the extraction folder is cleaned up
	expectedExtractedPath := strings.TrimSuffix(tempZipFile.Name(), defaultImportArchiveExtension)
	assert.NoDirExists(t, expectedExtractedPath)
}

func TestExtractErrorNonExistentZipFile(t *testing.T) {

	tempDir, err := ioutil.TempDir("", "test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	nonExistingZipFileName := "thereisnospoon.zip"

	p, err := NewZippedPackage(path.Join(tempDir, nonExistingZipFileName), testArchiveSize20MB)
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestErrorInvalidZipFile(t *testing.T) {

	tempDir, err := ioutil.TempDir("", "test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	invalidZipFile := path.Join(tempDir, "invalid.zip")
	ioutil.WriteFile(invalidZipFile, []byte("this is clearly not a zip file"), 0600)

	p, err := NewZippedPackage(invalidZipFile, testArchiveSize20MB)
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestExtractErrorNoManifest(t *testing.T) {

	sourceImportPackage := "../test/data/import/invalid-package"

	tempDir, err := ioutil.TempDir("", "test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tempZipFile, err := ioutil.TempFile(tempDir, "test-archive*"+defaultImportArchiveExtension)
	require.NoError(t, err)

	err = writeZip(tempZipFile, sourceImportPackage)
	require.NoError(t, err)

	err = tempZipFile.Close()
	require.NoError(t, err)

	p, err := NewZippedPackage(tempZipFile.Name(), testArchiveSize20MB)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Nil(t, p)

	// the extraction folder is cleaned up
	expectedExtractedPath := strings.TrimSuffix(tempZipFile.Name(), defaultImportArchiveExtension)
	assert.NoDirExists(t, expectedExtractedPath)
}

func assertDirEqual(t *testing.T, expected string, actual string) {
	filepath.WalkDir(
		expected, func(walkedPath string, d os.DirEntry, err error) error {

			relPath, err := filepath.Rel(expected, walkedPath)
			require.NoError(t, err)

			actualPath := path.Join(actual, relPath)
			if d.IsDir() {
				assert.DirExists(t, actualPath)
			} else {
				assert.FileExists(t, actualPath)
				actualFileBytes, err := ioutil.ReadFile(actualPath)
				assert.NoError(t, err)
				actualFileHash := sha256.Sum256(actualFileBytes)

				expectedFileBytes, err := ioutil.ReadFile(walkedPath)
				assert.NoError(t, err)
				expectedFileHash := sha256.Sum256(expectedFileBytes)

				assert.Equalf(
					t, hex.EncodeToString(expectedFileHash[:]), hex.EncodeToString(actualFileHash[:]),
					"files %s and %s should have the same content but their SHA256 is different!", walkedPath,
					actualPath,
				)
			}

			return nil
		},
	)
}

func writeZip(file *os.File, filePaths ...string) error {
	zipWriter := zip.NewWriter(file)
	for _, filePath := range filePaths {
		fileInfo, err := os.Lstat(filePath)
		if err != nil {
			return fmt.Errorf("error getting info on %s: %w", filePath, err)
		}
		if fileInfo.IsDir() {
			err = addDirectoryContentToZip(zipWriter, filePath)
		} else {
			err = addFileToZip(zipWriter, filePath, path.Base(filePath))
		}

		if err != nil {
			return fmt.Errorf("error adding %s to zip: %w", filePath, err)
		}
	}

	return zipWriter.Close()

}

func addDirectoryToZip(writer *zip.Writer, dirName string) error {
	if !strings.HasSuffix(dirName, "/") {
		dirName = dirName + "/"
	}
	_, err := writer.Create(dirName)
	if err != nil {
		return fmt.Errorf("error creating directory %s in zip file: %w", dirName, err)
	}
	return nil
}

func addFileToZip(writer *zip.Writer, srcPath string, dstPath string) error {
	fileWriter, err := writer.Create(dstPath)
	if err != nil {
		return fmt.Errorf("error creating file %s in zip file: %w", dstPath, err)
	}
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %w", srcPath, err)
	}

	defer src.Close()
	_, err = io.Copy(fileWriter, src)
	if err != nil {
		return fmt.Errorf(
			"error writing file content from %s to %s (in zip file): %w",
			srcPath, dstPath, err,
		)
	}
	return nil
}

func addDirectoryContentToZip(writer *zip.Writer, baseDir string) error {
	return filepath.WalkDir(
		baseDir, func(walkedPath string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			relativePath, err := filepath.Rel(baseDir, walkedPath)

			if err != nil {
				return fmt.Errorf("error calculating path %s as relative path of %s", walkedPath, baseDir)
			}

			if relativePath == "." {
				return nil
			}

			if d.IsDir() {
				return addDirectoryToZip(writer, relativePath)
			}

			return addFileToZip(writer, walkedPath, relativePath)
		},
	)
}
