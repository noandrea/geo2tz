package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/noandrea/geo2tz/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: start,
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func start(cmd *cobra.Command, args []string) {
	fmt.Printf(`
  _____           ___  _       
 / ____|         |__ \| |      
| |  __  ___  ___   ) | |_ ____
| | |_ |/ _ \/ _ \ / /| __|_  /
| |__| |  __/ (_) / /_| |_ / / 
 \_____|\___|\___/____|\__/___| version %s
`, rootCmd.Version)
	// Start server
	go func() {
		if err := server.Start(settings); err != nil {
			log.Error("Error starting the server ", err)
			return
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt)
	<-quit
	server.Teardown()
	fmt.Print("Goodbye")
}
