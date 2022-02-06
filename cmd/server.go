package cmd

import (
	"fmt"

	"github.com/emanuelhristea/lime/server"
	"github.com/emanuelhristea/lime/version"
	"github.com/spf13/cobra"
)

var banner = "license server\nversion " + version.Version + "\nhash:" + version.GitCommit

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start license server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(banner)
		fmt.Println(args)
		server.Start(args)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
