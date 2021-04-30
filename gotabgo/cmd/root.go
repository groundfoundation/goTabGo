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
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type opts struct {
	password string
	server   string
	tls      bool
	username string
}

var (
	cfgFile string

	options opts

	tabApi *gotabgo.TabApi

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "gotabgo",
		Short: "A cli tool for interacting with Tableau Server",
		Long: `
gotabgo is a CLI library for Tableau Server that enables administration from
the command line.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (e error) {

			connectOpt := opts{
				server: viper.GetString("server"),
				tls:    viper.GetBool("tls"),
			}
			tabApi, e = gotabgo.NewTabApi(connectOpt.server, "3.9", connectOpt.tls, gotabgo.Xml)
			fmt.Printf("tabapi: %v", tabApi)
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
	rootCmd.Flags().StringVarP(&options.username, "username", "u", "", "username to use when connecting to Tableau Server")
	rootCmd.Flags().StringVarP(&options.password, "password", "p", "", "password for the user")
	rootCmd.Flags().StringVarP(&options.server, "server", "s", "", "the hostname of the server")
	rootCmd.Flags().BoolVar(&options.tls, "tls", true, "whether to use TLS or not when connecting")

	viper.BindPFlag("username", rootCmd.Flags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.Flags().Lookup("password"))
	viper.BindPFlag("server", rootCmd.Flags().Lookup("server"))
	viper.BindPFlag("tls", rootCmd.Flags().Lookup("tls"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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
