package common

import (
	"os"
)

// WriteFile writes to a file in the filesystem
func WriteFile(path string, content []byte) error {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
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
