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
	Short: "Start the geo2tz server",
	Run:   start,
}

func init() {
	rootCmd.AddCommand(startCmd)
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
	if err := server.Teardown(); err != nil {
		log.Error("error stopping server: ", err)
	}
	fmt.Print("Goodbye")
}
