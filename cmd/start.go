package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/noandrea/geo2tz/server"
)

// startCmd represents the start command.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the geo2tz service",
	Long:  ``,
	Run:   start,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func start(*cobra.Command, []string) {
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
			log.Println("Error starting the server ", err)
			return
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	signalChannelLength := 2
	quit := make(chan os.Signal, signalChannelLength)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err := server.Teardown(); err != nil {
		log.Println("error stopping server: ", err)
	}
	fmt.Print("Goodbye")
}
