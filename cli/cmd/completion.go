package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long:  "Installing bash completion on macOS using homebrew \n" +
		"	If running Bash 3.2 included with macOS, run brew install bash-completion \n" +
		"	If running Bash 4.1+, run brew install bash-completion@2 \n" +
		"	If you've installed via other means, you may need add the completion to your completion directory, run keptn completion bash > $(brew --prefix)/etc/bash_completion.d/keptn \n\n" +
		"Installing bash completion on Linux \n" +
		"If bash-completion is not installed on Linux, please install the 'bash-completion' package via your distribution's package manager. \n" +
		"Load the keptn completion code for bash into the current shell \n" +
		"	source <(keptn completion bash) \n" +
		"Write bash completion code to a file and source it from .bash_profile \n" +
		"	keptn completion bash > /etc/bash_completion.d/keptn \n" +
		"If you are a ZSH User \n" +
		"Load the keptn completion code for zsh[1] into the current shell \n" +
		"	source <(keptn completion zsh) \n" +
		"Set the keptn completion code for zsh[1] to autoload on startup \n" +
		"	keptn completion zsh > \"${fpath[1]}/_keptn",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func (cmd *cobra.Command, args []string){
	switch args[0]{
case "bash":
	cmd.Root().GenBashCompletion(os.Stdout)
case "zsh":
	cmd.Root().GenZshCompletion(os.Stdout)
case "fish":
	cmd.Root().GenFishCompletion(os.Stdout, true)
case "powershell":
	cmd.Root().GenPowerShellCompletion(os.Stdout)
}
},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
