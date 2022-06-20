package _import

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	logger "github.com/sirupsen/logrus"
)

type ZipArchive struct {
	extractedDir string
}

func (m *ZipArchive) Close() error {
	if m.extractedDir != "" {
		return os.RemoveAll(m.extractedDir)
	}
	return nil
}

func (m *ZipArchive) extract(zipFile string) error {
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
	return extractZipFile(zipReader, extractionDir)
}

// This function could be part of the manifest interface to abstract the phisycal location of the files
// func (m *ZipArchive) GetResource(resourcePath string) (io.ReadCloser, error) {
// 	// TODO
// 	return nil, nil
// }

func NewManifest(zipFile string) (*ZipArchive, error) {
	m := new(ZipArchive)
	err := m.extract(zipFile)
	if err != nil {
		return nil, fmt.Errorf("error initializing zip import package: %w", err)
	}
	return m, nil
}

func extractZipFile(reader *zip.Reader, outputDir string) error {

	// TODO add constraint on maximum extracted archive size

	for _, zippedFile := range reader.File {
		logger.Debugf("Extracting file %+v to %s...", reader, outputDir)

		if zippedFile.FileInfo().IsDir() {
			fullOuputDirectoryName := path.Join(outputDir, zippedFile.Name)
			err := os.Mkdir(fullOuputDirectoryName, 0755)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %w", fullOuputDirectoryName, err)
			}
		} else {
			err := func() error {
				src, err := zippedFile.Open()
				if err != nil {
					return fmt.Errorf("error reading %s from archive: %w", zippedFile.Name, err)
				}
				defer src.Close()

				dst, err := os.Create(path.Join(outputDir, zippedFile.Name))
				if err != nil {
					return fmt.Errorf(
						"error creating file %s in output directory %s: %w", zippedFile.Name, outputDir, err,
					)
				}
				defer dst.Close()

				// TODO add safety checks on the actual written file size
				_, err = io.Copy(dst, src)

				if err != nil {
					return fmt.Errorf("error extracting %s from archive into %s: %w", zippedFile.Name, dst.Name(), err)
				}

				return nil
			}()

			if err != nil {
				return err
			}
		}
	}

	return nil
}
