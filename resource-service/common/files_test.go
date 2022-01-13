package common

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

// base64 encoded tgz file containing file test.txt with content "test"
const testTgzContent = "H4sIAOu80mEAA+2Vz0rDQBCHc84zeNgnSGf2b/ZQsFjFkxTRg+BltasGTAtJCjl46N2bZ5/SJ3ATrCW1iqCrlOx3mbAEMpvh+83E1MfWTG0xqGxZRV4AACUEiVTLusIbCEiQo+SCScYYAQqSiojUftrpsigrU7hWbu7nRWZmyZW5vrPFx/fGFyejs9PRweHleJ6bbEbOS1uUnUs6yHvdETAlt9l0iKBSppFzGruThTuhoEGCUKBiqkleZbkduhEhBUiFTkSqKUOJ8X9fIPAjGusHnr+x8n/pnh+e9x6bun/08gRruv6jUpRFRHjuq6Xn/rfzn3SWQFLVv7sI3P+QnH+V/3Qj/ymnEPL/L/hO/jPYkv+SKwUaZVgAO03rvxfr16z8X0af5b/Y8B85kxEBT/106Ln/zeiDwoFAINA/XgHzASEtABIAAA=="

func TestFileSystem_WriteAndReadFile(t *testing.T) {
	// create a tmp directory in test/tmp
	err := ensureDirectoryExists("../test/tmp")
	require.Nil(t, err)
	dir, err := ioutil.TempDir("../test/tmp/", "project-")
	require.Nil(t, err)
	defer func(name string) {
		err = os.RemoveAll(name)
		if err != nil {
			t.Logf("could not delete tmp directory: %v", err)
		}
	}(dir)

	fs := FileSystem{}

	filePath := dir + "/my-file"

	fileContent := "content"

	err = fs.WriteFile(filePath, []byte(fileContent))
	require.Nil(t, err)

	fileExists := fs.FileExists(filePath)
	require.True(t, fileExists)

	res, err := fs.ReadFile(filePath)
	require.Nil(t, err)

	require.Equal(t, fileContent, string(res))

	err = fs.DeleteFile(filePath)
	require.Nil(t, err)

	fileExists = fs.FileExists(filePath)
	require.False(t, fileExists)
}

func TestIsHelmChartPath(t *testing.T) {
	type args struct {
		resourcePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "helm chart file path",
			args: args{
				resourcePath: "helm/chart.tgz",
			},
			want: true,
		},
		{
			name: "not helm chart file path",
			args: args{
				resourcePath: "helm/chart.tgza",
			},
			want: false,
		},
		{
			name: "not helm chart file path",
			args: args{
				resourcePath: "helmy/chart.tgz",
			},
			want: false,
		},
		{
			name: "not helm chart file path",
			args: args{
				resourcePath: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHelmChartPath(tt.args.resourcePath); got != tt.want {
				t.Errorf("IsHelmChartPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileSystem_WriteHelmChart(t *testing.T) {
	// create a tmp directory in test/tmp
	dir, err := ioutil.TempDir("../test/tmp/", "project-")
	require.Nil(t, err)
	defer func(name string) {
		err = os.RemoveAll(name)
		if err != nil {
			t.Logf("could not delete tmp directory: %v", err)
		}
	}(dir)

	fs := NewFileSystem(dir)

	filePath := dir + "/my-file.tgz"

	err = fs.WriteBase64EncodedFile(filePath, testTgzContent)
	require.Nil(t, err)

	err = fs.WriteHelmChart(filePath)
	require.Nil(t, err)

	res, err := fs.ReadFile(dir + "/my-file/test.txt")
	require.Nil(t, err)

	require.Equal(t, "test\n", string(res))
}
