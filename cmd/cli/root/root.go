package root

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	PulsarURL = "pulsar://localhost:6650"
)

var rootCmd = &cobra.Command{
	Use:   "gopulsar",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gopulsar.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
