package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	debug bool
	host  string
	port  string
)

var rootCmd = &cobra.Command{
	Use: "ksk",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "set debug to true or false")
	rootCmd.PersistentFlags().StringVar(&host, "host", "0.0.0.0", "webserver host")
	rootCmd.PersistentFlags().StringVar(&port, "port", "8080", "webserver port")
	rootCmd.AddCommand(adminCmd)
	rootCmd.AddCommand(registerCmd)
}
