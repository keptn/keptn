package cespec

import (
	"fmt"
	"strings"
)

// MarkDown represents a MarkDown (*.md) file content
type MarkDown struct {
	builder *strings.Builder
}

// NewMarkDown creates a new MarkDown instance
func NewMarkDown() *MarkDown {
	md := new(MarkDown)
	md.builder = new(strings.Builder)
	return md
}

func (m *MarkDown) createTitle(title string, level int) string {
	return fmt.Sprintf("%s %s", strings.Repeat("#", level), title)
}
func (m *MarkDown) createLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

func (m *MarkDown) createCodeBlock(code, t string) string {
	return fmt.Sprintf("```%s\n%s\n```\n", t, code)
}

// Write appends a string to the md document
func (m *MarkDown) Write(str string) *MarkDown {
	m.builder.WriteString(str)
	return m
}

// Writeln appends a newline appended string to the md document
func (m *MarkDown) Writeln(str string) *MarkDown {
	m.builder.WriteString(fmt.Sprintf("%s\n", str))
	return m
}

// Title creates and appends a md title with given level to the md document
func (m *MarkDown) Title(title string, level int) *MarkDown {
	m.Write(m.createTitle(title, level))
	m.WriteLineBreak()
	return m
}

// Bullet creates and appends a bullet point to the md document
func (m *MarkDown) Bullet() *MarkDown {
	m.Write("* ")
	return m
}

// WriteLineBreak appends a empty line to the md document
func (m *MarkDown) WriteLineBreak() *MarkDown {
	m.Write("\n")
	return m
}

// MultiBr appends multiple newlines to the md document
func (m *MarkDown) MultiBr(lineBreaks int) *MarkDown {
	m.Write(strings.Repeat("\n", lineBreaks))
	return m
}

// Link creates and appends a md link to the md document
func (m *MarkDown) Link(text, url string) *MarkDown {
	m.Write(m.createLink(text, url))
	m.WriteLineBreak()
	return m
}

// CodeBlock creates and appends a Code block to the md document
func (m *MarkDown) CodeBlock(code, t string) *MarkDown {
	m.Write(m.createCodeBlock(code, t))
	return m
}

// UpLink creates and appends a link to the top of the file
func (m *MarkDown) UpLink() *MarkDown {
	m.Writeln(`([&uarr; up to index](#keptn-cloud-events))`)
	return m
}

// String prints the current state of the md document
func (m *MarkDown) String() string {
	return m.builder.String()
}
