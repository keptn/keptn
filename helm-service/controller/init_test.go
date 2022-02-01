package controller

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Printf("TestMain: error while setting up the tests: %v", err)
		os.Exit(-1)
	}
	code := m.Run()
	os.Exit(code)
}

func setup() error {
	_ = os.Setenv("USE_COMMITID", "true")
	return nil
}
