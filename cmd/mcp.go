package cmd

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/noandrea/geo2tz/v2/mcpserver"
	"github.com/spf13/cobra"
)

// mcpCmd represents the mcp command.
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "starts the geo2tz MCP server over stdio",
	Long: `Starts a Model Context Protocol (MCP) server over stdin/stdout
exposing the timezone lookup as tools for MCP clients (eg. AI assistants).
Note: to keep the protocol stream clean, log messages go to stderr only.`,
	Run: mcpServe,
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

func mcpServe(*cobra.Command, []string) {
	server, err := mcpserver.NewServer(settings)
	if err != nil {
		log.Println("Error creating the MCP server ", err)
		os.Exit(1)
	}
	if err = server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Println("Error running the MCP server ", err)
		os.Exit(1)
	}
}
