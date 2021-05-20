package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/pkg/config"

	"github.com/spf13/cobra"
)

var configMng *config.CLIConfigManager

// setConfig implements the config command
var setConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Sets flags of the CLI configuration",
	Long: `Sets flags of the CLI configuration, which is stored in $HOME/.keptn/config.

*	This command takes a key and a new value as arguments. 
`,
	Example:      `keptn set config AutomaticVersionCheck false`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			cmd.SilenceUsage = false
			return errors.New("required arguments KEY and VALUE")
		}
		configMng = config.NewCLIConfigManager()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		cliConfig, err := configMng.LoadCLIConfig()
		if err != nil {
			return err
		}

		key := strings.ToLower(args[0])
		newConfig := false
		switch key {
		case "automaticversioncheck":
			val, err := strconv.ParseBool(args[1])
			if err != nil {
				return fmt.Errorf("error when parsing value %w", err)
			}
			cliConfig.AutomaticVersionCheck = val
			newConfig = true
		case "lastversioncheck":
			val, err := time.Parse("RFC3339", args[1])
			if err != nil {
				return fmt.Errorf("error when parsing value %w", err)
			}
			cliConfig.LastVersionCheck = &val
			newConfig = true
		case "kubecontextcheck":
			val, err := strconv.ParseBool(args[1])
			if err != nil {
				return fmt.Errorf("error when parsing value %w", err)
			}
			cliConfig.KubeContextCheck = val
			newConfig = true
		default:
			return fmt.Errorf("Unsupported key %s", args[0])
		}

		if newConfig {
			return configMng.StoreCLIConfig(cliConfig)
		}

		return nil
	},
}

func init() {
	setCmd.AddCommand(setConfigCmd)
}
