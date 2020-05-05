/*
Copyright © 2020 Andrea Giacobino <no.andrea@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/noandrea/geo2tz/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool
var settings server.ConfigSchema

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "geo2tz",
	Short: "A rest API service to get the timezone from geo coordinates",
	Long: `Throwing around coordinates to online services seems like not 
  a great idea privacy wise.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	rootCmd.Version = v
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/geo2tz/config.yaml)")
	// for debug logging
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if debug {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/geo2tz")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	server.Defaults()

	// If a config file is found, read it in, else use the defaults
	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(&settings)
		server.Validate(&settings)
		log.Println("Using config file at ", viper.ConfigFileUsed())
	} else {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			viper.Unmarshal(&settings)
		}
	}
	// make the version available via settings
	settings.RuntimeVersion = rootCmd.Version
	log.Debug(fmt.Sprintf("config %#v", settings))
}
