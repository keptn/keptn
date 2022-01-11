package auth

import (
	"fmt"
	"os/exec"
	"runtime"
)

// URLOpener opens a given URL
type URLOpener interface {
	Open(url string) error
}

// NewBrowser creates a new Browser
func NewBrowser() *Browser {
	return &Browser{}
}

// Browser is an implementation of URLOpener which opens an URL
// using the local Browser
type Browser struct {
}

// Open opens the gven URL using the local browser
// This method is platform independent, thus works on
// Windows, Linux as well as OSX
func (b Browser) Open(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

// BrwoserMock is an implementation of URLOpener
// usable mocking in tests
type BrowserMock struct {
	openFn func(string) error
}

// Open calls the mocked function of the BrwoserMock
func (b *BrowserMock) Open(url string) error {
	if b != nil && b.openFn != nil {
		return b.openFn(url)
	}
	return nil
}
