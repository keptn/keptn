package auth

import (
	"fmt"
	"os/exec"
	"runtime"
)

type URLOpener interface {
	Open(url string) error
}

func NewBrowser() *Browser {
	return &Browser{}
}

type Browser struct {
}

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

type BrowserMock struct {
	openFn func(string) error
}

func (b *BrowserMock) Open(url string) error {
	if b != nil && b.openFn != nil {
		return b.openFn(url)
	}
	return nil
}
