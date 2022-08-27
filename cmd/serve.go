package cmd

import (
	"github.com/spf13/cobra"
	"github.com/st2projects/ssh-sentinel-server/app"
)

var devMode bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the CA server",
	Run: func(cmd *cobra.Command, args []string) {
		app.InitialiseApp(configPath, devMode)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "", "Config file")
	serveCmd.Flags().BoolVarP(&devMode, "dev-mode", "d", false, "Run in DEV mode. See README for implications")

	serveCmd.MarkFlagRequired("config")
}
