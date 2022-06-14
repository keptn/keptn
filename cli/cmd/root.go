package cmd

import (
	"fmt"
	"github.com/keptn/keptn/cli/pkg/config"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/pkg/version"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var cfgFile string
var verboseLogging bool
var quietLogging bool
var mocking bool
var insecureSkipTLSVerify bool
var kubectlOptions string
var namespace string
var assumeYes bool
var help bool

const authErrorMsg = "This command requires to be authenticated. See \"keptn auth\" for details"

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRootCommand(version.NewVersionChecker())
var rootCLIConfig config.CLIConfig

func NewRootCommand(vChecker *version.VersionChecker) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "keptn",
		Short: "The CLI for using Keptn",
		Long: `The CLI allows interaction with a Keptn installation to manage Keptn, to trigger workflows, and to get details.
	`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			runVersionCheck(vChecker, os.Args[1:], rootCLIConfig)
		},
	}
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	// Set LogLevel to QuietLevel
	currentLogLevel := logging.LogLevel
	logging.LogLevel = logging.QuietLevel
	// Set LogLevel back to previous state
	logging.LogLevel = currentLogLevel

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verboseLogging, "verbose", "v", false, "Enables verbose logging to print debug messages")
	rootCmd.PersistentFlags().BoolVarP(&quietLogging, "quiet", "q", false, "Suppresses debug and info messages")
	rootCmd.PersistentFlags().BoolVarP(&mocking, "mock", "", false, "Disables communication to a Keptn endpoint")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "keptn",
		"Specify the namespace where Keptn should be installed, used and uninstalled")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config-file", "", "",
		"Specify custom Keptn Config file path (default: ~/.keptn/config)")
	rootCmd.PersistentFlags().BoolVarP(&assumeYes, "yes", "y", false, "Assume yes for all user prompts")
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "help")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logging.LogLevel = logging.InfoLevel
	if verboseLogging && quietLogging {
		fmt.Println("Verbose logging and quiet output are mutually exclusive flags. Please use only one.")
		os.Exit(1)
	}
	if verboseLogging {
		logging.LogLevel = logging.VerboseLevel
	}
	if quietLogging {
		logging.LogLevel = logging.QuietLevel
	}

	cfgMgr := config.NewCLIConfigManager(cfgFile)

	logging.PrintLog(fmt.Sprintf("Using config file: %s", cfgMgr.CLIConfigPath), logging.VerboseLevel)

	var err error
	rootCLIConfig, err = cfgMgr.LoadCLIConfig()
	if err != nil {
		logging.PrintLog(err.Error(), logging.InfoLevel)
	}

}

type options []string

func (s *options) appendIfNotEmpty(newOption string) {
	if newOption != "" {
		*s = append(*s, newOption)
	}
}

// passing flags and cliConfig as arguments makes it easy to test this function
func runVersionCheck(vChecker *version.VersionChecker, flags []string, cliConfig config.CLIConfig) {
	// Don't check version if AutomaticVersionCheck is disabled
	// Server version won't be available during `install`
	// Server version should not be needed when using oauth
	if !cliConfig.AutomaticVersionCheck || skipVersionCheck(flags) {
		return
	}

	var cliMsgPrinted, cliChecked, keptnMsgPrinted, keptnChecked bool

	cliChecked, cliMsgPrinted = vChecker.CheckCLIVersion(Version, true)

	if cliMsgPrinted {
		fmt.Println("* Your Keptn CLI version: " + Version)
	}

	clusterVersion, err := getKeptnServerVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "* Warning: could not check Keptn server version: %s\n", err.Error())
	} else {
		kvChecker := version.NewKeptnVersionChecker()
		keptnChecked, keptnMsgPrinted = kvChecker.CheckKeptnVersion(Version, clusterVersion, true)
		if keptnMsgPrinted {
			fmt.Fprintf(os.Stderr, "* Your Keptn cluster version: %s\n", clusterVersion)
		}

		if clusterVersion != Version {
			fmt.Fprintf(os.Stderr, "* Warning: Your Keptn CLI version (%s) and Keptn cluster version (%s) don't match. This can lead to problems. Please make sure to use the same versions.\n", Version, clusterVersion)
		}
	}

	if cliMsgPrinted || keptnMsgPrinted {
		fmt.Fprintf(os.Stderr, setVersionCheckMsg, "disable", "false")
	}

	if cliChecked || keptnChecked {
		updateLastVersionCheck()
	}
}

// skipVersionCheck checks if the subcommand requires to skip the version check step
// (for now this is true in case of  `install` or `--oauth`)
// args here does not contain the main command
// e.g., For `keptn -q install`, args would be just ['-q', 'install']
func skipVersionCheck(args []string) bool {
	for _, arg := range args {
		switch {
		case arg == "--oauth":
			return true
		case strings.HasPrefix(arg, "-"):
			continue
		case arg == "install":
			return true
		default:
			return false
		}
	}
	return false
}
