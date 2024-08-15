/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mayron1806/gengo/constants"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if gengo is correctly installed",
	Long:  `Check if gengo is correctly installed`,
	Run: func(cmd *cobra.Command, args []string) {
		err := constants.CheckConstants()
		if err != nil {
			// Logger.Warn(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
