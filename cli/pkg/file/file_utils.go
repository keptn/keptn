package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

// PlaceholderReplacement is a helper type for replacing a placeholder with a desired value
type PlaceholderReplacement struct {
	PlaceholderValue string
	DesiredValue     string
}

// FileExists checks whether a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(keptnutils.ExpandTilde(filename))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ReadFile reads a file and returns the content as string
func ReadFile(fileName string) (string, error) {

	fileName = keptnutils.ExpandTilde(fileName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return "", fmt.Errorf("Cannot find file %s", fileName)
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DownloadFile will download a url to a local file.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Replace reads a file, replaces all placeholders and writes the file back
func Replace(filePath string, replacements ...PlaceholderReplacement) error {
	content, err := ReadFile(filePath)
	if err != nil {
		return err
	}
	for _, replacement := range replacements {
		content = strings.ReplaceAll(content, replacement.PlaceholderValue, replacement.DesiredValue)
	}

	return ioutil.WriteFile(filePath, []byte(content), 0666)
}

// Walks and Refactors Codebase
func RecursiveRefactor(old, new string, filePatterns []string) filepath.WalkFunc {
	return filepath.WalkFunc(func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !!fi.IsDir() {
			return nil
		}

		var matched bool
		for _, pattern := range filePatterns {
			var err error
			matched, err = filepath.Match(pattern, fi.Name())
			if err != nil {
				return err
			}

			if matched {
				read, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				newContents := strings.Replace(string(read), old, new, -1)
				err = ioutil.WriteFile(path, []byte(newContents), 0)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
