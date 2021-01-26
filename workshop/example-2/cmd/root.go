package cmd

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const (
	Iterations = 10
)

func Execute() {
	_ = godotenv.Load(".env")

	var rootCmd = &cobra.Command{
		Use: "kafka",
	}

	producerCommand.
		PersistentFlags().
		StringArrayP(
			"files",
			"f",
			viper.GetStringSlice("CONFIGURATION_FILE"),
			"file pipelines to be initialized")

	producerCommand.
		PersistentFlags().
		IntP(
			"iterations",
			"i",
			Iterations,
			fmt.Sprintf("number of iterations, default: %v", Iterations))

	consumerCommand.
		PersistentFlags().
		StringArrayP(
			"files",
			"f",
			viper.GetStringSlice("CONFIGURATION_FILE"),
			"file pipelines to be initialized")

	rootCmd.AddCommand(producerCommand)
	rootCmd.AddCommand(consumerCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("failed to execute terminal: %v", err.Error())
		os.Exit(1)
	}
}
