/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package check

import (
	"github.com/mayron1806/gengo/config"
	"github.com/mayron1806/gengo/constants"
	"github.com/spf13/cobra"
)

var logger *config.Logger

// CheckCmd represents the check command
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if gengo is correctly installed",
	Long:  `Check if gengo is correctly installed`,
	Run: func(cmd *cobra.Command, args []string) {
		err := constants.CheckConstants()
		if err != nil {
			logger.Warn(err)
		}
	},
}

func init() {
	logger = config.GetLogger()
}
