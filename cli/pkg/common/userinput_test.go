package common

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_AssumeTrue(t *testing.T) {
	opts := UserInputOptions{
		AssumeYes: true,
	}
	answer := NewUserInput().AskBool("hello are you there?", &opts)
	assert.True(t, answer)
}

func Test_UserInput(t *testing.T) {
	opts := UserInputOptions{}

	t.Run("test Y", func(t *testing.T) {
		input := UserInput{
			Writer: ioutil.Discard,
			Reader: bufio.NewReader(bytes.NewBufferString("Y\n")),
		}
		answer := input.AskBool("Are you sure?", &opts)
		assert.True(t, answer)
	})
	t.Run("test y", func(t *testing.T) {
		input := UserInput{
			Writer: ioutil.Discard,
			Reader: bufio.NewReader(bytes.NewBufferString("y\n")),
		}
		answer := input.AskBool("Are you sure?", &opts)
		assert.True(t, answer)
	})
	t.Run("test Yes", func(t *testing.T) {
		input := UserInput{
			Writer: ioutil.Discard,
			Reader: bufio.NewReader(bytes.NewBufferString("Yes\n")),
		}
		answer := input.AskBool("Are you sure?", &opts)
		assert.True(t, answer)
	})
	t.Run("test yes", func(t *testing.T) {
		input := UserInput{
			Writer: ioutil.Discard,
			Reader: bufio.NewReader(bytes.NewBufferString("yes\n")),
		}
		answer := input.AskBool("Are you sure?", &opts)
		assert.True(t, answer)
	})
	t.Run("test other than yes", func(t *testing.T) {
		input := UserInput{
			Writer: ioutil.Discard,
			Reader: bufio.NewReader(bytes.NewBufferString("Nooooooo\n")),
		}
		answer := input.AskBool("Are you sure?", &opts)
		assert.False(t, answer)
	})

}
