package cmd

import (
	"github.com/spf13/cobra"
	"github.com/st2projects/ssh-sentinel-server/app"
	"github.com/st2projects/ssh-sentinel-server/model"
)

var devMode bool

var httpConfig = model.HTTPConfig{}.Default()

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
	serveCmd.Flags().IntVarP(&httpConfig.Port, "port", "p", 8080, "HTTP Port")
	serveCmd.Flags().StringVarP(&httpConfig.ListenOn, "listen", "l", "0.0.0.0", "Listen Address")
	serveCmd.Flags().BoolVarP(&devMode, "dev-mode", "d", false, "Run in DEV mode. See README for implications")

	serveCmd.MarkFlagRequired("config")
}
