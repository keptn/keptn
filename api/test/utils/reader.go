package utils

import (
	"errors"
	"io"
)

// TestReader implements io.Reader interface and allows to generate data and/or throw errors for testing purposes.
type TestReader struct {
	data         []byte
	allowedLoops int
	loopCount    int
	pos          int
	throwError   bool
}

// NewTestReader will return a newly allocated TestReader which will return the data,
// looping over it for a certain number of times before optionally throwing an error or returning io.EOF
func NewTestReader(data []byte, allowedLoops int, throwError bool) *TestReader {
	return &TestReader{
		data:         data,
		allowedLoops: allowedLoops,
		throwError:   throwError,
	}
}

// Read implements io.Reader semantics: will return data if still allowed by configuration,
// return io.EOF (if throwError == false) or a custom error (
// throwError == true) once all the configured data has been returned
func (tr *TestReader) Read(p []byte) (n int, err error) {
	// check if we can reset
	if tr.pos >= len(tr.data) && tr.loopCount < tr.allowedLoops {
		tr.allowedLoops++
		tr.pos = 0
	}

	if tr.pos >= len(tr.data) {
		err := io.EOF

		if tr.throwError {
			err = errors.New("TestReader error")
		}

		// end of buffer
		return 0, err
	}

	// try to copy as much as possible in the buffer
	copied := copy(p, tr.data[tr.pos:])
	tr.pos += copied
	return copied, nil
}
