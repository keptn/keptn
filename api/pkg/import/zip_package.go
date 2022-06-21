package _import

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	logger "github.com/sirupsen/logrus"
)

var /*const*/ ErrorUncompressedSizeTooBig = errors.New("uncompressed size of package exceeds configured maximum size")

// ZippedPackage represents a zipped import package ready to be use (it is extracted to a temp directory)
type ZippedPackage struct {
	extractedDir string
}

// Close signals that the package resources can be freed (including any extracted files).
// Once Close has been called it's illegal to call any other ZippedPackage operation
func (m *ZippedPackage) Close() error {
	if m.extractedDir != "" {
		err := os.RemoveAll(m.extractedDir)
		m.extractedDir = ""
		return err
	}
	return nil
}

func (m *ZippedPackage) extract(zipFile string, maxSize uint64) error {
	srcInfo, err := os.Stat(zipFile)
	if err != nil {
		return fmt.Errorf("could not retrieve file info for %s: %w", zipFile, err)
	}

	srcFile, err := os.Open(zipFile)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", zipFile, err)
	}

	zipReader, err := zip.NewReader(srcFile, srcInfo.Size())
	if err != nil {
		return fmt.Errorf("invalid zip archive: %w", err)
	}

	extractionDir := strings.TrimSuffix(srcFile.Name(), defaultImportArchiveExtension)
	err = os.Mkdir(extractionDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating temporary extraction directory %s: %w", extractionDir, err)
	}

	m.extractedDir = extractionDir
	return extractZipArchive(zipReader, extractionDir, maxSize)
}

// This function could be part of the manifest interface to abstract the physical location of the files
// func (m *ZippedPackage) GetResource(resourcePath string) (io.ReadCloser, error) {
// 	// TODO
// 	return nil, nil
// }

// NewPackage creates a new ZippedPackage object ready to be used.
// The zip file contents will be extracted in a subDirectory with the same name as the file stripped of the .zip
// extension. During the extraction zip file uncompressed content is checked not to surpass maxSize.
// If any error occurs, the temporary folder is cleaned up and nil, error will be returned
func NewPackage(zipFile string, maxSize uint64) (*ZippedPackage, error) {
	m := new(ZippedPackage)
	err := m.extract(zipFile, maxSize)
	if err != nil {
		m.Close()
		return nil, fmt.Errorf("error initializing zip import package: %w", err)
	}
	return m, nil
}

func extractZipArchive(reader *zip.Reader, outputDir string, maxSize uint64) error {

	var extractedSize uint64
	for _, zippedFile := range reader.File {
		logger.Debugf("Extracting file %+v to %s...", reader, outputDir)

		if zippedFile.FileInfo().IsDir() {
			fullOutputDirectoryName := path.Join(outputDir, zippedFile.Name)
			err := os.Mkdir(fullOutputDirectoryName, 0755)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %w", fullOutputDirectoryName, err)
			}
		} else {
			err := func() error {
				if extractedSize+zippedFile.UncompressedSize64 > maxSize {
					return ErrorUncompressedSizeTooBig
				}

				dstFileName, written, err := extractZippedFile(zippedFile, outputDir)

				if err != nil {
					return fmt.Errorf("error extracting %s from archive into %s: %w", zippedFile.Name, dstFileName, err)
				}

				if written != zippedFile.UncompressedSize64 {
					logger.Warnf(
						"Wrong uncompressed size reported for file %s: expected %d, got %d",
						zippedFile.Name, zippedFile.UncompressedSize64, written,
					)
				}

				extractedSize += zippedFile.UncompressedSize64
				return nil
			}()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func extractZippedFile(zippedFile *zip.File, outputDir string) (string, uint64, error) {
	src, err := zippedFile.Open()
	if err != nil {
		return "", 0, fmt.Errorf("error reading %s from archive: %w", zippedFile.Name, err)
	}
	defer src.Close()

	dst, err := os.Create(path.Join(outputDir, zippedFile.Name))
	if err != nil {
		return "", 0, fmt.Errorf(
			"error creating file %s in output directory %s: %w", zippedFile.Name, outputDir, err,
		)
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)

	return dst.Name(), uint64(written), err
}
