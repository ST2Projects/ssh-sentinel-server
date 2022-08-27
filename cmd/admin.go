package cmd

import (
	"github.com/spf13/cobra"
	"github.com/st2projects/ssh-sentinel-server/app"
)

var create bool
var name string
var username string
var principals []string

// adminCmd represents the serve command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Create / delete users",
	Run: func(cmd *cobra.Command, args []string) {

		app.RunAdmin(configPath, create, name, username, principals)
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)
	adminCmd.Flags().StringVarP(&configPath, "config", "c", "", "Config file")
	adminCmd.Flags().StringVarP(&name, "name", "n", "", "User's name")
	adminCmd.Flags().BoolVarP(&create, "create", "C", false, "If set a new user will be created")
	adminCmd.Flags().StringSliceVarP(&principals, "principals", "P", nil, "A list of principals for the user")
	adminCmd.Flags().StringVarP(&username, "username", "U", "", "Username")

	adminCmd.MarkFlagRequired("config")
	adminCmd.MarkFlagRequired("name")
	adminCmd.MarkFlagRequired("create")
	adminCmd.MarkFlagRequired("principals")
}
