/*
Copyright © 2021 The Authors of gotabgo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/groundfoundation/gotabgo"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type opts struct {
	password         string
	server           string
	tls              bool
	username         string
	serverApiVersion string
}

var (
	cfgFile string

	debug bool

	options opts

	tabApi *gotabgo.TabApi

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "gotabgo [serverinfo]",
		Short: "A cli tool for interacting with Tableau Server",
		Long: `
gotabgo is a CLI library for Tableau Server that enables administration from
the command line.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (e error) {

			connectOpt := opts{
				server:           viper.GetString("server"),
				tls:              viper.GetBool("tls"),
				serverApiVersion: viper.GetString("apiversion"),
			}
			tabApi, e = gotabgo.NewTabApi(connectOpt.server,
				connectOpt.serverApiVersion, connectOpt.tls,
				gotabgo.Xml)
			log.Debugf("tabApi struct: %v", tabApi)
			return e
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotabgo.yaml)")
	rootCmd.PersistentFlags().StringVarP(&options.username, "username", "u", "", "username to use when connecting to Tableau Server")
	rootCmd.PersistentFlags().StringVarP(&options.password, "password", "p", "", "password for the user")
	rootCmd.PersistentFlags().StringVarP(&options.server, "server", "s", "", "the hostname of the server")
	rootCmd.Flags().StringVarP(&options.serverApiVersion, "apiversion", "a", "3.9", "specify which version of the api to user")
	rootCmd.Flags().BoolVar(&options.tls, "tls", true, "whether to use TLS or not when connecting")

	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("apiversion", rootCmd.Flags().Lookup("apiversion"))
	viper.BindPFlag("tls", rootCmd.Flags().Lookup("tls"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Turn on debuf output")

	if debug {
		log.SetLevel(log.DebugLevel)
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gotabgo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gotabgo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
