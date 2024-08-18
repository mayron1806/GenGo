/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package new

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mayron1806/gengo/constants"
	"github.com/mayron1806/gengo/internal/components"
	"github.com/mayron1806/gengo/internal/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type createMode uint8

const (
	MANUAL createMode = iota
	INTERACTIVE
	IMPORT // comming soon
)

var (
	name   string
	vars   []string
	global bool
	mode   createMode
)

type templateModel struct {
	Name      string   `json:"name"`
	Version   string   `json:"version"`
	Global    bool     `json:"global"`
	Variables []string `json:"variables"`
}

// TemplateCmd represents the template command
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// add name if not provided

		if name == "" {
			inputName := components.NewInputText(components.InputTextOptions{
				Placeholder: "my-template",
				Label:       "Please enter the name of the template:",
				MinLength:   1,
			})
			p := tea.NewProgram(inputName, tea.WithAltScreen())
			m, err := p.Run()
			if err != nil {
				fmt.Println(err)
			}

			if m, ok := m.(components.InputText); ok {
				if m.Quitting {
					return
				}
				name = m.TextInput.Value()
			}
		}

		items := make(map[createMode]string, 0)
		items[MANUAL] = "manual"
		items[INTERACTIVE] = "interactive"
		// set template creation mode
		list := components.NewList(components.ListOptions{
			Items: []string{items[MANUAL], items[INTERACTIVE]},
			Title: "Select Template Creation Mode",
		})

		p := tea.NewProgram(list, tea.WithAltScreen())
		m, err := p.Run()
		if err != nil {
			fmt.Println(err)
		}
		if m, ok := m.(components.List); ok {
			if m.Quitting {
				return
			}
			selectedMode := m.Choice
			switch selectedMode {
			case items[MANUAL]:
				mode = MANUAL
			case items[INTERACTIVE]:
				mode = INTERACTIVE
			}
		}
		gengoLocalRoot, err := helper.FindGengoLocalRoot()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// load local config
		viper.SetConfigName("gengo")
		viper.SetConfigType("json")
		viper.AddConfigPath(gengoLocalRoot)
		err = viper.ReadInConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// version
		version := viper.GetString("version")

		// list templates
		fmt.Printf("Gengo Version: %s\n", version)
		var templates []templateModel
		viper.UnmarshalKey("templates", &templates)

		var workDir string
		if global {
			var err error
			workDir, err = constants.GetGengoHome()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			workDir = gengoLocalRoot
		}
		newTemplate := templateModel{
			Name:      name,
			Version:   "1.0.0",
			Global:    global,
			Variables: vars,
		}
		// create template folder
		err = os.MkdirAll(filepath.Join(workDir, "templates", name), 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		templates = append(templates, newTemplate)
		viper.Set("templates", templates)
		err = viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Template %s created successfully in %s\n", name, workDir)
	},
}

func init() {
	NewCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the template")
	viper.BindPFlag("name", NewCmd.Flags().Lookup("name"))

	NewCmd.Flags().StringSliceVarP(&vars, "variables", "v", []string{}, "Template variables")
	viper.BindPFlag("vars", NewCmd.Flags().Lookup("vars"))

	NewCmd.Flags().BoolVarP(&global, "global", "g", false, "Global template")
	viper.BindPFlag("global", NewCmd.Flags().Lookup("global"))
}
