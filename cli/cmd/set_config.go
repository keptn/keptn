package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/utils/config"

	"github.com/spf13/cobra"
)

var configMng *config.CLIConfigManager

// setConfig implements the config command
var setConfigCmd = &cobra.Command{
	Use:   "set config",
	Short: "Allows to set the CLI configuration",
	Long: `Allows to set the CLI configuration, which is stored in $HOME/.keptn/config.
Therefore, this command takes a key and a new value as arguments. 
	
Example:
	keptn set config AutomaticVersionCheck false`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			cmd.SilenceUsage = false
			return errors.New("required arguments KEY and VALUE")
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		configMng = config.NewCLIConfigManager()
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
				return fmt.Errorf("error when parsing value %v", err)
			}
			cliConfig.AutomaticVersionCheck = val
			newConfig = true
		case "lastversioncheck":
			val, err := time.Parse("RFC3339", args[1])
			if err != nil {
				return fmt.Errorf("error when parsing value %v", err)
			}
			cliConfig.LastVersionCheck = &val
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
