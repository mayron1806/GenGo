/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package template

import (
	"github.com/mayron1806/gengo/cmd/template/new"
	"github.com/spf13/cobra"
)

// TemplateCmd represents the template command
var TemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "A brief description of your command",
	Long:  ``,
}

func init() {
	TemplateCmd.AddCommand(new.NewCmd)
}
