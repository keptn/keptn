// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/spf13/cobra"
)

// installCmd represents the version command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs keptn on your Kubernetes cluster",
	Long: `Installs keptn on your Kubernetes cluster

Example:
	keptn install`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing keptn...")
		execCmd := exec.Command("kubectl", "logs", "control-lxxt6-deployment-569499f6cb-gbzgp", "-n", "keptn", "-c", "user-container", "-f")

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutIn, _ := execCmd.StdoutPipe()
		stderrIn, _ := execCmd.StderrPipe()
		err := execCmd.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}

		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			copyAndCapture(os.Stdout, stdoutIn)
			wg.Done()
		}()

		copyAndCapture(os.Stderr, stderrIn)

		wg.Wait()

		err = execCmd.Wait()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		if errStdout != nil || errStderr != nil {
			log.Fatal("failed to capture stdout or stderr\n")
		}
		outStr, errStr := string(stdout), string(stderr)
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func copyAndCapture(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	/*
		var out []byte
		buf := make([]byte, 1024, 1024)
		for {
			n, err := r.Read(buf[:])
			if n > 0 {
				d := buf[:n]
				out = append(out, d...)
				str := string(d)
				strArr := strings.Split(str, "\n")

				_, err := w.Write(d)
				if err != nil {
					return out, err
				}
			}
			if err != nil {
				// Read returns io.EOF at the end of file, which is not an error for us
				if err == io.EOF {
					err = nil
				}
				return out, err
			}
	*/
}
