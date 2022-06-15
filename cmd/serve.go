package cmd

import (
	"github.com/spf13/cobra"
	"ssh-sentinel-server/server"
	"strconv"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the CA server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := strconv.Atoi(cmd.Flags().Lookup("port").Value.String())
		config := cmd.Flags().Lookup("config").Value.String()

		server.Serve(port, config)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 8080, "Port to run the service on")
	serveCmd.Flags().StringP("config", "c", "", "Config file")
	serveCmd.MarkFlagRequired("config")
}
