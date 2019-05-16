package cmd

import (
	"fmt"
	"os"

	"github.com/keptn/keptn/cli/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verboseLogging bool
var quietLogging bool
var mocking bool

const authErrorMsg = "This command requires to be authenticated. See \"keptn auth\" for details"

const logo = `                                                                                                                                     
                ##########*                                                                                                                                    
           ,#############    ##                                                                                                                                
       (###############    ####    *                                                                                                                           
    ##################    ###*    ###.                                                                                                                         
   #######      ####    ####    ####                                                                                                                           
   #####          ,   (###    ####    ##                 .&&&&                                                                                                 
  (####   #####      ####    ####    ###                 .&&&&                                                                                                 
  #####    ####    ####    ####    ####                  .&&&&                                                              &&&&&                              
 .######         .###    *###    ####                    .&&&&                                                              &&&&&                              
 ##########     ####    ####    ####    #(               .&&&&                                                              &&&&&                              
 #########    ####    ####    ####    ####               .&&&&       &&&&&/       &&&&&&&&&&/        &&&&&&&&&&&&&%         &&&&&&&&&&&&,     &&&&&&&&&&&&&&   
#########    ####    ####   .###/   /#####               .&&&&     &&&&&&       &&&&&&&&&&&&&&%      &&&&&&&&&&&&&&&&       &&&&&&&&&&&&,     &&&&&&&&&&&&&&&& 
#######    ####    ####    ####    ########              .&&&&   &&&&&&        &&&&&.     /&&&&&     &&&&&       &&&&&(     &&&&&             &&&&&      &&&&&&
 ####(   .###    (###    ####    #########               .&&&& &&&&&&         &&&&&        *&&&&     &&&&&        &&&&&     &&&&&             &&&&&       &&&&&
  ##    ####    ####    ####    ########                 .&&&&&&&&&           &&&&&&&&&&&&&&&&&&     &&&&&         &&&&&    &&&&&             &&&&&       &&&&&
      ####    ####    ####    #########                  .&&&&&&&&&&          &&&&&&&&&&&&&&&&&&     &&&&&         &&&&&    &&&&&             &&&&&       &&&&&
     ####    ###/   (###,   (########                    .&&&&  &&&&         &&&&&                   &&&&&         &&&&&    &&&&&             &&&&&       &&&&&
           ####    ####    ########*                     .&&&&   .&&&&&       #&&&&&                 &&&&&        &&&&&     /&&&&             &&&&&       &&&&&
         ####    ####    #########                       .&&&&     &&&&&&      &&&&&&&%    ,&&&      &&&&&&&( .&&&&&&&       &&&&&&/  %&&     &&&&&       &&&&&
          ##    ####    ########                         .&&&&       &&&&      &&&&&&&&&&&&&&        &&&&&&&&&&&&&&&&         &&&&&&&&&&&     &&&&&       &&&&&
                                                                                    .&&&&&&&*        &&&&&  *&&&&                 (&&%                         
                                                                                                     &&&&&                                                     
                                                                                                     &&&&&                                                     
                                                                                                     &&&&&`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "keptn",
	Short: "This is a CLI for using keptn",
	Long: `This is a CLI for using keptn. The CLI allows to authenticate against keptn, to configure your Github organization,
to create projects, and to onboard services.
	
	` + logo,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseLogging, "verbose", "v", false, "verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&quietLogging, "quiet", "q", false, "suppress debug and info output")
	rootCmd.PersistentFlags().BoolVarP(&mocking, "mock", "", false, "mocking of server communication - ATTENTION: your commands won't be sent to the keptn server")

	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	utils.LogLevel = utils.InfoLevel
	if verboseLogging && quietLogging {
		fmt.Println("Verbose logging and quiet output are mutually exclusive flags. Please use only one.")
		os.Exit(1)
	}
	if verboseLogging {
		utils.LogLevel = utils.VerboseLevel
	}
	if quietLogging {
		utils.LogLevel = utils.QuietLevel
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		utils.PrintLog(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()), utils.InfoLevel)
	}
}
