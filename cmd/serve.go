package cmd

import (
	"github.com/spf13/cobra"
	"github.com/st2projects/ssh-sentinel-server/app"
	"github.com/st2projects/ssh-sentinel-server/model"
)

var devMode bool

var httpConfig = model.HTTPConfig{}.Default()

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the CA server",
	Run: func(cmd *cobra.Command, args []string) {
		app.InitialiseApp(configPath, devMode, httpConfig)
	},
}

func init() {

	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "", "Config file")
	serveCmd.Flags().IntVarP(&httpConfig.HttpsPort, "https-port", "s", 443, "HTTPS Port")
	serveCmd.Flags().IntVarP(&httpConfig.HttpPort, "http-port", "i", 80, "HTTP Port")
	serveCmd.Flags().BoolVarP(&devMode, "dev-mode", "d", false, "Run in DEV mode. See README for implications")

	serveCmd.MarkFlagRequired("config")
}
