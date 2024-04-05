package main

import (
	"fmt"
	"gopulsar/cmd/publisher/publish"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool

var RootCmd = &cobra.Command{
	Use:   "publisher",
	Short: "CLI to publish new message",
	Long:  "This CLI is to pushing a new message to Pulsar cluster",
	//Run: func(cmd *cobra.Command, args []string) {
	//	log.Println("zo")
	//},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	log.Println("zo1")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	publish.AddCommand(RootCmd)
}
